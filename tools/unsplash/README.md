# unsplash

A CLI tool to search for photos on [Unsplash](https://unsplash.com) directly from your terminal.

## Requirements

An Unsplash API access key. Register your application at [unsplash.com/developers](https://unsplash.com/developers) to get one.

## Authentication

The API key is resolved in the following order — the first match wins:

1. `--api-key` / `-k` flag
2. `UNSPLASH_ACCESS_KEY` environment variable
3. `~/.config/unsplash/config.json`

```json
{
  "access_key": "your-access-key"
}
```

## Usage

```
unsplash --search <query> [options]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--search` | `-s` | — | Search query **(required)** |
| `--api-key` | `-k` | — | Unsplash API access key |
| `--page` | | `1` | Page number |
| `--per-page` | | `10` | Results per page (max 30) |
| `--order-by` | | `relevant` | Sort order: `relevant`, `latest` |
| `--orientation` | | — | Filter by orientation: `landscape`, `portrait`, `squarish` |
| `--color` | | — | Filter by color: `black_and_white`, `black`, `white`, `yellow`, `orange`, `red`, `purple`, `magenta`, `green`, `teal`, `blue` |
| `--content-filter` | | `low` | Content safety level: `low`, `high` |
| `--collections` | | — | Limit results to comma-separated collection IDs |

### Examples

```bash
# Simple search
unsplash -s "mountains"

# Latest photos of cats, landscape only
unsplash -s "cats" --order-by latest --orientation landscape

# High content filter, blue tones, page 2
unsplash -s "ocean" --content-filter high --color blue --page 2 --per-page 20

# Limit to specific collections
unsplash -s "architecture" --collections "123456,789012"
```

## Output

Each result shows the photo ID, description, and a direct link to the Unsplash page:

```
Found 1042 photos (page 1 of 35)

1. [abc123] Sunrise over the Alps
   https://unsplash.com/photos/abc123

2. [def456] (no description)
   https://unsplash.com/photos/def456
```
