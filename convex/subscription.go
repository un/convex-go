package convex

import "sync"

type QuerySubscription struct {
    UpdatesCh <-chan Value
    closeFn   func()
    once      sync.Once
}

func (s *QuerySubscription) Updates() <-chan Value {
    return s.UpdatesCh
}

func (s *QuerySubscription) Close() {
    s.once.Do(func() {
        if s.closeFn != nil {
            s.closeFn()
        }
    })
}

type QuerySetSubscription struct {
    UpdatesCh <-chan map[int64]Value
    closeFn   func()
    once      sync.Once
}

func (s *QuerySetSubscription) Updates() <-chan map[int64]Value {
    return s.UpdatesCh
}

func (s *QuerySetSubscription) Close() {
    s.once.Do(func() {
        if s.closeFn != nil {
            s.closeFn()
        }
    })
}
