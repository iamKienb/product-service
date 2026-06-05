package port

import (
	"context"
	"encoding/json"

	kafkax "github.com/iamKienb/go-core/kafka"
)

type Message = kafkax.Message
type EventProcessor interface {
	Handle(context.Context, Message) error
}
type EventHandler interface {
	Handle(context.Context, json.RawMessage) error
}
