package repository

import (
	"fmt"
	"log"
	"time"

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

// FilterBooksByDate filters books by year and/or month.
// Pass 0 to skip filtering by that field.
// Books whose date cannot be parsed are excluded when any filter is active.
func FilterBooksByDate(books []entity.Book, year, month int) []entity.Book {
	if year == 0 && month == 0 {
		return books
	}
	var filtered []entity.Book
	no := 1
	for _, b := range books {
		t, err := parseBookDate(b.Date)
		if err != nil {
			continue
		}
		if year != 0 && t.Year() != year {
			continue
		}
		if month != 0 && int(t.Month()) != month {
			continue
		}
		b.No = no
		filtered = append(filtered, b)
		no++
	}
	return filtered
}

// parseBookDate tries several date formats used by Bookmeter.
func parseBookDate(date string) (time.Time, error) {
	formats := []string{"2006/01/02", "2006-01-02"}
	for _, f := range formats {
		t, err := time.Parse(f, date)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %q", date)
}

func readListURL(userID string, page int) string {
	return fmt.Sprintf("https://bookmeter.com/users/%s/books/read?page=%d", userID, page)
}
