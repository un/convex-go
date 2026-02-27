# Convex Go Live Demo

This folder contains a real integration demo for the Go Convex client.

## What this demo validates

- Real websocket sync connection to Convex.
- `Query`, `Mutation`, `Action`, `Subscribe`, and `WatchAll` behavior.
- Realtime updates when rows change in Convex dashboard.
- Connection diagnostics (`/api/live/connection`) and event timeline (`/api/events`).

## Matching backend

Use the scaffold in `demo/convex-app`.

Default function mapping in this demo:

- `messages:list` (query)
- `messages:create` (mutation)
- `messages:ping` (action)

## Environment

Required:

- `CONVEX_DEPLOYMENT_URL` (or `CONVEX_URL`)
- `DEMO_QUERY_FUNCTION` (defaults to `messages:list` in `demo/.env`)

Optional:

- `CONVEX_AUTH_TOKEN` / `DEMO_AUTH_TOKEN`
- `DEMO_MUTATION_FUNCTION` / `DEMO_ACTION_FUNCTION`
- `DEMO_DEFAULT_ARGS` (default `{}`)
- `PORT` (default `8080`)
- `DEMO_CLIENT_ID` (default `demo-web`)

Load env and run:

```bash
set -a
source demo/.env
set +a
go run ./demo
```

Open: `http://localhost:8080`

## Manual realtime flow

1. Seed data once:

```bash
cd demo/convex-app
npx convex run messages:seedIfEmpty '{}'
```

2. Open `http://localhost:8080`.
3. Use the in-page controls:

- `Refresh messages` loads `messages:list` into the table view.
- `Insert random row` calls `messages:create` directly (no form required).
- The table uses the same live endpoint: `/api/live/query?function=messages:list&args={}`.
- Connection status is shown in-page and updates automatically (connected/connecting/reconnecting/disconnected/timeout).

4. Optional baseline query endpoint:

```text
http://localhost:8080/api/live/query?function=messages:list&args={}
```

5. Open stream in another tab:

```text
http://localhost:8080/api/live/stream?function=messages:list&args={}&timeout_seconds=300
```

6. Edit a `messages` row in Convex dashboard.
7. Confirm `event: update` appears in stream output.

## Useful endpoints

- `GET /api/debug`
- `GET /api/live/connection`
- `GET /api/events`
- `GET /api/live/query?function=messages:list&args={}`
- `GET /api/live/mutation?function=messages:create&args={}`
- `GET /api/live/action?function=messages:ping&args={}`
- `GET /api/live/subscribe?function=messages:list&args={}`
- `GET /api/live/watch?function=messages:list&args={}`
- `GET /api/live/run-suite`
- `GET /api/live/reset`

## Troubleshooting

- If query returns `value: null`, check Convex logs for `Invalid arguments provided`.
- Ensure you run a build that includes the function-args wire fix (args are array-wrapped over sync protocol).
- Call `/api/live/connection` and verify `deployment_url` and `resolved_ws_url`.
- Reset client and retry:

```text
http://localhost:8080/api/live/reset
```

## Tests

```bash
go test ./demo -v
```

Live test (`TestLiveDemoFlow`) runs only when `CONVEX_INTEGRATION=1` and required env vars are set.
