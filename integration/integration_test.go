package integration

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/get-convex/convex-go/convex"
)

type liveConfig struct {
	deploymentURL     string
	queryName         string
	mutationName      string
	actionName        string
	authToken         string
	authRefreshToken  string
	reconnectQuery    string
	expectReconnect   bool
	expectAuthRefresh bool
	probeDuration     time.Duration
}

func TestLiveConvexIntegration(t *testing.T) {
	config, ok := loadLiveConfig(t)
	if !ok {
		return
	}

	t.Run("query", func(t *testing.T) {
		client := newLiveClient(config)
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		result, err := client.Query(ctx, config.queryName, map[string]any{})
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if _, err := result.Unwrap(); err != nil {
			t.Fatalf("query unwrap failed: %v", err)
		}
	})

	t.Run("subscribe", func(t *testing.T) {
		client := newLiveClient(config)
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		sub, err := client.Subscribe(ctx, config.queryName, map[string]any{})
		if err != nil {
			t.Fatalf("subscribe failed: %v", err)
		}
		defer sub.Close()

		select {
		case <-sub.Updates():
		case <-time.After(20 * time.Second):
			t.Fatalf("timed out waiting for subscription update")
		}
	})

	t.Run("mutation", func(t *testing.T) {
		if config.mutationName == "" {
			t.Skip("set CONVEX_TEST_MUTATION to run mutation live test")
		}

		client := newLiveClient(config)
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		result, err := client.Mutation(ctx, config.mutationName, map[string]any{})
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}
		if _, err := result.Unwrap(); err != nil {
			t.Fatalf("mutation unwrap failed: %v", err)
		}
	})

	t.Run("action", func(t *testing.T) {
		if config.actionName == "" {
			t.Skip("set CONVEX_TEST_ACTION to run action live test")
		}

		client := newLiveClient(config)
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		result, err := client.Action(ctx, config.actionName, map[string]any{})
		if err != nil {
			t.Fatalf("action failed: %v", err)
		}
		if _, err := result.Unwrap(); err != nil {
			t.Fatalf("action unwrap failed: %v", err)
		}
	})

	t.Run("auth-refresh-callback", func(t *testing.T) {
		if config.authRefreshToken == "" {
			t.Skip("set CONVEX_AUTH_REFRESH_TOKEN to run auth refresh callback scenario")
		}

		stateEvents := make(chan convex.WebSocketState, 64)
		client := convex.NewClientBuilder().
			WithDeploymentURL(config.deploymentURL).
			WithClientID("integration-test-auth-refresh").
			WithWebSocketStateCallback(func(state convex.WebSocketState) {
				stateEvents <- state
			}).
			Build()
		defer client.Close()

		calls := make(chan bool, 8)
		err := client.SetAuthCallback(func(forceRefresh bool) (*string, error) {
			calls <- forceRefresh
			if forceRefresh {
				token := config.authRefreshToken
				return &token, nil
			}
			if config.authToken != "" {
				token := config.authToken
				return &token, nil
			}
			fallback := config.authRefreshToken
			return &fallback, nil
		})
		if err != nil {
			t.Fatalf("set auth callback failed: %v", err)
		}

		select {
		case force := <-calls:
			if force {
				t.Fatalf("expected initial auth callback with forceRefresh=false")
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for initial auth callback call")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		result, err := client.Query(ctx, config.queryName, map[string]any{})
		if err != nil {
			t.Fatalf("query with auth callback failed: %v", err)
		}
		if _, err := result.Unwrap(); err != nil {
			t.Fatalf("query unwrap failed: %v", err)
		}

		if config.expectAuthRefresh {
			deadline := time.NewTimer(config.probeDuration)
			defer deadline.Stop()
			for {
				select {
				case force := <-calls:
					if force {
						return
					}
				case state := <-stateEvents:
					if state == convex.WebSocketStateReconnecting {
						// continue waiting for forceRefresh callback
					}
				case <-deadline.C:
					t.Fatalf("expected forceRefresh auth callback during probe window")
				}
			}
		}
	})

	t.Run("reconnect-probe", func(t *testing.T) {
		stateEvents := make(chan convex.WebSocketState, 128)
		client := convex.NewClientBuilder().
			WithDeploymentURL(config.deploymentURL).
			WithClientID("integration-test-reconnect").
			WithWebSocketStateCallback(func(state convex.WebSocketState) {
				stateEvents <- state
			}).
			Build()
		defer client.Close()

		if config.authToken != "" {
			client.SetAuth(&config.authToken)
		}

		deadline := time.Now().Add(config.probeDuration)
		for time.Now().Before(deadline) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			_, err := client.Query(ctx, config.reconnectQuery, map[string]any{})
			cancel()
			if err != nil {
				t.Fatalf("reconnect probe query failed: %v", err)
			}
			time.Sleep(250 * time.Millisecond)
		}

		reconnectingSeen := false
		for {
			select {
			case state := <-stateEvents:
				if state == convex.WebSocketStateReconnecting {
					reconnectingSeen = true
				}
			default:
				if config.expectReconnect && !reconnectingSeen {
					t.Fatalf("expected reconnect state during reconnect probe")
				}
				return
			}
		}
	})
}

func loadLiveConfig(t *testing.T) (liveConfig, bool) {
	t.Helper()
	if os.Getenv("CONVEX_INTEGRATION") == "" {
		t.Skip("set CONVEX_INTEGRATION=1 to run live integration tests")
		return liveConfig{}, false
	}

	deploymentURL := os.Getenv("CONVEX_DEPLOYMENT_URL")
	queryName := os.Getenv("CONVEX_TEST_QUERY")
	if deploymentURL == "" || queryName == "" {
		t.Skip("set CONVEX_DEPLOYMENT_URL and CONVEX_TEST_QUERY")
		return liveConfig{}, false
	}

	probeSeconds := 10
	if rawProbe := os.Getenv("CONVEX_RECONNECT_PROBE_SECONDS"); rawProbe != "" {
		if parsed, err := strconv.Atoi(rawProbe); err == nil && parsed > 0 {
			probeSeconds = parsed
		}
	}

	config := liveConfig{
		deploymentURL:     deploymentURL,
		queryName:         queryName,
		mutationName:      os.Getenv("CONVEX_TEST_MUTATION"),
		actionName:        os.Getenv("CONVEX_TEST_ACTION"),
		authToken:         os.Getenv("CONVEX_AUTH_TOKEN"),
		authRefreshToken:  os.Getenv("CONVEX_AUTH_REFRESH_TOKEN"),
		reconnectQuery:    queryName,
		expectReconnect:   os.Getenv("CONVEX_EXPECT_RECONNECT") == "1",
		expectAuthRefresh: os.Getenv("CONVEX_EXPECT_AUTH_REFRESH") == "1",
		probeDuration:     time.Duration(probeSeconds) * time.Second,
	}
	if customReconnectQuery := os.Getenv("CONVEX_TEST_RECONNECT_QUERY"); customReconnectQuery != "" {
		config.reconnectQuery = customReconnectQuery
	}
	return config, true
}

func newLiveClient(config liveConfig) *convex.Client {
	client := convex.NewClientBuilder().
		WithDeploymentURL(config.deploymentURL).
		WithClientID("integration-test").
		Build()
	if config.authToken != "" {
		client.SetAuth(&config.authToken)
	}
	return client
}
