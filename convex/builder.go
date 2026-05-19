package convex

type StateCallback func(WebSocketState)

type FailureCallback func(err error)

type ClientBuilder struct {
	deploymentURL  string
	clientID       string
	stateCallback  StateCallback
	failureCallback FailureCallback
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func (b *ClientBuilder) WithDeploymentURL(deploymentURL string) *ClientBuilder {
	b.deploymentURL = deploymentURL
	return b
}

func (b *ClientBuilder) WithClientID(clientID string) *ClientBuilder {
	b.clientID = clientID
	return b
}

func (b *ClientBuilder) WithWebSocketStateCallback(callback StateCallback) *ClientBuilder {
	b.stateCallback = callback
	return b
}

func (b *ClientBuilder) WithFailureCallback(callback FailureCallback) *ClientBuilder {
	b.failureCallback = callback
	return b
}

func (b *ClientBuilder) Build() *Client {
	c := newClient()
	c.clientID = b.clientID
	c.stateCallback = b.stateCallback
	c.failureCallback = b.failureCallback
	c.deploymentURL = b.deploymentURL
	return c
}
