# RSS Feed Aggregator in GO

## Prerequisits
- Go 1.23+
- a running instance of postgres (docker compose file included)

## Setup 

### database
Launch postgres database - e.g. in docker using `docker compose up` within the project directory.
Services configured in `docker-compose.yaml`

Connection string
```
postgres://postgres:postgres@localhost:5432/gator?sslmode=disable
```

### GO packages
Run `go install` within the project root directory

### Configuration files
Copy the example configuration file `config.example.json` into `~/.config/gator/config.json` and update the database connect string as needed.

If you are running postgres in docker using the docker compose file the connect string above is already included in the example.

### Deploy database schema
With a running database push the database schema into the database using goose. 
Therefore install goose and from within `sql/schema` run 
```
goose postgres <connect-string> up
```
Use your connect string - e.g. `postgres://postgres:postgres@localhost:5432/gator`

Now you are ready to launch the CLI. Assuming you have created an executable - e.g. by `go build -o gator .`.
## register
Registers a new user with given username. The new user gets logged in immediately.
```
gator register <username>
```

## login
Logs in a registered user.
```
gator login <username>
```

## users
List all registered users - including the current logged in user.
```
gator users
```

## addfeed
Adds a new rss feed to watch.
```
gator addfeed <name> <url>
```
e.g.
```
gator addfeed "Hacker News RSS" "https://hnrss.org/newest"
```

## feeds
List all watched feeds.
```
gator feeds
```

## agg
Launch a scraper to check for newest updates in the stored feeds. 
The scraper always takes the least recently scraped feed first.
Specify an optional frequency between scraping the next feed.
```
gator agg [frequency]
```
Optional frequency to be specified in GO duration format - e.g. 2m or 1h30m.
Defaults to 30s. Valid time units are "ms", "s", "m", "h".
Note that gator will not accept scraping faster than every 5 seconds (5000ms).
Stop scraping by hitting Ctrl-C.

## follow
Registers current logged in user as follower of the specified feed (by feed URL).
You need to follow feeds to get them reported in `browse`.
```
gator follow <feed url>
```
Note that you can only follow feeds which are already added to gator. See `addfeed` command.

## unfollow
Stop following the specified feed.
```
gator unfollow <feed url>
```

## following
List all feeds you are following.
```
gator following
```

## browse
List the newset posts across all feeds you are following. 
Specify optional `limit` to show more or less than 2 posts. 
```
gator browse [limit]
```