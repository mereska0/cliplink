# ClipLink

ClipLink is a small URL shortener built with Go.

It provides a terminal UI, gRPC API, HTTP redirects, PostgreSQL storage, Docker Compose setup, and unit tests. The project was made as a backend pet project to practice real service architecture in Go.


## Tech Stack

* Go
* gRPC / Protocol Buffers
* PostgreSQL
* Docker Compose
* Bubble Tea for TUI
* pgx for PostgreSQL
* Standard Go testing package

## How It Works

When a user creates a short link, cliplink saves the original URL in PostgreSQL and receives a unique numeric ID.

This ID is encoded into a Base62 string and stored as the short code.

Example:

```text
ID: 1025 -> short code: gw
```

When someone opens:

```text
http://localhost:8080/gw
```

## Run with Docker Compose

Start PostgreSQL:

```bash
docker compose up -d postgres
```

Or run the default mode:

```bash
go run ./cmd/cliplink
```


## Status

MVP is implemented:

* Shorten link
* Redirect
* List links
* Delete link
* Click counter
* gRPC
* TUI
* PostgreSQL
* Docker Compose
* Unit tests

