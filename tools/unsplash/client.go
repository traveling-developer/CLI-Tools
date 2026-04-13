package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const baseURL = "https://api.unsplash.com"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func newClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type SearchParams struct {
	Query         string
	Page          int
	PerPage       int
	OrderBy       string
	Collections   string
	ContentFilter string
	Color         string
	Orientation   string
}

type Photo struct {
	ID             string            `json:"id"`
	Description    string            `json:"description"`
	AltDescription string            `json:"alt_description"`
	Links          map[string]string `json:"links"`
}

type SearchResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}

func (c *Client) searchPhotos(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	q := url.Values{}
	q.Set("query", params.Query)
	q.Set("page", strconv.Itoa(params.Page))
	q.Set("per_page", strconv.Itoa(params.PerPage))
	q.Set("order_by", params.OrderBy)
	if params.Collections != "" {
		q.Set("collections", params.Collections)
	}
	if params.ContentFilter != "" {
		q.Set("content_filter", params.ContentFilter)
	}
	if params.Color != "" {
		q.Set("color", params.Color)
	}
	if params.Orientation != "" {
		q.Set("orientation", params.Orientation)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/search/photos?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Client-ID "+c.apiKey)
	req.Header.Set("Accept-Version", "v1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, body)
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}
