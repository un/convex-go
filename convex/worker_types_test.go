package convex

import (
	"context"
	"errors"
	"testing"

	"github.com/get-convex/convex-go/internal/protocol"
	syncproto "github.com/get-convex/convex-go/internal/sync"
)

func TestWorkerCommandCancellationSemantics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := workerCommand{kind: workerCommandSubscribe, ctx: ctx, result: make(chan workerCommandResult, 1)}

	if cmd.cancelled() {
		t.Fatalf("expected command not cancelled before cancel")
	}
	cancel()
	if !cmd.cancelled() {
		t.Fatalf("expected command cancelled after cancel")
	}
	if !errors.Is(cmd.cancelErr(), context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", cmd.cancelErr())
	}

	cmd.resolve("ok", nil)
	result := <-cmd.result
	if got, ok := result.value.(string); !ok || got != "ok" {
		t.Fatalf("unexpected command result payload: %#v", result.value)
	}
}

func TestWorkerEventFromProtocolResponse(t *testing.T) {
	errEvent := workerEventFromProtocolResponse(syncproto.ProtocolResponse{Err: errors.New("boom")})
	if errEvent.kind != workerEventTransportErr {
		t.Fatalf("expected transport error event, got %s", errEvent.kind)
	}

	doneEvent := workerEventFromProtocolResponse(syncproto.ProtocolResponse{})
	if doneEvent.kind != workerEventTransportDone {
		t.Fatalf("expected transport done event, got %s", doneEvent.kind)
	}

	msg := protocol.ServerMessage{Type: "Ping"}
	msgEvent := workerEventFromProtocolResponse(syncproto.ProtocolResponse{Message: &msg})
	if msgEvent.kind != workerEventTransportMsg {
		t.Fatalf("expected transport message event, got %s", msgEvent.kind)
	}
	if msgEvent.message == nil || msgEvent.message.Type != "Ping" {
		t.Fatalf("expected ping message, got %#v", msgEvent.message)
	}
}
