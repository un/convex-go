package baseclient

type QueryResults struct {
    bySubscriber map[int64]any
}

func NewQueryResults(values map[int64]any) QueryResults {
    copied := map[int64]any{}
    for k, v := range values {
        copied[k] = v
    }
    return QueryResults{bySubscriber: copied}
}

func (r QueryResults) Get(subscriberID int64) (any, bool) {
    v, ok := r.bySubscriber[subscriberID]
    return v, ok
}

func (r QueryResults) Len() int {
    return len(r.bySubscriber)
}

func (r QueryResults) IsEmpty() bool {
    return len(r.bySubscriber) == 0
}

func (r QueryResults) Snapshot() map[int64]any {
    out := map[int64]any{}
    for k, v := range r.bySubscriber {
        out[k] = v
    }
    return out
}
