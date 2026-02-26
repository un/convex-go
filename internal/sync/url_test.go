package sync

import "testing"

func TestDeploymentURLToWebSocketURL(t *testing.T) {
    got, err := DeploymentURLToWebSocketURL("https://example.convex.cloud/path")
    if err != nil {
        t.Fatalf("conversion failed: %v", err)
    }
    if got != "wss://example.convex.cloud/api/sync" {
        t.Fatalf("unexpected ws url: %s", got)
    }
}
