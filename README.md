# fetch-bookmeter

[![Test](../../actions/workflows/test.yaml/badge.svg)](../../actions/workflows/test.yaml)

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

### Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `user-id` | Bookmeter user ID | No | `104` |
| `output` | Output file path for the wish list JSON | No | `wish.json` |

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

      - name: Commit and push
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add wish.json
          git diff --staged --quiet || git commit -m "chore: update wish list"
          git push
```

## CLI Usage

You can also run the tool directly using Go:

```bash
# Fetch wish list
go run . fetch-wish -user-id 104 -output wish.json

# Update README with a random book from wish list
go run . update-readme -wish-file wish.json -readme README.md
```
