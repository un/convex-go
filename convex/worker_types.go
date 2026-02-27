package convex

import (
	"context"

	"github.com/get-convex/convex-go/internal/protocol"
	syncproto "github.com/get-convex/convex-go/internal/sync"
)

type workerCommandKind string

const (
	workerCommandSubscribe   workerCommandKind = "subscribe"
	workerCommandUnsubscribe workerCommandKind = "unsubscribe"
	workerCommandMutation    workerCommandKind = "mutation"
	workerCommandAction      workerCommandKind = "action"
	workerCommandCancelReq   workerCommandKind = "cancel_request"
	workerCommandSetAuth     workerCommandKind = "set_auth"
	workerCommandSetAuthCB   workerCommandKind = "set_auth_callback"
	workerCommandWatchAll    workerCommandKind = "watch_all"
	workerCommandClose       workerCommandKind = "close"
)

type workerCommand struct {
	kind   workerCommandKind
	ctx    context.Context
	value  any
	result chan workerCommandResult
}

type workerSubscribePayload struct {
	name string
	args map[string]any
}

type workerSubscribeResult struct {
	subID   int64
	updates chan Value
}

type workerUnsubscribePayload struct {
	subID int64
}

type workerRunRequestPayload struct {
	name string
	args map[string]any
}

type workerRunRequestResult struct {
	requestID protocol.RequestSequenceNumber
	response  <-chan FunctionResult
}

type workerCancelRequestPayload struct {
	requestID protocol.RequestSequenceNumber
	err       error
}

type workerCommandResult struct {
	value any
	err   error
}

func (cmd workerCommand) cancelled() bool {
	if cmd.ctx == nil {
		return false
	}
	select {
	case <-cmd.ctx.Done():
		return true
	default:
		return false
	}
}

func (cmd workerCommand) cancelErr() error {
	if cmd.ctx == nil {
		return nil
	}
	return cmd.ctx.Err()
}

func (cmd workerCommand) resolve(value any, err error) {
	if cmd.result == nil {
		return
	}
	cmd.result <- workerCommandResult{value: value, err: err}
}

type workerEventKind string

const (
	workerEventCommand       workerEventKind = "command"
	workerEventTransportMsg  workerEventKind = "transport_message"
	workerEventTransportErr  workerEventKind = "transport_error"
	workerEventTransportDone workerEventKind = "transport_done"
)

type workerEvent struct {
	kind    workerEventKind
	command *workerCommand
	message *protocol.ServerMessage
	err     error
}

func workerEventFromProtocolResponse(response syncproto.ProtocolResponse) workerEvent {
	if response.Err != nil {
		return workerEvent{kind: workerEventTransportErr, err: response.Err}
	}
	if response.Message == nil {
		return workerEvent{kind: workerEventTransportDone}
	}
	return workerEvent{kind: workerEventTransportMsg, message: response.Message}
}
