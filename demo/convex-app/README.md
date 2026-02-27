# Convex Backend Scaffold For Go Demo

This is a minimal Convex backend designed to match the Go demo in `demo/`.

It includes:

- Query: `messages:list`
- Mutations: `messages:create`, `messages:updateText`, `messages:seedIfEmpty`
- Action: `messages:ping`
- Internal functions used by action: `messages_internal:summary`, `messages_internal:insertActionLog`, `messages_internal:echoAction`
- HTTP actions: `GET /health`, `POST /seed`

## 1) Install and run locally

```bash
cd demo/convex-app
npm install
npx convex dev
```

`convex dev` will prompt login/project setup.

If it does not print the deployment URL directly, get it from the Convex dashboard for that project and set it in `demo/.env` as `CONVEX_DEPLOYMENT_URL`.

## 2) Connect the Go demo app

Use the deployment URL from Convex and set in `demo/.env`:

- `CONVEX_DEPLOYMENT_URL`
- `DEMO_QUERY_FUNCTION=messages:list`
- `DEMO_MUTATION_FUNCTION=messages:create`
- `DEMO_ACTION_FUNCTION=messages:ping`

Then run the Go demo from repo root:

```bash
set -a
source demo/.env
set +a
go run ./demo
```

The Go demo homepage includes an in-page live table and an `Insert random row` button that calls `messages:create`.

## 3) Realtime manual test

1. Call `messages:seedIfEmpty` once from Convex dashboard (or use `POST /seed` HTTP action).
2. Open stream endpoint from Go demo:

```text
http://localhost:8080/api/live/stream?function=messages:list&args={}
```

3. In Convex dashboard, edit a row in `messages` table (text/source/updatedAt).
4. The stream should emit `event: update` with updated query values.

## 4) Optional production deploy

After local validation, deploy with:

```bash
npx convex deploy
```

Then update `CONVEX_DEPLOYMENT_URL` in `demo/.env` to the production deployment URL.

## 5) Validate function wiring

From `demo/convex-app`:

```bash
npx convex run messages:seedIfEmpty '{}'
npx convex run messages:list '{}'
npx convex run messages:ping '{}'
```

If these succeed, the Go demo should work with:

```text
http://localhost:8080/api/live/query?function=messages:list&args={}
```
