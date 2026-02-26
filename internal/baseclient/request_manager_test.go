package baseclient

import (
	"testing"

	"github.com/get-convex/convex-go/internal/protocol"
)

func TestRequestManagerMutationVisibilityAndActionCompletion(t *testing.T) {
	manager := NewRequestManager()
	manager.Add(1, RequestKindMutation)
	manager.Add(2, RequestKindAction)

	if ok := manager.HandleMutationResponse(1, protocol.NewTimestamp(5), false); !ok {
		t.Fatalf("expected mutation response to resolve existing request")
	}
	if pending, ok := manager.Pending(1); !ok || pending.Completed {
		t.Fatalf("expected mutation to wait on visibility transition")
	}

	completed := manager.ApplyTransition(protocol.NewTimestamp(4))
	if len(completed) != 0 {
		t.Fatalf("did not expect completion before visibility timestamp")
	}
	completed = manager.ApplyTransition(protocol.NewTimestamp(5))
	if len(completed) != 1 || completed[0] != 1 {
		t.Fatalf("expected mutation completion at visibility timestamp, got %v", completed)
	}

	if ok := manager.HandleActionResponse(2, false); !ok {
		t.Fatalf("expected action response to resolve existing request")
	}
	if pending, ok := manager.Pending(2); !ok || !pending.Completed {
		t.Fatalf("expected action to complete immediately")
	}
}

func TestRequestManagerReplayOrderDeterministic(t *testing.T) {
	manager := NewRequestManager()
	manager.Add(10, RequestKindMutation)
	manager.Add(11, RequestKindMutation)
	manager.Add(12, RequestKindAction)
	manager.Add(11, RequestKindMutation)

	if ok := manager.HandleActionResponse(12, false); !ok {
		t.Fatalf("expected action completion")
	}
	if ok := manager.HandleMutationResponse(10, protocol.NewTimestamp(1), false); !ok {
		t.Fatalf("expected mutation response for 10")
	}
	manager.ApplyTransition(protocol.NewTimestamp(1))

	replay := manager.ReplayOrder()
	if len(replay) != 1 || replay[0] != 11 {
		t.Fatalf("expected only unresolved request 11 in replay order, got %v", replay)
	}
}
