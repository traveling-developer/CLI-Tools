package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func main() {
	help := flag.Bool("help", false, "show this help message")
	flag.BoolVar(help, "h", false, "shorthand for --help")
	feeds := flag.String("feeds", "", "path to file with RSS feed URLs (one per line)")
	flag.StringVar(feeds, "f", "", "shorthand for --feeds")
	dateFlag := flag.String("date", "", "filter by date in YYYY-MM-DD format (default: today)")
	flag.StringVar(dateFlag, "d", "", "shorthand for --date")
	page := flag.Int("page", 1, "page number")
	perPage := flag.Int("per-page", 30, "items per page (max 100)")
	format := flag.String("format", "json", "output format: json, text")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: rss [options]

Fetch and display RSS/Atom feed items from a list of feed URLs.

Options:
  -f, --feeds <file>    Path to feeds file (default: ~/.config/rss/feeds)
  -d, --date <date>     Filter by date in YYYY-MM-DD format (default: today)
      --page <n>        Page number (default: 1)
      --per-page <n>    Items per page, max 100 (default: 30)
      --format <fmt>    Output format: json, text (default: json)
  -h, --help            Show this help message

Feed file format:
  One RSS/Atom feed URL per line.
  Empty lines and lines starting with # are ignored.

  Example:
    # Tech blogs
    https://go.dev/blog/feed.atom
    https://blog.example.org/rss.xml

Examples:
  rss
  rss --date 2026-04-14
  rss -f feeds.txt --per-page 10
  rss -f feeds.txt --page 2
  rss --format text
`)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *format != "json" && *format != "text" {
		fmt.Fprintln(os.Stderr, "error: --format must be json or text")
		os.Exit(1)
	}

	if *perPage < 1 || *perPage > 100 {
		fmt.Fprintln(os.Stderr, "error: --per-page must be between 1 and 100")
		os.Exit(1)
	}

	if *page < 1 {
		fmt.Fprintln(os.Stderr, "error: --page must be >= 1")
		os.Exit(1)
	}

	feedsFile, err := resolveFeedsFile(*feeds)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	urls, err := readFeedURLs(feedsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	date := time.Now().UTC().Truncate(24 * time.Hour)
	if *dateFlag != "" {
		parsed, err := time.Parse("2006-01-02", *dateFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: --date must be in YYYY-MM-DD format")
			os.Exit(1)
		}
		date = parsed
	}

	result := Fetch(context.Background(), urls, date)

	items, total := Page(result.Items, *page, *perPage)
	totalPages := int(math.Ceil(float64(total) / float64(*perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	warnings := make([]string, 0, len(result.Errors))
	for _, e := range result.Errors {
		warnings = append(warnings, e.Error())
	}

	if *format == "json" {
		printJSON(date, *page, totalPages, total, items, warnings)
	} else {
		printText(date, *page, totalPages, total, items, warnings, feedsFile, *perPage)
	}
}

type jsonItem struct {
	Title       string `json:"title"`
	Feed        string `json:"feed"`
	Published   string `json:"published"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

type jsonOutput struct {
	Date       string     `json:"date"`
	Page       int        `json:"page"`
	TotalPages int        `json:"total_pages"`
	Total      int        `json:"total"`
	Items      []jsonItem `json:"items"`
	Warnings   []string   `json:"warnings"`
}

func printJSON(date time.Time, page, totalPages, total int, items []Item, warnings []string) {
	out := jsonOutput{
		Date:       date.Format("2006-01-02"),
		Page:       page,
		TotalPages: totalPages,
		Total:      total,
		Items:      make([]jsonItem, 0, len(items)),
		Warnings:   warnings,
	}

	for _, item := range items {
		out.Items = append(out.Items, jsonItem{
			Title:       item.Title,
			Feed:        item.FeedTitle,
			Published:   item.Published.Format(time.RFC3339),
			Link:        item.Link,
			Description: item.Description,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}

func printText(date time.Time, page, totalPages, total int, items []Item, warnings []string, feedsFile string, perPage int) {
	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, "warning:", w)
	}

	if total == 0 {
		fmt.Printf("No items found for %s.\n", date.Format("2006-01-02"))
		return
	}

	fmt.Printf("Found %d items for %s (page %d of %d)\n\n",
		total, date.Format("2006-01-02"), page, totalPages)

	for i, item := range items {
		n := (page-1)*perPage + i + 1
		fmt.Printf("%d. %s\n", n, item.Title)
		if item.FeedTitle != "" {
			fmt.Printf("   Feed: %s\n", item.FeedTitle)
		}
		fmt.Printf("   %s\n", item.Published.Format("2006-01-02 15:04"))
		if item.Link != "" {
			fmt.Printf("   %s\n", item.Link)
		}
		if item.Description != "" {
			desc := strings.TrimSpace(item.Description)
			if len(desc) > 120 {
				desc = desc[:117] + "..."
			}
			fmt.Printf("   %s\n", desc)
		}
		fmt.Println()
	}

	if page < totalPages {
		fmt.Printf("Next page: rss -f %s --date %s --page %d\n",
			feedsFile, date.Format("2006-01-02"), page+1)
	}
}
