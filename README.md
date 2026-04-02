# Go HTTP Service

A small Go HTTP service that can be configured with environment variables, logs requests, handles graceful shutdown, and can be run in a Docker container.


## How To Build

Build the Docker image (example tag `go-http-service`):

```bash
docker build -t go-http-service .
```

## How To Run

Quick start (no config):

```bash
docker run --rm -p 9000:8000 go-http-service
```

Run with environment variables:

```bash
docker run --rm -p 9000:8000 \
  -e RESPONSE_MESSAGE="Server is up!" \
  -e ALLOW_ORIGIN="*" \
  -e PORT=8000 \
  --name cybernetica-devops-server \
  go-http-service
```

## Supported environment variables

`PORT` - Port the server listens to inside the container 
`RESPONSE_MESSAGE` - Text returned by `GET /` 
`ALLOW_ORIGIN` -  Value for `Access-Control-Allow-Origin` header 

## Test endpoints (open in your browser)

Open these URLs in a browser (replace `9000` with the host port you published):

- `http://localhost:9000/` — expected body: `Service request succeeded!` (or the value of `RESPONSE_MESSAGE` if overridden)
- `http://localhost:9000/health` — expected: `ok`
- `http://localhost:9000/ready` — expected: `ready`


