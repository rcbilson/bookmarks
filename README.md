# bookmarks

Richard's quirky bookmarks manager.

## What is it?

A simple web app that maintains a database of Internet bookmarks. You add a URL,
it remembers it.

To be very clear, no one needs this. I'm just fed up with Internet browsers
burying bookmarks, which I still find very useful. Really I should just sign up
for [Pinboard](https://pinboard.in/) but I already had written most of the
relevant code as part of my [recipe manager](https://github.com/rcbilson/recipe).

## Building and running

`make docker` will build a container. I use a docker-compose fragment something
like this to run it:

```
  bookmarks:
    image: rcbilson/bookmarks:latest
    pull_policy: never
    ports:
      - 80:9093
    volumes:
      - ./bookmarks/data:/app/data
    restart: unless-stopped
```

## What's under the hood

The frontend is Vite + TypeScript + React with some chakra-ui. The backend is
Go + Sqlite.

## Known issues

The server attempts to obtain the title of the web page, which is the primary
thing displayed in the UI to identify each site. In our benighted age there are
no small number of websites that don't actually have a `<title>` element
defined in their static HTML. For this reason you may not actually get a title
for every bookmark, so all that will be displayed in the UI is the domain in
small print. Too bad!
