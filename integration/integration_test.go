package integration

import (
    "os"
    "testing"
)

func TestLiveConvexIntegration(t *testing.T) {
    if os.Getenv("CONVEX_INTEGRATION") == "" {
        t.Skip("set CONVEX_INTEGRATION=1 to run live integration tests")
    }
}
