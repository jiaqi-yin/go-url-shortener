# Golang URL Shortener

This application provides the service to generate short links of original long URLs and can redirect to the original URL from the short version.

## Get Started

### Build and Run

Run `docker-compose up --build`

### Create a short link

```
Request:

curl -X POST http://localhost:8080/api/shorten -d '{"url":"https://www.google.com"}'
```

```
Response: 

{"shortlink":"4c94"}
```

### Redirect from short links

```
Request:

curl -sL -o /dev/null -D - http://localhost:8080/4c94
```

```
Response:

HTTP/1.1 307 Temporary Redirect
Location: https://www.google.com
...
HTTP/2 200
...
```