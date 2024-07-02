# go-urlshortner
A simple golang based url shortner


# How to use

```bash
## Shorten

curl --location 'localhost:9090/shorten' \
--header 'Content-Type: application/json' \
--data '{
    "shortcode" : "something",
    "url" : "github.com/labstack/echo"
}'

## Expand

curl --location 'localhost:9090/resolve/H2eBb6CN'

```
