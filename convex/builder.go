package convex

type StateCallback func(WebSocketState)

type ClientBuilder struct {
    clientID      string
    stateCallback StateCallback
}

func NewClientBuilder() *ClientBuilder {
    return &ClientBuilder{}
}

func (b *ClientBuilder) WithClientID(clientID string) *ClientBuilder {
    b.clientID = clientID
    return b
}

func (b *ClientBuilder) WithWebSocketStateCallback(callback StateCallback) *ClientBuilder {
    b.stateCallback = callback
    return b
}

func (b *ClientBuilder) Build() *Client {
    c := newClient()
    c.clientID = b.clientID
    c.stateCallback = b.stateCallback
    return c
}
