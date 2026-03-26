package service

import (
	"os"
	"strings"
	"testing"

	"github.com/kotaoue/fetch-bookmeter/entity"
)

func TestFilterValidBooks(t *testing.T) {
	books := []entity.Book{
		{No: 1, Title: "Valid", URL: "https://bookmeter.com/1", Thumb: "https://example.com/1.jpg"},
		{No: 2, Title: "No URL", URL: "", Thumb: "https://example.com/2.jpg"},
		{No: 3, Title: "No Thumb", URL: "https://bookmeter.com/3", Thumb: ""},
		{No: 4, Title: "", URL: "https://bookmeter.com/4", Thumb: "https://example.com/4.jpg"},
	}

	valid := filterValidBooks(books)
	if len(valid) != 1 {
		t.Errorf("expected 1 valid book, got %d", len(valid))
	}
	if valid[0].No != 1 {
		t.Errorf("expected valid book No=1, got %d", valid[0].No)
	}
}

func TestBuildBookHTML(t *testing.T) {
	book := entity.Book{
		URL:   "https://bookmeter.com/books/123",
		Thumb: "https://example.com/thumb.jpg",
		Title: "Test Book",
	}

	got := buildBookHTML(book)
	want := `<a href="https://bookmeter.com/books/123"><img src="https://example.com/thumb.jpg" alt="Test Book" width="128px"></a>`
	if got != want {
		t.Errorf("buildBookHTML = %q, want %q", got, want)
	}
}

func TestBuildBookHTMLEscaping(t *testing.T) {
	book := entity.Book{
		URL:   `https://bookmeter.com/books/123?a=1&b=2`,
		Thumb: "https://example.com/thumb.jpg",
		Title: `Book with <special> "chars"`,
	}

	got := buildBookHTML(book)
	if strings.Contains(got, "&") && !strings.Contains(got, "&amp;") {
		t.Errorf("buildBookHTML did not escape ampersand in URL")
	}
	if strings.Contains(got, `<special>`) {
		t.Errorf("buildBookHTML did not escape angle brackets in title")
	}
}

func TestReplaceBetweenMarkers(t *testing.T) {
	content := "before<!-- WISH_BOOK_START -->old content<!-- WISH_BOOK_END -->after"
	got, err := replaceBetweenMarkers(content, "<!-- WISH_BOOK_START -->", "<!-- WISH_BOOK_END -->", "new content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "before<!-- WISH_BOOK_START -->new content<!-- WISH_BOOK_END -->after"
	if got != want {
		t.Errorf("replaceBetweenMarkers = %q, want %q", got, want)
	}
}

func TestReplaceBetweenMarkersNoStart(t *testing.T) {
	content := "before<!-- WISH_BOOK_END -->after"
	_, err := replaceBetweenMarkers(content, "<!-- WISH_BOOK_START -->", "<!-- WISH_BOOK_END -->", "new")
	if err == nil {
		t.Error("expected error when start marker missing")
	}
}

func TestReplaceBetweenMarkersNoEnd(t *testing.T) {
	content := "before<!-- WISH_BOOK_START -->after"
	_, err := replaceBetweenMarkers(content, "<!-- WISH_BOOK_START -->", "<!-- WISH_BOOK_END -->", "new")
	if err == nil {
		t.Error("expected error when end marker missing")
	}
}

func TestReplaceBetweenMarkersWrongOrder(t *testing.T) {
	content := "before<!-- WISH_BOOK_END -->middle<!-- WISH_BOOK_START -->after"
	_, err := replaceBetweenMarkers(content, "<!-- WISH_BOOK_START -->", "<!-- WISH_BOOK_END -->", "new")
	if err == nil {
		t.Error("expected error when end marker appears before start marker")
	}
}

func TestUpdateReadmeEmptyWishFile(t *testing.T) {
	wishFile, err := os.CreateTemp("", "wish-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(wishFile.Name())
	wishFile.WriteString("[]")
	wishFile.Close()

	readmeFile, err := os.CreateTemp("", "readme-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(readmeFile.Name())
	readmeFile.WriteString("# README\n<!-- WISH_BOOK_START --><!-- WISH_BOOK_END -->\n")
	readmeFile.Close()

	if err := updateReadme(wishFile.Name(), readmeFile.Name()); err != nil {
		t.Errorf("unexpected error with empty wish list: %v", err)
	}
}

func TestUpdateReadme(t *testing.T) {
	wishFile, err := os.CreateTemp("", "wish-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(wishFile.Name())
	wishFile.WriteString(`[{"no":1,"title":"Test Book","url":"https://bookmeter.com/1","author":"Author","authorUrl":"","thumb":"https://example.com/1.jpg","date":""}]`)
	wishFile.Close()

	readmeFile, err := os.CreateTemp("", "readme-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(readmeFile.Name())
	readmeFile.WriteString("# README\n<!-- WISH_BOOK_START --><!-- WISH_BOOK_END -->\n")
	readmeFile.Close()

	if err := updateReadme(wishFile.Name(), readmeFile.Name()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(readmeFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "Test Book") {
		t.Errorf("README should contain book title, got: %s", string(content))
	}
}
