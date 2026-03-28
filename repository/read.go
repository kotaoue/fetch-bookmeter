package repository

import (
	"fmt"
	"log"

	"github.com/kotaoue/fetch-bookmeter/entity"
)

// FetchReadList fetches and parses all pages of the read books list from Bookmeter
func FetchReadList(userID string) ([]entity.Book, error) {
	var allBooks []entity.Book
	no := 1

	for page := 1; ; page++ {
		url := readListURL(userID, page)
		htmlContent, err := fetchHTML(url, 3)
		if err != nil {
			log.Printf("Stopping read list pagination at page %d: %v", page, err)
			break
		}

		books := parseBooks(htmlContent)
		if len(books) == 0 {
			break
		}

		for _, b := range books {
			b.No = no
			allBooks = append(allBooks, b)
			no++
		}
		log.Printf("Parsed %d books from read list page %d", len(books), page)
	}

	log.Printf("Parsed %d books total from read list", len(allBooks))
	return allBooks, nil
}

func readListURL(userID string, page int) string {
	return fmt.Sprintf("https://bookmeter.com/users/%s/books/read?page=%d", userID, page)
}
