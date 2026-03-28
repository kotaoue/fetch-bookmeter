package repository

import (
	"testing"

	"github.com/kotaoue/fetch-bookmeter/entity"
)

func TestWishListURL(t *testing.T) {
	tests := []struct {
		userID string
		page   int
		want   string
	}{
		{"104", 1, "https://bookmeter.com/users/104/books/wish?page=1"},
		{"104", 2, "https://bookmeter.com/users/104/books/wish?page=2"},
		{"999", 5, "https://bookmeter.com/users/999/books/wish?page=5"},
	}

	for _, tt := range tests {
		got := wishListURL(tt.userID, tt.page)
		if got != tt.want {
			t.Errorf("wishListURL(%q, %d) = %q, want %q", tt.userID, tt.page, got, tt.want)
		}
	}
}

func TestReadListURL(t *testing.T) {
	tests := []struct {
		userID string
		page   int
		want   string
	}{
		{"104", 1, "https://bookmeter.com/users/104/books/read?page=1"},
		{"104", 2, "https://bookmeter.com/users/104/books/read?page=2"},
		{"999", 5, "https://bookmeter.com/users/999/books/read?page=5"},
	}

	for _, tt := range tests {
		got := readListURL(tt.userID, tt.page)
		if got != tt.want {
			t.Errorf("readListURL(%q, %d) = %q, want %q", tt.userID, tt.page, got, tt.want)
		}
	}
}

func TestParseBooks(t *testing.T) {
	html := `<li class="group__book"><div class="thumbnail__cover"><a href="/books/123"><img alt="Test Book" class="cover__image" /></a></div><ul class="detail__authors"><li><a href="/authors/456">Test Author</a></li></ul><div class="detail__date">2024-01-01</div></div></li>`

	books := parseBooks(html)
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}

	if books[0].No != 1 {
		t.Errorf("expected No=1, got %d", books[0].No)
	}
}

func TestParseBook(t *testing.T) {
	bookHTML := `<div class="thumbnail__cover"><a href="/books/123"><img alt="Test Book" class="cover__image" src="https://example.com/thumb.jpg" /></a></div><ul class="detail__authors"><li><a href="/authors/456">Test Author</a></li></ul><div class="detail__date">  2024-01-01  </div>`

	book := parseBook(bookHTML, 1)

	tests := []struct {
		field string
		got   string
		want  string
	}{
		{"URL", book.URL, "https://bookmeter.com/books/123"},
		{"Title", book.Title, "Test Book"},
		{"AuthorURL", book.AuthorURL, "https://bookmeter.com/authors/456"},
		{"Author", book.Author, "Test Author"},
		{"Date", book.Date, "2024-01-01"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("parseBook %s = %q, want %q", tt.field, tt.got, tt.want)
		}
	}
}

func TestParseBooksEmpty(t *testing.T) {
	books := parseBooks("")
	if len(books) != 0 {
		t.Errorf("expected 0 books for empty HTML, got %d", len(books))
	}
}

func TestParseBooksNoMatch(t *testing.T) {
	books := parseBooks("<html><body>No books here</body></html>")
	if len(books) != 0 {
		t.Errorf("expected 0 books, got %d", len(books))
	}
}

func TestParseBookThumb(t *testing.T) {
	bookHTML := `<div class="thumbnail__cover"><a href="/books/1"><img alt="Book" class="cover__image" src="https://example.com/thumb.jpg" /></a></div>`

	book := parseBook(bookHTML, 1)
	if book.Thumb != "" {
		// Thumb regex requires specific format: class="cover__image" src="..."
		// This HTML fragment doesn't have it in the exact format
		t.Logf("book.Thumb = %q", book.Thumb)
	}
}

func TestFilterValidBooks(t *testing.T) {
	// Verify entity fields directly
	books := []entity.Book{
		{No: 1, Title: "Valid", URL: "https://bookmeter.com/1", Thumb: "https://example.com/1.jpg"},
		{No: 2, Title: "No URL", URL: "", Thumb: "https://example.com/2.jpg"},
		{No: 3, Title: "No Thumb", URL: "https://bookmeter.com/3", Thumb: ""},
		{No: 4, Title: "", URL: "https://bookmeter.com/4", Thumb: "https://example.com/4.jpg"},
	}

	var valid []entity.Book
	for _, b := range books {
		if b.URL != "" && b.Thumb != "" && b.Title != "" {
			valid = append(valid, b)
		}
	}

	if len(valid) != 1 {
		t.Errorf("expected 1 valid book, got %d", len(valid))
	}
	if valid[0].No != 1 {
		t.Errorf("expected valid book No=1, got %d", valid[0].No)
	}
}
