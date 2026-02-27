package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/get-convex/convex-go/convex"
	convsync "github.com/get-convex/convex-go/internal/sync"
)

type appConfig struct {
	Port              string
	ClientID          string
	DeploymentURL     string
	AuthToken         string
	DefaultQueryFn    string
	DefaultMutationFn string
	DefaultActionFn   string
	DefaultArgsRaw    string
}

type debugEvent struct {
	ID         int64  `json:"id"`
	Timestamp  string `json:"timestamp"`
	Endpoint   string `json:"endpoint"`
	Success    bool   `json:"success"`
	DurationMS int64  `json:"duration_ms"`
	Message    string `json:"message,omitempty"`
}

type suiteCheckResult struct {
	Name       string         `json:"name"`
	Passed     bool           `json:"passed"`
	Skipped    bool           `json:"skipped"`
	Error      string         `json:"error,omitempty"`
	DurationMS int64          `json:"duration_ms"`
	Data       map[string]any `json:"data,omitempty"`
}

type demoApp struct {
	mu          sync.Mutex
	cfg         appConfig
	startedAt   time.Time
	defaultArgs map[string]any
	configErr   string
	client      *convex.Client
	events      []debugEvent
	nextID      int64
}

type homeView struct {
	Port              string
	ClientID          string
	StartedAt         string
	Uptime            string
	DeploymentURL     string
	HasAuthToken      bool
	Ready             bool
	ConfigError       string
	DefaultQueryFn    string
	DefaultMutationFn string
	DefaultActionFn   string
	DefaultArgsRaw    string
	EventCount        int
}

