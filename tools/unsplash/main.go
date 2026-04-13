package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

var validOrientations = map[string]bool{
	"landscape": true, "portrait": true, "squarish": true,
}

var validOrderBy = map[string]bool{
	"relevant": true, "latest": true,
}

var validContentFilter = map[string]bool{
	"low": true, "high": true,
}

var validColors = map[string]bool{
	"black_and_white": true, "black": true, "white": true,
	"yellow": true, "orange": true, "red": true, "purple": true,
	"magenta": true, "green": true, "teal": true, "blue": true,
}

func main() {
	search := flag.String("search", "", "search query")
	flag.StringVar(search, "s", "", "shorthand for --search")
	apiKey := flag.String("api-key", "", "Unsplash API access key")
	flag.StringVar(apiKey, "k", "", "shorthand for --api-key")
	page := flag.Int("page", 1, "page number")
	perPage := flag.Int("per-page", 10, "results per page (max 30)")
	orderBy := flag.String("order-by", "relevant", "sort order: relevant, latest")
	collections := flag.String("collections", "", "filter by collection IDs (comma-separated)")
	contentFilter := flag.String("content-filter", "low", "content safety filter: low, high")
	color := flag.String("color", "", "filter by color: black_and_white, black, white, yellow, orange, red, purple, magenta, green, teal, blue")
	orientation := flag.String("orientation", "", "filter by orientation: landscape, portrait, squarish")
	flag.Parse()

	if *search == "" {
		fmt.Fprintln(os.Stderr, "error: --search / -s is required")
		flag.Usage()
		os.Exit(1)
	}

	if *perPage < 1 || *perPage > 30 {
		fmt.Fprintln(os.Stderr, "error: --per-page must be between 1 and 30")
		os.Exit(1)
	}

	if !validOrderBy[*orderBy] {
		fmt.Fprintln(os.Stderr, "error: --order-by must be relevant or latest")
		os.Exit(1)
	}

	if !validContentFilter[*contentFilter] {
		fmt.Fprintln(os.Stderr, "error: --content-filter must be low or high")
		os.Exit(1)
	}

	if *color != "" && !validColors[*color] {
		fmt.Fprintln(os.Stderr, "error: --color must be one of: black_and_white, black, white, yellow, orange, red, purple, magenta, green, teal, blue")
		os.Exit(1)
	}

	if *orientation != "" && !validOrientations[*orientation] {
		fmt.Fprintln(os.Stderr, "error: --orientation must be landscape, portrait, or squarish")
		os.Exit(1)
	}

	key, err := resolveAPIKey(*apiKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	c := newClient(key)
	result, err := c.searchPhotos(context.Background(), SearchParams{
		Query:         *search,
		Page:          *page,
		PerPage:       *perPage,
		OrderBy:       *orderBy,
		Collections:   *collections,
		ContentFilter: *contentFilter,
		Color:         *color,
		Orientation:   *orientation,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if len(result.Results) == 0 {
		fmt.Println("No photos found.")
		return
	}

	fmt.Printf("Found %d photos (page %d of %d)\n\n", result.Total, *page, result.TotalPages)
	for i, photo := range result.Results {
		desc := photo.Description
		if desc == "" {
			desc = photo.AltDescription
		}
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("%d. [%s] %s\n", i+1, photo.ID, desc)
		fmt.Printf("   %s\n\n", photo.Links["html"])
	}
}
