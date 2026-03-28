package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kotaoue/fetch-bookmeter/repository"
)

// RunFetchRead parses flags and fetches the read books list from Bookmeter
func RunFetchRead(args []string) error {
	fs := flag.NewFlagSet("fetch-read", flag.ExitOnError)
	userID := fs.String("user-id", "104", "Bookmeter user ID")
	output := fs.String("output", "read.json", "Output file path for read.json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return fetchAndSaveReadList(*userID, *output)
}

func fetchAndSaveReadList(userID, outputFile string) error {
	books, err := repository.FetchReadList(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch read list: %w", err)
	}

	jsonData, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	log.Printf("✓ Success! Read list saved to %s", outputFile)
	log.Printf("✓ Total read books: %d", len(books))

	return nil
}
