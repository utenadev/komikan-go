package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kench/komikan-go/internal/api"
	"github.com/kench/komikan-go/internal/db"
	"github.com/kench/komikan-go/internal/manga"
)

func main() {
	var (
		isbn    = flag.String("isbn", "", "ISBN code to add")
		list    = flag.Bool("list", false, "List all manga")
		latest  = flag.String("latest", "", "Check latest volume for a title")
		dbPath  = flag.String("db", "data/komikan.db", "Database path")
		appID   = flag.String("app-id", "", "Rakuten Application ID (or set RAKUTEN_APP_ID env var)")
	)

	flag.Parse()

	// Initialize database
	database, err := db.NewDB(db.Config{Path: *dbPath})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	mgr := manga.NewManager(database)

	if *list {
		// List all manga
		books, err := mgr.List()
		if err != nil {
			log.Fatalf("Failed to list manga: %v", err)
		}

		if len(books) == 0 {
			fmt.Println("No manga registered yet.")
			return
		}

		fmt.Println("Registered Manga:")
		fmt.Println("==================")
		for _, b := range books {
			if b.Series != "" {
				fmt.Printf("- %s Vol.%d [%s] (%s) - %s\n", b.Title, b.Volume, b.Series, b.Author, b.ISBN)
			} else {
				fmt.Printf("- %s (%s) - %s\n", b.Title, b.Author, b.ISBN)
			}
		}
		return
	}

	if *latest != "" {
		// Get Rakuten App ID
		appID := getRakutenAppID(*appID)
		if appID == "" {
			log.Fatal("Rakuten Application ID is required. Use -app-id flag or set RAKUTEN_APP_ID env var")
		}

		// Check latest volume
		fmt.Printf("Checking latest volume for: %s\n", *latest)

		client := api.NewRakutenClient(appID)
		books, err := client.SearchByTitleSorted(*latest, "-releaseDate", 30)
		if err != nil {
			log.Fatalf("Failed to search: %v", err)
		}

		if len(books) == 0 {
			fmt.Println("No results found.")
			return
		}

		// Find latest numbered volume
		var latestVolume int
		var latestBook *api.BookInfo
		for _, book := range books {
			info := manga.ExtractVolumeInfo(book.Title)
			if info.HasVolume && info.Volume > latestVolume {
				latestVolume = info.Volume
				latestBook = &book
			}
		}

		if latestBook != nil {
			info := manga.ExtractVolumeInfo(latestBook.Title)
			fmt.Printf("\n📚 Latest Volume Found:\n")
			fmt.Printf("  Title: %s\n", latestBook.Title)
			fmt.Printf("  Volume: %d\n", info.Volume)
			fmt.Printf("  Author: %s\n", latestBook.Author)
			fmt.Printf("  Publisher: %s\n", latestBook.Publisher)
			fmt.Printf("  ISBN: %s\n", latestBook.Isbn)
			fmt.Printf("  Release Date: %s\n", latestBook.SalesDate)
			fmt.Printf("  URL: %s\n", latestBook.ItemURL)
		} else {
			fmt.Println("No numbered volumes found.")
		}

		// Show total count
		fmt.Printf("\nFound %d total results for \"%s\"\n", len(books), *latest)
		return
	}

	if *isbn != "" {
		// Get Rakuten App ID
		appID := getRakutenAppID(*appID)
		if appID == "" {
			log.Fatal("Rakuten Application ID is required. Use -app-id flag or set RAKUTEN_APP_ID env var")
		}

		// Add manga by ISBN
		fmt.Printf("Looking up ISBN: %s\n", *isbn)

		client := api.NewRakutenClient(appID)

		book, err := client.SearchByISBN(*isbn)
		if err != nil {
			log.Fatalf("Failed to find book: %v", err)
		}

		m := manga.Manga{
			ID:          *isbn,
			Title:       book.Title,
			Author:      book.Author,
			Publisher:   book.Publisher,
			ISBN:        book.Isbn,
			PublishDate: book.SalesDate,
			URL:         book.ItemURL,
		}

		// Extract volume info
		volInfo := manga.ExtractVolumeInfo(book.Title)
		if volInfo.HasVolume {
			m.Volume = volInfo.Volume
			m.Series = volInfo.Title
		}

		if err := mgr.Add(m); err != nil {
			log.Fatalf("Failed to add manga: %v", err)
		}

		fmt.Printf("Added: %s (%s)\n", m.Title, m.Author)
		return
	}

	// Show usage
	fmt.Println("Komikan CLI - Manga Management Tool")
	fmt.Println("\nUsage:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  komikan-cli -isbn 9784088818791")
	fmt.Println("  komikan-cli -list")
	fmt.Println("  komikan-cli -latest ダンダダン")
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("  RAKUTEN_APP_ID  Rakuten Application ID")
	os.Exit(1)
}

func getRakutenAppID(fromFlag string) string {
	if fromFlag != "" {
		return fromFlag
	}
	return os.Getenv("RAKUTEN_APP_ID")
}
