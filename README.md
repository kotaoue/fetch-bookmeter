# fetch-bookmeter

[![Test](../../actions/workflows/test.yaml/badge.svg)](../../actions/workflows/test.yaml)

Fetch from Bookmeter (読書メーター).

A GitHub Action that fetches book data from [Bookmeter](https://bookmeter.com/) and saves it as JSON.

## Usage

### Fetch wish list

```yaml
- name: Fetch Bookmeter wish list
  uses: kotaoue/fetch-bookmeter@v1
  with:
    user-id: '104'
    output: 'wish.json'
```

### Fetch read books history

```yaml
- name: Fetch Bookmeter read books history
  uses: kotaoue/fetch-bookmeter@v1
  with:
    user-id: '104'
    type: 'read'
    output: 'read.json'
```

To filter by a specific year and/or month, pass the additional inputs:

```yaml
- name: Fetch Bookmeter read books (2024 January)
  uses: kotaoue/fetch-bookmeter@v1
  with:
    user-id: '104'
    type: 'read'
    output: 'read-2024-01.json'
    year: '2024'
    month: '1'
```

### Inputs

| Input | Description | Required | Default |
| ----- | ----------- | -------- | ------- |
| `user-id` | Bookmeter user ID | No | `104` |
| `output` | Output file path for the JSON | No | `wish.json` |
| `type` | Type of book list to fetch: `wish` for wish list, `read` for read books history | No | `wish` |
| `year` | Filter read books by year (e.g. `2024`). Only used when `type` is `read`. `0` means no filter. | No | `0` |
| `month` | Filter read books by month (1-12). Only used when `type` is `read`. `0` means no filter. | No | `0` |

### Full workflow example

```yaml
name: Update README

on:
  schedule:
    - cron: '0 15 * * *'
  workflow_dispatch:

permissions:
  contents: write

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Fetch Bookmeter wish list
        uses: kotaoue/fetch-bookmeter@v1
        with:
          user-id: '104'
          output: ${{ github.workspace }}/wish.json

      - name: Fetch Bookmeter read books history
        uses: kotaoue/fetch-bookmeter@v1
        with:
          user-id: '104'
          type: 'read'
          output: ${{ github.workspace }}/read.json

      - name: Commit and push
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add wish.json read.json
          git diff --staged --quiet || git commit -m "chore: update book lists"
          git push
```

## CLI Usage

You can also run the tool directly using Go:

```bash
# Fetch wish list
go run . fetch-wish -user-id 104 -output wish.json

# Fetch read books history (all)
go run . fetch-read -user-id 104 -output read.json

# Fetch read books history filtered by year
go run . fetch-read -user-id 104 -year 2024 -output read-2024.json

# Fetch read books history filtered by year and month
go run . fetch-read -user-id 104 -year 2024 -month 1 -output read-2024-01.json
```
