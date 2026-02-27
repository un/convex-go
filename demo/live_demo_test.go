package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/get-convex/convex-go/convex"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	args, err := parseArgs(`{"count":1,"name":"demo"}`)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if args["name"] != "demo" {
		t.Fatalf("expected parsed name demo, got %#v", args["name"])
	}

	if _, err := parseArgs(`[]`); err == nil {
		t.Fatalf("expected parseArgs to reject non-object JSON")
	}
}

func TestLiveDemoFlow(t *testing.T) {
	if os.Getenv("CONVEX_INTEGRATION") == "" {
		t.Skip("set CONVEX_INTEGRATION=1 to run live demo integration checks")
	}

	deploymentURL := envOrDefault("CONVEX_DEPLOYMENT_URL", os.Getenv("CONVEX_URL"))
	queryFn := envOrDefault("DEMO_QUERY_FUNCTION", os.Getenv("CONVEX_TEST_QUERY"))
	if deploymentURL == "" || queryFn == "" {
		t.Skip("set CONVEX_DEPLOYMENT_URL/CONVEX_URL and DEMO_QUERY_FUNCTION/CONVEX_TEST_QUERY")
	}

	client := convex.NewClientBuilder().WithDeploymentURL(deploymentURL).WithClientID("demo-live-test").Build()
	defer client.Close()

	if token := envOrDefault("DEMO_AUTH_TOKEN", os.Getenv("CONVEX_AUTH_TOKEN")); token != "" {
		client.SetAuth(&token)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.Query(ctx, queryFn, map[string]any{})
	if err != nil {
		t.Fatalf("live query failed: %v", err)
	}
	if _, err := result.Unwrap(); err != nil {
		t.Fatalf("live query unwrap failed: %v", err)
	}

	sub, err := client.Subscribe(ctx, queryFn, map[string]any{})
	if err != nil {
		t.Fatalf("live subscribe failed: %v", err)
	}
	defer sub.Close()

	select {
	case <-sub.Updates():
	case <-time.After(20 * time.Second):
		t.Fatalf("timed out waiting for first subscription update")
	}

	watch := client.WatchAll()
	defer watch.Close()

	select {
	case snapshot := <-watch.Updates():
		if len(snapshot) == 0 {
			t.Fatalf("expected non-empty watch snapshot")
		}
	case <-time.After(20 * time.Second):
		t.Fatalf("timed out waiting for watch snapshot")
	}

	if mutationFn := os.Getenv("DEMO_MUTATION_FUNCTION"); mutationFn != "" {
		result, err := client.Mutation(ctx, mutationFn, map[string]any{})
		if err != nil {
			t.Fatalf("live mutation failed: %v", err)
		}
		if _, err := result.Unwrap(); err != nil {
			t.Fatalf("live mutation unwrap failed: %v", err)
		}
	}

	if actionFn := os.Getenv("DEMO_ACTION_FUNCTION"); actionFn != "" {
		result, err := client.Action(ctx, actionFn, map[string]any{})
		if err != nil {
			t.Fatalf("live action failed: %v", err)
		}
		if _, err := result.Unwrap(); err != nil {
			t.Fatalf("live action unwrap failed: %v", err)
		}
	}
}
