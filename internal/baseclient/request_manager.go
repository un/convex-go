package baseclient

import "github.com/get-convex/convex-go/internal/protocol"

type RequestKind string

const (
	RequestKindMutation RequestKind = "mutation"
	RequestKindAction   RequestKind = "action"
)

type PendingRequest struct {
	ID          uint64
	Kind        RequestKind
	VisibleAt   protocol.Timestamp
	WaitingOnTS bool
	Errored     bool
	Completed   bool
}

type RequestManager struct {
	pending map[uint64]*PendingRequest
	order   []uint64
}

func NewRequestManager() *RequestManager {
	return &RequestManager{pending: map[uint64]*PendingRequest{}}
}

func (m *RequestManager) Add(id uint64, kind RequestKind) {
	if _, exists := m.pending[id]; exists {
		return
	}
	m.pending[id] = &PendingRequest{ID: id, Kind: kind}
	m.order = append(m.order, id)
}

func (m *RequestManager) HandleActionResponse(id uint64, err bool) bool {
	request, ok := m.pending[id]
	if !ok {
		return false
	}
	request.Errored = err
	request.Completed = true
	return true
}

func (m *RequestManager) HandleMutationResponse(id uint64, ts protocol.Timestamp, err bool) bool {
	request, ok := m.pending[id]
	if !ok {
		return false
	}
	request.VisibleAt = ts
	request.Errored = err
	if err {
		request.Completed = true
		request.WaitingOnTS = false
	} else {
		request.WaitingOnTS = true
	}
	return true
}

func (m *RequestManager) ApplyTransition(ts protocol.Timestamp) []uint64 {
	completed := []uint64{}
	for _, request := range m.pending {
		if request.Kind == RequestKindMutation && !request.Completed && !request.Errored && request.WaitingOnTS {
			if request.VisibleAt <= ts {
				request.Completed = true
				request.WaitingOnTS = false
				completed = append(completed, request.ID)
			}
		}
	}
	return completed
}

func (m *RequestManager) ReplayOrder() []uint64 {
	out := make([]uint64, 0, len(m.order))
	for _, id := range m.order {
		request := m.pending[id]
		if request != nil && !request.Completed {
			out = append(out, id)
		}
	}
	return out
}

func (m *RequestManager) Pending(id uint64) (*PendingRequest, bool) {
	request, ok := m.pending[id]
	if !ok {
		return nil, false
	}
	copy := *request
	return &copy, true
}
