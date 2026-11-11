# Gator

Gator is a command-line RSS feed aggregator written in Go. It allows users to register, add RSS feeds, follow feeds, collect posts, and browse their latest articles.

## Requirements

Before using Gator, make sure you have installed:

- [Go](https://go.dev/doc/install) (latest version recommended)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

Install the Gator CLI using Go:

```bash
go install github.com/ayad/gator@latest
```

Make sure your Go binary directory is included in your system's `PATH`.

You can verify the installation with:

```bash
gator
```

## Configuration

Gator uses a configuration file located at:

```text
~/.gatorconfig.json
```

Create this file and add the following configuration:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace the database username, password, host, and database name with your own PostgreSQL settings.

## Database Setup

Create a PostgreSQL database named `gator`:

```bash
createdb gator
```

Install Goose, the database migration tool:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

From the project directory, run the database migrations:

```bash
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

Adjust the connection string if your PostgreSQL credentials differ.

## Usage

Gator provides several commands.

### User management

Register a new user:

```bash
gator register ayad
```

Log in as an existing user:

```bash
gator login ayad
```

List registered users:

```bash
gator users
```

### Feed management

Add an RSS feed:

```bash
gator addfeed "Hacker News RSS" "https://hnrss.org/newest"
```

List available feeds:

```bash
gator feeds
```

Follow a feed:

```bash
gator follow "https://hnrss.org/newest"
```

List feeds you follow:

```bash
gator following
```

Unfollow a feed:

```bash
gator unfollow "https://hnrss.org/newest"
```

### Collecting and browsing posts

Start collecting posts from RSS feeds every minute:

```bash
gator agg 1m
```

Browse the latest two posts:

```bash
gator browse
```

Browse the latest five posts:

```bash
gator browse 5
```

## Development

If you have cloned the source code, you can run the application during development with:

```bash
go run .
```

To build a standalone executable:

```bash
go build
```

The resulting executable can be run without using `go run`.

## License

Add a license if you intend to distribute this project publicly.
