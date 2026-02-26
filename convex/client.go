package convex

import (
    "context"
    "fmt"
    "sync"
)

type AuthTokenFetcher func(forceRefresh bool) (*string, error)

type Client struct {
    mu            sync.RWMutex
    closed        bool
    nextSubID     int64
    clientID      string
    stateCallback StateCallback
    authToken     *string
    authFetcher   AuthTokenFetcher
    subscriptions map[int64]chan Value
    watchers      map[int64]chan map[int64]Value
    nextWatcherID int64
    latestValue   Value
}

func newClient() *Client {
    return &Client{
        subscriptions: map[int64]chan Value{},
        watchers:      map[int64]chan map[int64]Value{},
        latestValue:   NewNullValue(),
    }
}

func NewClient() *Client {
    return NewClientBuilder().Build()
}

func (c *Client) Clone() *Client {
    return c
}

func (c *Client) Subscribe(ctx context.Context, name string, args map[string]any) (*QuerySubscription, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.closed {
        return nil, fmt.Errorf("client closed")
    }
    _ = ctx
    _ = name
    _ = args
    c.nextSubID++
    subID := c.nextSubID
    updates := make(chan Value, 4)
    updates <- c.latestValue
    c.subscriptions[subID] = updates
    return &QuerySubscription{
        UpdatesCh: updates,
        closeFn: func() {
            c.mu.Lock()
            defer c.mu.Unlock()
            if ch, ok := c.subscriptions[subID]; ok {
                close(ch)
                delete(c.subscriptions, subID)
            }
        },
    }, nil
}

func (c *Client) Query(ctx context.Context, name string, args map[string]any) (FunctionResult, error) {
    sub, err := c.Subscribe(ctx, name, args)
    if err != nil {
        return Failure(err), err
    }
    defer sub.Close()

    select {
    case value := <-sub.Updates():
        return Success(value), nil
    case <-ctx.Done():
        return Failure(ctx.Err()), ctx.Err()
    }
}

func (c *Client) Mutation(ctx context.Context, name string, args map[string]any) (FunctionResult, error) {
    _ = ctx
    c.mu.Lock()
    c.latestValue = NewValue(map[string]any{"operation": "mutation", "name": name, "args": args})
    c.broadcastLocked(c.latestValue)
    c.mu.Unlock()
    return Success(c.latestValue), nil
}

func (c *Client) Action(ctx context.Context, name string, args map[string]any) (FunctionResult, error) {
    _ = ctx
    c.mu.Lock()
    c.latestValue = NewValue(map[string]any{"operation": "action", "name": name, "args": args})
    c.broadcastLocked(c.latestValue)
    c.mu.Unlock()
    return Success(c.latestValue), nil
}

func (c *Client) WatchAll() *QuerySetSubscription {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.nextWatcherID++
    watcherID := c.nextWatcherID
    updates := make(chan map[int64]Value, 4)
    updates <- c.snapshotLocked()
    c.watchers[watcherID] = updates
    return &QuerySetSubscription{
        UpdatesCh: updates,
        closeFn: func() {
            c.mu.Lock()
            defer c.mu.Unlock()
            if ch, ok := c.watchers[watcherID]; ok {
                close(ch)
                delete(c.watchers, watcherID)
            }
        },
    }
}

func (c *Client) SetAuth(token *string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.authToken = token
}

func (c *Client) SetAuthCallback(fetcher AuthTokenFetcher) error {
    c.mu.Lock()
    c.authFetcher = fetcher
    c.mu.Unlock()
    token, err := fetcher(false)
    if err != nil {
        return err
    }
    c.SetAuth(token)
    return nil
}

func (c *Client) Close() {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.closed {
        return
    }
    c.closed = true
    for id, ch := range c.subscriptions {
        close(ch)
        delete(c.subscriptions, id)
    }
    for id, ch := range c.watchers {
        close(ch)
        delete(c.watchers, id)
    }
}

func (c *Client) broadcastLocked(value Value) {
    for _, sub := range c.subscriptions {
        select {
        case sub <- value:
        default:
        }
    }
    snapshot := c.snapshotLocked()
    for _, watcher := range c.watchers {
        select {
        case watcher <- snapshot:
        default:
        }
    }
}

func (c *Client) snapshotLocked() map[int64]Value {
    out := map[int64]Value{}
    for subID := range c.subscriptions {
        out[subID] = c.latestValue
    }
    return out
}
