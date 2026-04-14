package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

type Item struct {
	Title       string
	Link        string
	Published   time.Time
	Description string
	FeedTitle   string
}

type FetchResult struct {
	Items  []Item
	Errors []error
}

func Fetch(ctx context.Context, urls []string, date time.Time) FetchResult {
	type result struct {
		items []Item
		err   error
	}

	ch := make(chan result, len(urls))
	var wg sync.WaitGroup

	targetDate := date.UTC().Truncate(24 * time.Hour)

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			items, err := fetchOne(ctx, url, targetDate)
			ch <- result{items: items, err: err}
		}(u)
	}

	wg.Wait()
	close(ch)

	var fr FetchResult
	for r := range ch {
		if r.err != nil {
			fr.Errors = append(fr.Errors, r.err)
			continue
		}
		fr.Items = append(fr.Items, r.items...)
	}

	return fr
}

func fetchOne(ctx context.Context, url string, date time.Time) ([]Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}

	var items []Item
	for _, entry := range feed.Items {
		pub := publishedTime(entry)
		if pub == nil {
			continue
		}
		if pub.UTC().Truncate(24*time.Hour) != date {
			continue
		}

		link := entry.Link
		if link == "" && len(entry.Links) > 0 {
			link = entry.Links[0]
		}

		desc := stripHTML(entry.Description)
		if desc == "" {
			desc = stripHTML(entry.Content)
		}

		items = append(items, Item{
			Title:       entry.Title,
			Link:        link,
			Published:   pub.UTC(),
			Description: desc,
			FeedTitle:   feed.Title,
		})
	}

	return items, nil
}

func publishedTime(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	}
	if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}
	return nil
}

// stripHTML removes HTML tags from s using a simple state machine.
func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// Page returns the items for the requested page and the total item count.
// Page numbers are 1-based.
func Page(items []Item, page, perPage int) ([]Item, int) {
	total := len(items)
	if total == 0 {
		return nil, 0
	}

	start := (page - 1) * perPage
	if start >= total {
		return nil, total
	}

	end := start + perPage
	if end > total {
		end = total
	}

	return items[start:end], total
}
