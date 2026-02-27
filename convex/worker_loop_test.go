package convex

import (
	"testing"
	"time"

	"github.com/get-convex/convex-go/internal/protocol"
	syncproto "github.com/get-convex/convex-go/internal/sync"
)

func TestWorkerLoopProcessesCommandAndTransportMessage(t *testing.T) {
	client := newClient()
	responses := make(chan syncproto.ProtocolResponse, 1)

	client.mu.Lock()
	client.responses = responses
	client.workerStarted = true
	client.mu.Unlock()

	go client.workerLoop()

	responses <- syncproto.ProtocolResponse{Message: &protocol.ServerMessage{Type: "Ping"}}

	resultCh := make(chan workerCommandResult, 1)
	client.workerCommands <- workerCommand{kind: workerCommandClose, result: resultCh}
	result := <-resultCh
	if result.err != nil {
		t.Fatalf("close command failed: %v", result.err)
	}

	client.Close()
	select {
	case <-client.workerDone:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for worker shutdown")
	}
}

func TestWorkerCommandUnsupportedKind(t *testing.T) {
	client := newClient()
	resultCh := make(chan workerCommandResult, 1)
	client.handleWorkerCommand(workerCommand{kind: "unknown", result: resultCh})
	result := <-resultCh
	if result.err == nil {
		t.Fatalf("expected unsupported command error")
	}
}
