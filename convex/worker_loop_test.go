package convex

import (
	"context"
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

func TestWorkerFlushesOutboundBeforeSelect(t *testing.T) {
	client := newClient()
	responses := make(chan syncproto.ProtocolResponse)
	sent := make(chan string, 1)

	client.mu.Lock()
	client.responses = responses
	client.connected = true
	client.workerStarted = true
	client.sendFn = func(_ context.Context, message protocol.ClientMessage) error {
		sent <- message.Type
		return nil
	}
	client.mu.Unlock()

	go client.workerLoop()

	client.enqueueOutbound(protocol.ClientMessage{Type: "Event", Event: []byte(`"flush-first"`)})
	resultCh := make(chan workerCommandResult)
	client.workerCommands <- workerCommand{kind: "unknown", result: resultCh}

	select {
	case messageType := <-sent:
		if messageType != "Event" {
			t.Fatalf("unexpected message type flushed: %s", messageType)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("expected outbound queue flush before command handling")
	}

	result := <-resultCh
	if result.err == nil {
		t.Fatalf("expected unsupported command error")
	}

	client.Close()
}

func TestWorkerCommunicatesDuringSend(t *testing.T) {
	client := newClient()
	responses := make(chan syncproto.ProtocolResponse, 1)
	releaseSend := make(chan struct{})
	events := make(chan workerEventKind, 2)

	client.mu.Lock()
	client.responses = responses
	client.connected = true
	client.workerStarted = true
	client.sendFn = func(_ context.Context, _ protocol.ClientMessage) error {
		<-releaseSend
		return nil
	}
	client.workerEventHook = func(event workerEvent) {
		events <- event.kind
	}
	client.mu.Unlock()

	go client.workerLoop()

	client.enqueueOutbound(protocol.ClientMessage{Type: "Event", Event: []byte(`"blocking-send"`)})
	responses <- syncproto.ProtocolResponse{Message: &protocol.ServerMessage{Type: "Ping"}}

	select {
	case kind := <-events:
		if kind != workerEventTransportMsg {
			t.Fatalf("unexpected worker event kind: %s", kind)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("expected worker to process transport event while send was blocked")
	}

	close(releaseSend)
	client.Close()
}
