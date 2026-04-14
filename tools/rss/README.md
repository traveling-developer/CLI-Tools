# rss

A CLI tool to fetch and display RSS/Atom feed items from a list of feed URLs, with date filtering and pagination.

## Usage

```
rss [options]
```

### Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--feeds <file>` | `-f` | `~/.config/rss/feeds` | Path to feeds file (see below) |
| `--date <date>` | `-d` | today | Filter by date in `YYYY-MM-DD` format |
| `--page <n>` | | `1` | Page number |
| `--per-page <n>` | | `30` | Items per page (max 100) |
| `--format <fmt>` | | `json` | Output format: `json`, `text` |
| `--help` | `-h` | | Show help message |

### Examples

```bash
# Fetch today's items (JSON output)
rss

# Fetch items for a specific date
rss --date 2026-04-14

# Human-readable output
rss --format text

# Paginate through results
rss -f feeds.txt --page 2 --per-page 10
```

## Feeds File

The feeds file is a plain text file with one RSS or Atom feed URL per line.

- Empty lines are ignored
- Lines starting with `#` are treated as comments and ignored

**Example `feeds.txt`:**

```
# Tech blogs
https://go.dev/blog/feed.atom
https://blog.example.org/rss.xml

# News
https://news.example.com/feed
```

## Configuration

If `--feeds` is not provided, `rss` automatically looks for a feeds file at:

```
~/.config/rss/feeds
```

This allows you to set up your feeds once and run `rss` without any arguments:

```bash
mkdir -p ~/.config/rss
cat > ~/.config/rss/feeds <<EOF
https://go.dev/blog/feed.atom
https://blog.example.org/rss.xml
EOF

rss
```

## Output

The default output format is JSON, designed for use by AI agents and scripts.

```json
{
  "date": "2026-04-14",
  "page": 1,
  "total_pages": 3,
  "total": 42,
  "items": [
    {
      "title": "Example Post",
      "feed": "Example Blog",
      "published": "2026-04-14T08:35:00Z",
      "link": "https://example.com/posts/example",
      "description": "A short summary of the post."
    }
  ],
  "warnings": []
}
```

Use `--format text` for human-readable output.