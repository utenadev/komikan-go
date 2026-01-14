package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kench/komikan-go/internal/api"
	"github.com/kench/komikan-go/internal/manga"
)

func main() {
	var (
		isbn   = flag.String("isbn", "", "ISBN code to add")
		list   = flag.Bool("list", false, "List all manga")
	)

	flag.Parse()

	mgr := manga.NewManager()

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
			fmt.Printf("- %s Vol.%d (%s) - %s\n", b.Title, b.Volume, b.Author, b.ISBN)
		}
		return
	}

	if *isbn != "" {
		// Add manga by ISBN
		fmt.Printf("Looking up ISBN: %s\n", *isbn)

		// TODO: Load app ID from config
		client := api.NewRakutenClient("YOUR_APP_ID_HERE")

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
	os.Exit(1)
}
