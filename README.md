# Convex Go

Convex Go client with websocket sync transport.

## Quickstart

Set deployment URL and run the example:

```bash
CONVEX_DEPLOYMENT_URL="https://<your-deployment>.convex.cloud" go run ./examples/quickstart
```

## Live demo app

There is a full live demo in `demo/` plus a matching Convex backend scaffold in `demo/convex-app`.

- Demo docs: `demo/README.md`
- Backend scaffold docs: `demo/convex-app/README.md`

## Testing

- Unit tests: `go test ./...`
- Race detector: `go test -race ./...`
- Live integration tests (opt-in):

```bash
CONVEX_INTEGRATION=1 \
CONVEX_DEPLOYMENT_URL="https://<your-deployment>.convex.cloud" \
CONVEX_TEST_QUERY="messages:list" \
go test ./integration -v
```
