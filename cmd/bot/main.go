package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kench/komikan-go/internal/config"
	"github.com/kench/komikan-go/internal/db"
	"github.com/kench/komikan-go/internal/manga"
	"github.com/kench/komikan-go/internal/nostr"
)

var (
	version = "dev"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Configuration file path")
	printVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Printf("komikan-bot v%s\n", version)
		os.Exit(0)
	}

	log.Printf("Starting Komikan Bot v%s...", version)

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override with environment variables
	cfg.LoadFromEnv()

	// Validate configuration
	if cfg.Nostr.SecretKey == "" {
		log.Fatal("Nostr secret key is required. Set it in config.yaml or NOSTR_SECRET_KEY env var")
	}
	if cfg.Rakuten.ApplicationID == "" {
		log.Fatal("Rakuten Application ID is required. Set it in config.yaml or RAKUTEN_APP_ID env var")
	}

	// Initialize database
	database, err := db.NewDB(db.Config{Path: cfg.Database.Path})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize Nostr client
	client, err := nostr.NewClient(nostr.Config{
		SecretKey: cfg.Nostr.SecretKey,
		Relays:    cfg.Nostr.Relays,
	})
	if err != nil {
		log.Fatalf("Failed to create Nostr client: %v", err)
	}

	// Connect to relays
	log.Println("Connecting to Nostr relays...")
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to relays: %v", err)
	}

	// Get and display public key
	npub, err := client.GetPublicKey()
	if err != nil {
		log.Printf("Warning: failed to get public key: %v", err)
	} else {
		log.Printf("Bot public key: %s", npub)
	}

	// Post startup announcement
	if err := client.Publish("📚 Komikan Bot is now running!"); err != nil {
		log.Printf("Warning: failed to publish startup message: %v", err)
	}

	log.Println("Bot is running. Press Ctrl+C to stop.")

	// Start periodic checks if enabled
	if cfg.Bot.AnnounceNewReleases {
		go runPeriodicChecks(client, database, cfg)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	client.Disconnect()
	fmt.Println("Bye!")
}

func runPeriodicChecks(client *nostr.Client, database *db.DB, cfg *config.Config) {
	// Parse check interval
	interval, err := time.ParseDuration(cfg.Bot.CheckInterval)
	if err != nil {
		log.Printf("Invalid check interval: %v, using 1 hour", err)
		interval = time.Hour
	}

	// Initial check on startup
	checkAndAnnounceNewReleases(client, database, cfg.Rakuten.ApplicationID)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		checkAndAnnounceNewReleases(client, database, cfg.Rakuten.ApplicationID)
	}
}

func checkAndAnnounceNewReleases(client *nostr.Client, database *db.DB, rakutenAPIKey string) {
	log.Println("Checking for new releases...")

	mgr := manga.NewManager(database)
	newReleases, err := mgr.CheckNewReleases(rakutenAPIKey)
	if err != nil {
		log.Printf("Failed to check for new releases: %v", err)
		return
	}

	if len(newReleases) == 0 {
		log.Println("No new releases found.")
		return
	}

	log.Printf("Found %d new release(s)!", len(newReleases))

	// Announce each new release
	for _, release := range newReleases {
		message := formatNewReleaseMessage(release)
		if err := client.Publish(message); err != nil {
			log.Printf("Failed to publish announcement: %v", err)
		} else {
			log.Printf("Announced: %s Vol.%d", release.SeriesTitle, release.NewVolume)
		}
	}
}

func formatNewReleaseMessage(release manga.NewReleaseCheckResult) string {
	return fmt.Sprintf("📖 新刊情報！\n\n"+
		"%s Vol.%d が発売予定です！\n"+
		"📅 発売日: %s\n"+
		"👨‍🎨 作者: %s\n"+
		"🔗 %s",
		release.SeriesTitle,
		release.NewVolume,
		release.SalesDate,
		release.Author,
		release.URL)
}