var homeTmpl = template.Must(template.New("home").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Convex Go Live Demo</title>
  <style>
    :root {
      --bg: #f6f0e6;
      --card: #fffdf8;
      --ink: #2f211a;
      --muted: #695244;
      --line: #ddc7b8;
      --accent: #9f3d2f;
      --accent-soft: #f8e3d7;
    }
    body {
      margin: 0;
      font-family: "Avenir Next", "Segoe UI", sans-serif;
      color: var(--ink);
      background:
        radial-gradient(circle at 0% 0%, #f4d8c5 0%, transparent 35%),
        radial-gradient(circle at 100% 0%, #e8d5f0 0%, transparent 28%),
        var(--bg);
    }
    main {
      max-width: 1024px;
      margin: 0 auto;
      padding: 24px;
      display: grid;
      gap: 16px;
    }
    .card {
      background: var(--card);
      border: 1px solid var(--line);
      border-radius: 14px;
      padding: 18px;
      box-shadow: 0 8px 24px rgba(47, 20, 10, 0.06);
    }
    h1, h2 {
      margin: 0 0 10px;
    }
    p, li {
      line-height: 1.5;
    }
    code {
      background: var(--accent-soft);
      border: 1px solid #ebc1ac;
      border-radius: 6px;
      padding: 0.1rem 0.35rem;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 12px;
    }
    .kvs {
      display: grid;
      grid-template-columns: 180px 1fr;
      gap: 8px;
      margin: 0;
    }
    .kvs dt {
      color: var(--muted);
      font-weight: 600;
    }
    .kvs dd {
      margin: 0;
      word-break: break-word;
    }
    a {
      color: var(--accent);
      text-decoration: none;
      font-weight: 600;
    }
    a:hover {
      text-decoration: underline;
    }
    .links {
      display: grid;
      gap: 8px;
    }
    .actions {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin-bottom: 12px;
    }
    .btn {
      border: 1px solid var(--line);
      background: var(--accent-soft);
      color: var(--ink);
      border-radius: 10px;
      padding: 8px 12px;
      font-weight: 600;
      cursor: pointer;
    }
    .btn:hover {
      border-color: #d7a88f;
      background: #f6d5c6;
    }
    .status {
      margin: 0 0 10px;
      color: var(--muted);
      font-size: 0.95rem;
    }
    .status.ok {
      color: #176b2c;
    }
    .status.warn {
      color: #8a5a00;
    }
    .status.error {
      color: #b10000;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 0.95rem;
    }
    th, td {
      border-bottom: 1px solid var(--line);
      padding: 8px;
      text-align: left;
      vertical-align: top;
      word-break: break-word;
    }
    th {
      color: var(--muted);
      font-weight: 700;
    }
    pre {
      margin: 0;
      padding: 12px;
      background: #faf4ec;
      border: 1px solid var(--line);
      border-radius: 10px;
      overflow: auto;
      max-height: 280px;
    }
    @media (max-width: 640px) {
      .kvs {
        grid-template-columns: 1fr;
      }
    }
  </style>
</head>
<body>
  <main>
    <section class="card">
      <h1>Convex Go Live Demo</h1>
      <p>This demo uses the real <code>convex.Client</code> (websocket sync, query, mutation, action, subscribe, watch). No mocked transport paths are used in the demo flow.</p>
    </section>

    <section class="card">
      <h2>Runtime</h2>
      <dl class="kvs">
        <dt>Port</dt><dd>{{.Port}}</dd>
        <dt>Client ID</dt><dd>{{.ClientID}}</dd>
        <dt>Started At</dt><dd>{{.StartedAt}}</dd>
        <dt>Uptime</dt><dd>{{.Uptime}}</dd>
        <dt>Deployment URL</dt><dd>{{if .DeploymentURL}}{{.DeploymentURL}}{{else}}(not set){{end}}</dd>
        <dt>Auth Token</dt><dd>{{if .HasAuthToken}}set{{else}}not set{{end}}</dd>
        <dt>Client Ready</dt><dd>{{if .Ready}}yes{{else}}no{{end}}</dd>
        <dt>Config Error</dt><dd>{{if .ConfigError}}{{.ConfigError}}{{else}}(none){{end}}</dd>
        <dt>Default Query</dt><dd>{{if .DefaultQueryFn}}{{.DefaultQueryFn}}{{else}}(required for convenience){{end}}</dd>
        <dt>Default Mutation</dt><dd>{{if .DefaultMutationFn}}{{.DefaultMutationFn}}{{else}}(optional){{end}}</dd>
        <dt>Default Action</dt><dd>{{if .DefaultActionFn}}{{.DefaultActionFn}}{{else}}(optional){{end}}</dd>
        <dt>Default Args</dt><dd><code>{{.DefaultArgsRaw}}</code></dd>
        <dt>Debug Events</dt><dd>{{.EventCount}}</dd>
      </dl>
    </section>

    <section class="card">
      <h2>Live Data Panel</h2>
      <p class="status" id="statusText">Loading live data...</p>
      <p class="status" id="connectionStatusText">Connection: checking...</p>
      <div class="actions">
        <button class="btn" id="btnRefresh" type="button">Refresh messages</button>
        <button class="btn" id="btnInsert" type="button">Insert random row</button>
        <button class="btn" id="btnConnection" type="button">Refresh connection</button>
        <a href="/api/live/stream?function=messages:list&amp;args={}&amp;timeout_seconds=300" target="_blank" rel="noopener noreferrer"><code>Open live stream</code></a>
      </div>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Text</th>
            <th>Source</th>
            <th>Updated At</th>
            <th>Created At</th>
          </tr>
        </thead>
        <tbody id="messagesBody">
          <tr><td colspan="5">Loading...</td></tr>
        </tbody>
      </table>
    </section>

    <section class="card grid">
      <div>
        <h2>Connection Snapshot</h2>
        <pre id="connectionOutput">Loading...</pre>
      </div>
      <div>
        <h2>Last API Payload</h2>
        <pre id="payloadOutput">No request yet.</pre>
      </div>
    </section>

    <section class="card">
      <h2>Raw Endpoints</h2>
      <div class="links">
        <a href="/api/debug" target="_blank" rel="noopener noreferrer"><code>/api/debug</code></a>
        <a href="/api/live/query?function=messages:list&amp;args={}" target="_blank" rel="noopener noreferrer"><code>/api/live/query?function=messages:list&amp;args={}</code></a>
        <a href="/api/live/mutation?function=messages:create&amp;args={}" target="_blank" rel="noopener noreferrer"><code>/api/live/mutation?function=messages:create&amp;args={}</code></a>
        <a href="/api/live/action?function=messages:ping&amp;args={}" target="_blank" rel="noopener noreferrer"><code>/api/live/action?function=messages:ping&amp;args={}</code></a>
        <a href="/api/live/subscribe?function=messages:list&amp;args={}" target="_blank" rel="noopener noreferrer"><code>/api/live/subscribe?function=messages:list&amp;args={}</code></a>
        <a href="/api/live/watch?function=messages:list&amp;args={}" target="_blank" rel="noopener noreferrer"><code>/api/live/watch?function=messages:list&amp;args={}</code></a>
        <a href="/api/live/connection" target="_blank" rel="noopener noreferrer"><code>/api/live/connection</code></a>
        <a href="/api/live/run-suite" target="_blank" rel="noopener noreferrer"><code>/api/live/run-suite</code></a>
        <a href="/api/live/reset" target="_blank" rel="noopener noreferrer"><code>/api/live/reset</code></a>
        <a href="/api/events" target="_blank" rel="noopener noreferrer"><code>/api/events</code></a>
      </div>
    </section>
  </main>
  <script>
    (function () {
      const statusText = document.getElementById("statusText");
      const connectionStatusText = document.getElementById("connectionStatusText");
      const messagesBody = document.getElementById("messagesBody");
      const payloadOutput = document.getElementById("payloadOutput");
      const connectionOutput = document.getElementById("connectionOutput");
      const btnRefresh = document.getElementById("btnRefresh");
      const btnInsert = document.getElementById("btnInsert");
      const btnConnection = document.getElementById("btnConnection");

      function setStatus(text, isError) {
        statusText.textContent = text;
        statusText.className = isError ? "status error" : "status";
      }

      function setConnectionStatus(kind, text) {
        connectionStatusText.textContent = text;
        if (kind === "ok") {
          connectionStatusText.className = "status ok";
          return;
        }
        if (kind === "warn") {
          connectionStatusText.className = "status warn";
          return;
        }
        if (kind === "error") {
          connectionStatusText.className = "status error";
          return;
        }
        connectionStatusText.className = "status";
      }

      function toJSON(value) {
        try {
          return JSON.stringify(value, null, 2);
        } catch (_err) {
          return String(value);
        }
      }

      function formatTimestamp(value) {
        if (typeof value !== "number") {
          return "-";
        }
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) {
          return String(value);
        }
        return date.toLocaleString();
      }

      function renderMessages(rows) {
        messagesBody.innerHTML = "";
        if (!Array.isArray(rows) || rows.length === 0) {
          const tr = document.createElement("tr");
          const td = document.createElement("td");
          td.colSpan = 5;
          td.textContent = "No rows found.";
          tr.appendChild(td);
          messagesBody.appendChild(tr);
          return;
        }

        for (const row of rows) {
          const tr = document.createElement("tr");

          const id = document.createElement("td");
          id.textContent = row && row._id ? String(row._id) : "-";
          tr.appendChild(id);

          const text = document.createElement("td");
          text.textContent = row && row.text ? String(row.text) : "";
          tr.appendChild(text);

          const source = document.createElement("td");
          source.textContent = row && row.source ? String(row.source) : "";
          tr.appendChild(source);

          const updatedAt = document.createElement("td");
          updatedAt.textContent = formatTimestamp(row ? row.updatedAt : undefined);
          tr.appendChild(updatedAt);

          const creation = document.createElement("td");
          creation.textContent = formatTimestamp(row ? row._creationTime : undefined);
          tr.appendChild(creation);

          messagesBody.appendChild(tr);
        }
      }

      async function callJSON(url, init) {
        const response = await fetch(url, init || {});
        const text = await response.text();
        let parsed;
        try {
          parsed = JSON.parse(text);
        } catch (_err) {
          throw new Error("Non-JSON response: " + text);
        }
        if (!response.ok || parsed.ok === false) {
          throw new Error(parsed.error || ("HTTP " + response.status));
        }
        return parsed;
      }

      async function refreshMessages() {
        setStatus("Loading messages...", false);
        const payload = await callJSON("/api/live/query?function=messages:list&args={}&timeout_ms=12000");
        payloadOutput.textContent = toJSON(payload);
        const rows = payload && payload.data ? payload.data.value : [];
        if (rows === null) {
          renderMessages([]);
          setStatus("Query returned null. Check Convex logs for argument or function errors.", true);
          return;
        }
        renderMessages(rows);
        const count = Array.isArray(rows) ? rows.length : 0;
        setStatus("Loaded " + count + " row(s).", false);
      }

      async function refreshConnection() {
        const payload = await callJSON("/api/live/connection");
        connectionOutput.textContent = toJSON(payload.data);
        const data = payload.data || {};
        const state = typeof data.last_ws_state === "string" ? data.last_ws_state : "";
        if (!state) {
          setConnectionStatus("warn", "Connection: no websocket state seen yet");
          return;
        }
        if (state === "connected") {
          setConnectionStatus("ok", "Connection: connected");
          return;
        }
        if (state === "connecting" || state === "reconnecting") {
          setConnectionStatus("warn", "Connection: " + state);
          return;
        }
        setConnectionStatus("error", "Connection: " + state);
      }

      async function insertRandomRow() {
        setStatus("Inserting random row...", false);
        const payload = await callJSON("/api/live/insert", { method: "POST" });
        payloadOutput.textContent = toJSON(payload);
        await refreshMessages();
      }

      async function runSafe(fn) {
        try {
          await fn();
        } catch (err) {
          const message = String(err && err.message ? err.message : err);
          if (message.indexOf("context deadline exceeded") !== -1 || message.indexOf("timeout") !== -1) {
            setConnectionStatus("error", "Connection: request timed out");
          }
          setStatus(message, true);
        }
      }

      btnRefresh.addEventListener("click", function () { runSafe(refreshMessages); });
      btnInsert.addEventListener("click", function () { runSafe(insertRandomRow); });
      btnConnection.addEventListener("click", function () { runSafe(refreshConnection); });

      runSafe(async function () {
        await refreshConnection();
        await refreshMessages();
      });

      setInterval(function () {
        runSafe(refreshConnection);
      }, 5000);
    })();
  </script>
</body>
</html>
`))

func main() {
	cfg := loadConfig()
	app := newDemoApp(cfg)
	defer app.close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.serveHome)
	mux.HandleFunc("/api/debug", app.jsonHandler("debug", app.handleDebug))
	mux.HandleFunc("/api/events", app.jsonHandler("events", app.handleEvents))
	mux.HandleFunc("/api/live/connection", app.jsonHandler("live_connection", app.handleLiveConnection))
	mux.HandleFunc("/api/live/query", app.jsonHandler("live_query", app.handleLiveQuery))
	mux.HandleFunc("/api/live/insert", app.jsonHandler("live_insert", app.handleLiveInsert))
	mux.HandleFunc("/api/live/mutation", app.jsonHandler("live_mutation", app.handleLiveMutation))
	mux.HandleFunc("/api/live/action", app.jsonHandler("live_action", app.handleLiveAction))
	mux.HandleFunc("/api/live/subscribe", app.jsonHandler("live_subscribe", app.handleLiveSubscribe))
	mux.HandleFunc("/api/live/watch", app.jsonHandler("live_watch", app.handleLiveWatch))
	mux.HandleFunc("/api/live/run-suite", app.jsonHandler("live_run_suite", app.handleLiveRunSuite))
	mux.HandleFunc("/api/live/reset", app.jsonHandler("live_reset", app.handleLiveReset))
	mux.HandleFunc("/api/live/stream", app.handleLiveStream)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           requestLogMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("demo server listening on http://localhost:%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func loadConfig() appConfig {
	deployment := envOrDefault("CONVEX_DEPLOYMENT_URL", os.Getenv("CONVEX_URL"))
	return appConfig{
		Port:              envOrDefault("PORT", "8080"),
		ClientID:          envOrDefault("DEMO_CLIENT_ID", "demo-web"),
		DeploymentURL:     deployment,
		AuthToken:         envOrDefault("DEMO_AUTH_TOKEN", os.Getenv("CONVEX_AUTH_TOKEN")),
		DefaultQueryFn:    envOrDefault("DEMO_QUERY_FUNCTION", os.Getenv("CONVEX_TEST_QUERY")),
		DefaultMutationFn: envOrDefault("DEMO_MUTATION_FUNCTION", os.Getenv("CONVEX_TEST_MUTATION")),
		DefaultActionFn:   envOrDefault("DEMO_ACTION_FUNCTION", os.Getenv("CONVEX_TEST_ACTION")),
		DefaultArgsRaw:    envOrDefault("DEMO_DEFAULT_ARGS", "{}"),
	}
}

func newDemoApp(cfg appConfig) *demoApp {
	app := &demoApp{
		cfg:       cfg,
		startedAt: time.Now(),
		events:    make([]debugEvent, 0, 256),
	}
	args, err := parseArgs(cfg.DefaultArgsRaw)
	if err != nil {
		app.defaultArgs = map[string]any{}
		app.configErr = fmt.Sprintf("invalid DEMO_DEFAULT_ARGS: %v", err)
	} else {
		app.defaultArgs = args
	}
	app.resetClient("startup")
	return app
}

func (a *demoApp) close() {
	a.mu.Lock()
	client := a.client
	a.mu.Unlock()
	if client != nil {
		client.Close()
	}
}

func (a *demoApp) resetClient(reason string) {
	a.mu.Lock()
	oldClient := a.client
	a.client = nil
	cfg := a.cfg
	configErr := a.configErr
	a.mu.Unlock()

	if oldClient != nil {
		oldClient.Close()
	}

	if configErr != "" {
		a.recordEvent("reset_client", false, 0, configErr)
		return
	}
	if strings.TrimSpace(cfg.DeploymentURL) == "" {
		a.mu.Lock()
		a.configErr = "CONVEX_DEPLOYMENT_URL (or CONVEX_URL) is required"
		a.mu.Unlock()
		a.recordEvent("reset_client", false, 0, "missing deployment URL")
		return
	}

	client := convex.NewClientBuilder().
		WithDeploymentURL(cfg.DeploymentURL).
		WithClientID(cfg.ClientID).
		WithWebSocketStateCallback(func(state convex.WebSocketState) {
			a.recordEvent("websocket_state", true, 0, string(state))
		}).
		Build()

	if cfg.AuthToken != "" {
		token := cfg.AuthToken
		client.SetAuth(&token)
	}

	a.mu.Lock()
	a.client = client
	a.mu.Unlock()
	a.recordEvent("reset_client", true, 0, "reason="+reason)
}

func (a *demoApp) getClient() (*convex.Client, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.client != nil {
		return a.client, nil
	}
	if a.configErr != "" {
		return nil, errors.New(a.configErr)
	}
	return nil, errors.New("client is not configured")
}

func (a *demoApp) defaultsSnapshot() (queryFn, mutationFn, actionFn string, args map[string]any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cloned := make(map[string]any, len(a.defaultArgs))
	for key, value := range a.defaultArgs {
		cloned[key] = value
	}
	return a.cfg.DefaultQueryFn, a.cfg.DefaultMutationFn, a.cfg.DefaultActionFn, cloned
}

func (a *demoApp) statusSnapshot() (ready bool, configErr string, eventCount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.client != nil, a.configErr, len(a.events)
}

func (a *demoApp) recordEvent(endpoint string, success bool, duration time.Duration, message string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.nextID++
	a.events = append(a.events, debugEvent{
		ID:         a.nextID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Endpoint:   endpoint,
		Success:    success,
		DurationMS: duration.Milliseconds(),
		Message:    message,
	})
	if len(a.events) > 400 {
		a.events = append([]debugEvent(nil), a.events[len(a.events)-400:]...)
	}
}

func (a *demoApp) eventsSnapshot() []debugEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]debugEvent, len(a.events))
	copy(out, a.events)
	return out
}

func (a *demoApp) serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	queryFn, mutationFn, actionFn, _ := a.defaultsSnapshot()
	ready, configErr, eventCount := a.statusSnapshot()
	view := homeView{
		Port:              a.cfg.Port,
		ClientID:          a.cfg.ClientID,
		StartedAt:         a.startedAt.UTC().Format(time.RFC3339),
		Uptime:            time.Since(a.startedAt).Round(time.Second).String(),
		DeploymentURL:     a.cfg.DeploymentURL,
		HasAuthToken:      a.cfg.AuthToken != "",
		Ready:             ready,
		ConfigError:       configErr,
		DefaultQueryFn:    queryFn,
		DefaultMutationFn: mutationFn,
		DefaultActionFn:   actionFn,
		DefaultArgsRaw:    a.cfg.DefaultArgsRaw,
		EventCount:        eventCount,
	}
	if err := homeTmpl.Execute(w, view); err != nil {
		http.Error(w, "failed to render home", http.StatusInternalServerError)
	}
}

func (a *demoApp) jsonHandler(name string, handler func(*http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		payload, err := handler(r)
		duration := time.Since(start)
		status := http.StatusOK
		response := map[string]any{
			"ok":          err == nil,
			"endpoint":    name,
			"timestamp":   time.Now().UTC().Format(time.RFC3339Nano),
			"duration_ms": duration.Milliseconds(),
			"request": map[string]any{
				"method": r.Method,
				"path":   r.URL.Path,
				"query":  r.URL.Query(),
			},
		}
		if err != nil {
			status = http.StatusInternalServerError
			response["error"] = err.Error()
			a.recordEvent(name, false, duration, err.Error())
		} else {
			response["data"] = payload
			a.recordEvent(name, true, duration, "ok")
		}
		writeJSON(w, status, response)
	}
}

func (a *demoApp) handleDebug(*http.Request) (any, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	ready, configErr, eventCount := a.statusSnapshot()
	queryFn, mutationFn, actionFn, args := a.defaultsSnapshot()
	connection := a.connectionDebugInfo()

	return map[string]any{
		"started_at": a.startedAt.UTC().Format(time.RFC3339Nano),
		"uptime":     time.Since(a.startedAt).Round(time.Millisecond).String(),
		"go_version": runtime.Version(),
		"goroutines": runtime.NumGoroutine(),
		"pid":        os.Getpid(),
		"config": map[string]any{
			"port":                 a.cfg.Port,
			"client_id":            a.cfg.ClientID,
			"deployment_url":       a.cfg.DeploymentURL,
			"resolved_ws_url":      connection["resolved_ws_url"],
			"resolved_ws_error":    connection["resolved_ws_error"],
			"auth_token_set":       a.cfg.AuthToken != "",
			"ready":                ready,
			"config_error":         configErr,
			"default_query":        queryFn,
			"default_mutation":     mutationFn,
			"default_action":       actionFn,
			"default_args":         args,
			"captured_event_count": eventCount,
		},
		"memory": map[string]any{
			"alloc_bytes":      mem.Alloc,
			"heap_alloc_bytes": mem.HeapAlloc,
			"sys_bytes":        mem.Sys,
			"num_gc":           mem.NumGC,
		},
		"connection": connection,
	}, nil
}

func (a *demoApp) handleEvents(*http.Request) (any, error) {
	events := a.eventsSnapshot()
	return map[string]any{"count": len(events), "events": events}, nil
}

func (a *demoApp) handleLiveConnection(*http.Request) (any, error) {
	return a.connectionDebugInfo(), nil
}

func (a *demoApp) handleLiveReset(*http.Request) (any, error) {
	a.resetClient("manual_reset")
	ready, configErr, _ := a.statusSnapshot()
	if !ready {
		return map[string]any{"ready": false, "config_error": configErr}, errors.New("client not ready after reset")
	}
	return map[string]any{"ready": true}, nil
}

func (a *demoApp) handleLiveQuery(r *http.Request) (any, error) {
	return a.handleLiveCall(r, "query")
}

func (a *demoApp) handleLiveMutation(r *http.Request) (any, error) {
	return a.handleLiveCall(r, "mutation")
}

func (a *demoApp) handleLiveAction(r *http.Request) (any, error) {
	return a.handleLiveCall(r, "action")
}

func (a *demoApp) handleLiveInsert(r *http.Request) (any, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	mutationFn := a.defaultFunction("mutation")
	if mutationFn == "" {
		mutationFn = "messages:create"
	}

	seed := time.Now().UTC().Format("20060102T150405.000000000Z")
	args := map[string]any{
		"text":   "Demo insert " + seed,
		"source": "mutation",
	}

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	result, err := client.Mutation(ctx, mutationFn, args)
	if err != nil {
		return map[string]any{
			"mutation_function": mutationFn,
			"args":              args,
		}, err
	}
	value, err := result.Unwrap()
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"mutation_function": mutationFn,
		"args":              args,
		"value":             value.Raw(),
		"value_type":        typeName(value.Raw()),
	}, nil
}

func (a *demoApp) handleLiveCall(r *http.Request, kind string) (any, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	function := strings.TrimSpace(r.URL.Query().Get("function"))
	if function == "" {
		function = a.defaultFunction(kind)
	}
	if function == "" {
		return nil, fmt.Errorf("missing function path for %s; provide ?function=module:name or set env default", kind)
	}

	args, err := a.argsFromRequest(r)
	if err != nil {
		return nil, err
	}

	timeoutMS := parsePositiveInt(r.URL.Query().Get("timeout_ms"), 20000)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	start := time.Now()
	result, err := callConvex(client, ctx, kind, function, args)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return map[string]any{
			"kind":       kind,
			"function":   function,
			"args":       args,
			"timeout_ms": timeoutMS,
			"latency_ms": latency,
		}, err
	}

	value, err := result.Unwrap()
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"kind":       kind,
		"function":   function,
		"args":       args,
		"timeout_ms": timeoutMS,
		"latency_ms": latency,
		"value":      value.Raw(),
		"value_type": typeName(value.Raw()),
		"hint":       "If value is null unexpectedly, call /api/live/reset and retry. Also verify /api/live/connection URL matches your active Convex deployment.",
	}, nil
}

func (a *demoApp) handleLiveSubscribe(r *http.Request) (any, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	function := strings.TrimSpace(r.URL.Query().Get("function"))
	if function == "" {
		function = a.defaultFunction("query")
	}
	if function == "" {
		return nil, errors.New("missing query function; set DEMO_QUERY_FUNCTION or pass ?function=")
	}

	args, err := a.argsFromRequest(r)
	if err != nil {
		return nil, err
	}

	timeoutMS := parsePositiveInt(r.URL.Query().Get("timeout_ms"), 20000)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	sub, err := client.Subscribe(ctx, function, args)
	if err != nil {
		return nil, err
	}
	defer sub.Close()

	select {
	case value, ok := <-sub.Updates():
		if !ok {
			return nil, errors.New("subscription channel closed before first update")
		}
		return map[string]any{
			"function":   function,
			"args":       args,
			"timeout_ms": timeoutMS,
			"value":      value.Raw(),
			"value_type": typeName(value.Raw()),
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (a *demoApp) handleLiveWatch(r *http.Request) (any, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	function := strings.TrimSpace(r.URL.Query().Get("function"))
	if function == "" {
		function = a.defaultFunction("query")
	}
	if function == "" {
		return nil, errors.New("missing query function; set DEMO_QUERY_FUNCTION or pass ?function=")
	}

	args, err := a.argsFromRequest(r)
	if err != nil {
		return nil, err
	}

	timeoutMS := parsePositiveInt(r.URL.Query().Get("timeout_ms"), 20000)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	sub, err := client.Subscribe(ctx, function, args)
	if err != nil {
		return nil, err
	}
	defer sub.Close()

	select {
	case <-sub.Updates():
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	watch := client.WatchAll()
	defer watch.Close()

	select {
	case snapshot, ok := <-watch.Updates():
		if !ok {
			return nil, errors.New("watch channel closed before first snapshot")
		}
		return map[string]any{
			"function":         function,
			"subscriber_count": len(snapshot),
			"snapshot":         snapshotToRaw(snapshot),
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (a *demoApp) handleLiveRunSuite(r *http.Request) (any, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	queryFn := strings.TrimSpace(r.URL.Query().Get("query_function"))
	if queryFn == "" {
		queryFn = a.defaultFunction("query")
	}
	if queryFn == "" {
		return nil, errors.New("missing query function; set DEMO_QUERY_FUNCTION or pass ?query_function=")
	}

	mutationFn := strings.TrimSpace(r.URL.Query().Get("mutation_function"))
	if mutationFn == "" {
		mutationFn = a.defaultFunction("mutation")
	}
	actionFn := strings.TrimSpace(r.URL.Query().Get("action_function"))
	if actionFn == "" {
		actionFn = a.defaultFunction("action")
	}

	args, err := a.argsFromRequest(r)
	if err != nil {
		return nil, err
	}

	timeoutMS := parsePositiveInt(r.URL.Query().Get("timeout_ms"), 20000)
	run := func(name string, fn func(context.Context) (map[string]any, error), optional bool) suiteCheckResult {
		result := suiteCheckResult{Name: name}
		start := time.Now()
		defer func() { result.DurationMS = time.Since(start).Milliseconds() }()
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMS)*time.Millisecond)
		defer cancel()
		data, err := fn(ctx)
		if err != nil {
			if optional {
				result.Skipped = true
				result.Error = err.Error()
				return result
			}
			result.Passed = false
			result.Error = err.Error()
			return result
		}
		result.Passed = true
		result.Data = data
		return result
	}

	results := []suiteCheckResult{
		run("query", func(ctx context.Context) (map[string]any, error) {
			result, err := client.Query(ctx, queryFn, args)
			if err != nil {
				return nil, err
			}
			value, err := result.Unwrap()
			if err != nil {
				return nil, err
			}
			return map[string]any{"value": value.Raw(), "value_type": typeName(value.Raw())}, nil
		}, false),
		run("subscribe", func(ctx context.Context) (map[string]any, error) {
			sub, err := client.Subscribe(ctx, queryFn, args)
			if err != nil {
				return nil, err
			}
			defer sub.Close()
			select {
			case value, ok := <-sub.Updates():
				if !ok {
					return nil, errors.New("subscription channel closed")
				}
				return map[string]any{"value": value.Raw(), "value_type": typeName(value.Raw())}, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}, false),
		run("watch", func(ctx context.Context) (map[string]any, error) {
			sub, err := client.Subscribe(ctx, queryFn, args)
			if err != nil {
				return nil, err
			}
			defer sub.Close()
			select {
			case <-sub.Updates():
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			watch := client.WatchAll()
			defer watch.Close()
			select {
			case snapshot, ok := <-watch.Updates():
				if !ok {
					return nil, errors.New("watch channel closed")
				}
				return map[string]any{"subscriber_count": len(snapshot), "snapshot": snapshotToRaw(snapshot)}, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}, false),
		run("mutation", func(ctx context.Context) (map[string]any, error) {
			if mutationFn == "" {
				return nil, errors.New("mutation function not configured")
			}
			result, err := client.Mutation(ctx, mutationFn, args)
			if err != nil {
				return nil, err
			}
			value, err := result.Unwrap()
			if err != nil {
				return nil, err
			}
			return map[string]any{"value": value.Raw(), "value_type": typeName(value.Raw())}, nil
		}, true),
		run("action", func(ctx context.Context) (map[string]any, error) {
			if actionFn == "" {
				return nil, errors.New("action function not configured")
			}
			result, err := client.Action(ctx, actionFn, args)
			if err != nil {
				return nil, err
			}
			value, err := result.Unwrap()
			if err != nil {
				return nil, err
			}
			return map[string]any{"value": value.Raw(), "value_type": typeName(value.Raw())}, nil
		}, true),
	}

	passed := 0
	failed := 0
	skipped := 0
	for _, result := range results {
		if result.Skipped {
			skipped++
			continue
		}
		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	return map[string]any{
		"query_function":    queryFn,
		"mutation_function": mutationFn,
		"action_function":   actionFn,
		"args":              args,
		"timeout_ms":        timeoutMS,
		"passed":            passed,
		"failed":            failed,
		"skipped":           skipped,
		"results":           results,
	}, nil
}

func (a *demoApp) handleLiveStream(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	client, err := a.getClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		a.recordEvent("live_stream", false, time.Since(start), err.Error())
		return
	}

	function := strings.TrimSpace(r.URL.Query().Get("function"))
	if function == "" {
		function = a.defaultFunction("query")
	}
	if function == "" {
		http.Error(w, "missing query function; set DEMO_QUERY_FUNCTION or pass ?function=", http.StatusBadRequest)
		a.recordEvent("live_stream", false, time.Since(start), "missing query function")
		return
	}

	args, err := a.argsFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.recordEvent("live_stream", false, time.Since(start), err.Error())
		return
	}

	timeoutSeconds := parsePositiveInt(r.URL.Query().Get("timeout_seconds"), 300)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	sub, err := client.Subscribe(ctx, function, args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.recordEvent("live_stream", false, time.Since(start), err.Error())
		return
	}
	defer sub.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		a.recordEvent("live_stream", false, time.Since(start), "flusher unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	_ = writeSSEEvent(w, "start", map[string]any{"function": function, "args": args, "timeout_seconds": timeoutSeconds})
	flusher.Flush()

	updateCount := 0
	for {
		select {
		case value, ok := <-sub.Updates():
			if !ok {
				_ = writeSSEEvent(w, "done", map[string]any{"updates": updateCount, "reason": "subscription closed"})
				flusher.Flush()
				a.recordEvent("live_stream", true, time.Since(start), fmt.Sprintf("subscription closed after %d updates", updateCount))
				return
			}
			updateCount++
			_ = writeSSEEvent(w, "update", map[string]any{
				"sequence":   updateCount,
				"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
				"value":      value.Raw(),
				"value_type": typeName(value.Raw()),
			})
			flusher.Flush()
		case <-ctx.Done():
			_ = writeSSEEvent(w, "done", map[string]any{"updates": updateCount, "reason": ctx.Err().Error()})
			flusher.Flush()
			a.recordEvent("live_stream", true, time.Since(start), fmt.Sprintf("stream finished after %d updates", updateCount))
			return
		}
	}
}

func (a *demoApp) argsFromRequest(r *http.Request) (map[string]any, error) {
	argsRaw := strings.TrimSpace(r.URL.Query().Get("args"))
	if argsRaw == "" {
		_, _, _, defaults := a.defaultsSnapshot()
		return defaults, nil
	}
	return parseArgs(argsRaw)
}

func (a *demoApp) defaultFunction(kind string) string {
	queryFn, mutationFn, actionFn, _ := a.defaultsSnapshot()
	switch kind {
	case "query":
		return queryFn
	case "mutation":
		return mutationFn
	case "action":
		return actionFn
	default:
		return ""
	}
}

func (a *demoApp) connectionDebugInfo() map[string]any {
	resolvedWSURL := ""
	resolvedWSError := ""
	if strings.TrimSpace(a.cfg.DeploymentURL) != "" {
		wsURL, err := convsync.DeploymentURLToWebSocketURL(a.cfg.DeploymentURL)
		if err != nil {
			resolvedWSError = err.Error()
		} else {
			resolvedWSURL = wsURL
		}
	}

	events := a.eventsSnapshot()
	states := make([]string, 0, 12)
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		if e.Endpoint != "websocket_state" {
			continue
		}
		states = append(states, e.Message)
		if len(states) >= 12 {
			break
		}
	}

	lastState := ""
	if len(states) > 0 {
		lastState = states[0]
	}

	return map[string]any{
		"deployment_url":    a.cfg.DeploymentURL,
		"resolved_ws_url":   resolvedWSURL,
		"resolved_ws_error": resolvedWSError,
		"last_ws_state":     lastState,
		"recent_ws_states":  states,
		"env_keys": []string{
			"CONVEX_DEPLOYMENT_URL",
			"CONVEX_URL",
			"CONVEX_AUTH_TOKEN",
			"DEMO_QUERY_FUNCTION",
			"DEMO_DEFAULT_ARGS",
		},
	}
}

func callConvex(client *convex.Client, ctx context.Context, kind, function string, args map[string]any) (convex.FunctionResult, error) {
	switch kind {
	case "query":
		return client.Query(ctx, function, args)
	case "mutation":
		return client.Mutation(ctx, function, args)
	case "action":
		return client.Action(ctx, function, args)
	default:
		return convex.Failure(fmt.Errorf("unsupported call kind %q", kind)), fmt.Errorf("unsupported call kind %q", kind)
	}
}

func snapshotToRaw(snapshot map[int64]convex.Value) map[string]any {
	out := make(map[string]any, len(snapshot))
	for id, value := range snapshot {
		out[fmt.Sprintf("%d", id)] = value.Raw()
	}
	return out
}

func parseArgs(raw string) (map[string]any, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return map[string]any{}, nil
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil, fmt.Errorf("args must be a JSON object: %w", err)
	}
	if parsed == nil {
		return map[string]any{}, nil
	}
	return parsed, nil
}

func writeSSEEvent(w http.ResponseWriter, event string, payload any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\n", event); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", string(encoded)); err != nil {
		return err
	}
	return nil
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func typeName(value any) string {
	if value == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%T", value)
}
