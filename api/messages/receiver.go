package messages

import (
	"encoding/json"
	"errors"
)

type Receiver struct {
	Orders     chan SubmitOrdersMessage
	Support    chan GiveSupportMessage
	Quit       chan QuitMessage
	Kick       chan KickMessage
	WinterVote chan WinterVoteMessage
	Errors     chan error
}

func NewReceiver() Receiver {
	return Receiver{
		Orders:     make(chan SubmitOrdersMessage),
		Support:    make(chan GiveSupportMessage),
		Quit:       make(chan QuitMessage),
		Kick:       make(chan KickMessage),
		WinterVote: make(chan WinterVoteMessage),
		Errors:     make(chan error),
	}
}

func (receiver *Receiver) HandleMessage(rawMessage []byte) {
	var baseMessage BaseMessage

	err := json.Unmarshal(rawMessage, &baseMessage)
	if err != nil {
		receiver.Errors <- err
		return
	}
	if baseMessage.Type == "" {
		receiver.Errors <- errors.New("error in deserializing message")
		return
	}

	switch baseMessage.Type {

	case SubmitOrdersMessageType:
		var ordersMessage SubmitOrdersMessage
		err := json.Unmarshal(rawMessage, &ordersMessage)
		if err != nil {
			receiver.Errors <- err
			return
		}

		receiver.Orders <- ordersMessage

	case GiveSupportMessageType:
		var supportMessage GiveSupportMessage
		err := json.Unmarshal(rawMessage, &supportMessage)
		if err != nil {
			receiver.Errors <- err
			return
		}

		receiver.Support <- supportMessage

	case QuitMessageType:
		var quitMessage QuitMessage
		err := json.Unmarshal(rawMessage, &quitMessage)
		if err != nil {
			receiver.Errors <- err
			return
		}

		receiver.Quit <- quitMessage

	case WinterVoteMessageType:
		var winterVoteMessage WinterVoteMessage
		err := json.Unmarshal(rawMessage, &winterVoteMessage)
		if err != nil {
			receiver.Errors <- err
			return
		}

		receiver.WinterVote <- winterVoteMessage

	default:
		receiver.Errors <- errors.New("unrecognized message type")
		return
	}
}