package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/get-convex/convex-go/convex"
)

func TestLiveConvexIntegration(t *testing.T) {
	if os.Getenv("CONVEX_INTEGRATION") == "" {
		t.Skip("set CONVEX_INTEGRATION=1 to run live integration tests")
	}

	deploymentURL := os.Getenv("CONVEX_DEPLOYMENT_URL")
	queryName := os.Getenv("CONVEX_TEST_QUERY")
	if deploymentURL == "" || queryName == "" {
		t.Skip("set CONVEX_DEPLOYMENT_URL and CONVEX_TEST_QUERY")
	}

	client := convex.NewClientBuilder().WithDeploymentURL(deploymentURL).WithClientID("integration-test").Build()
	defer client.Close()

	if token := os.Getenv("CONVEX_AUTH_TOKEN"); token != "" {
		client.SetAuth(&token)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := client.Query(ctx, queryName, map[string]any{})
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if _, err := result.Unwrap(); err != nil {
		t.Fatalf("query unwrap failed: %v", err)
	}

	sub, err := client.Subscribe(ctx, queryName, map[string]any{})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}
	defer sub.Close()

	select {
	case <-sub.Updates():
	case <-time.After(20 * time.Second):
		t.Fatalf("timed out waiting for subscription update")
	}

	if mutationName := os.Getenv("CONVEX_TEST_MUTATION"); mutationName != "" {
		mutationResult, err := client.Mutation(ctx, mutationName, map[string]any{})
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}
		if _, err := mutationResult.Unwrap(); err != nil {
			t.Fatalf("mutation unwrap failed: %v", err)
		}
	}

	if actionName := os.Getenv("CONVEX_TEST_ACTION"); actionName != "" {
		actionResult, err := client.Action(ctx, actionName, map[string]any{})
		if err != nil {
			t.Fatalf("action failed: %v", err)
		}
		if _, err := actionResult.Unwrap(); err != nil {
			t.Fatalf("action unwrap failed: %v", err)
		}
	}
}
