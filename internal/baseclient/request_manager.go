package baseclient

import "github.com/get-convex/convex-go/internal/protocol"

type RequestKind string

const (
    RequestKindMutation RequestKind = "mutation"
    RequestKindAction   RequestKind = "action"
)

type PendingRequest struct {
    ID        uint64
    Kind      RequestKind
    VisibleAt protocol.Timestamp
    Errored   bool
    Completed bool
}

type RequestManager struct {
    pending map[uint64]*PendingRequest
    order   []uint64
}

func NewRequestManager() *RequestManager {
    return &RequestManager{pending: map[uint64]*PendingRequest{}}
}

func (m *RequestManager) Add(id uint64, kind RequestKind) {
    m.pending[id] = &PendingRequest{ID: id, Kind: kind}
    m.order = append(m.order, id)
}

func (m *RequestManager) HandleActionResponse(id uint64, err bool) {
    request, ok := m.pending[id]
    if !ok {
        return
    }
    request.Errored = err
    request.Completed = true
}

func (m *RequestManager) HandleMutationResponse(id uint64, ts protocol.Timestamp, err bool) {
    request, ok := m.pending[id]
    if !ok {
        return
    }
    request.VisibleAt = ts
    request.Errored = err
    if err {
        request.Completed = true
    }
}

func (m *RequestManager) ApplyTransition(ts protocol.Timestamp) {
    for _, request := range m.pending {
        if request.Kind == RequestKindMutation && !request.Completed && !request.Errored {
            if request.VisibleAt <= ts {
                request.Completed = true
            }
        }
    }
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
