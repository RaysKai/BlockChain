package peer_error

import (
	"errors"
	"fmt"
)

const (
	ErrInvalidMsgCode = iota
	ErrInvalidMsg
)

var errorToString = map[int]string{
	ErrInvalidMsgCode: "invalid message code",
	ErrInvalidMsg:     "invalid message",
}

type peerError struct {
	code    int
	message string
}

func NewPeerError(code int, format string, v ...interface{}) *peerError {
	desc, ok := errorToString[code]
	if !ok {
		panic("invalid error code")
	}
	err := &peerError{code, desc}
	if format != "" {
		err.message += ": " + fmt.Sprintf(format, v...)
	}
	return err
}

func (self *peerError) Error() string {
	return self.message
}

var ErrProtocolReturned = errors.New("protocol returned")

type DiscReason uint

const (
	DiscRequested DiscReason = iota
	DiscNetworkError
	DiscProtocolError
	DiscUselessPeer
	DiscTooManyPeers
	DiscAlreadyConnected
	DiscIncompatibleVersion
	DiscInvalidIdentity
	DiscQuitting
	DiscUnexpectedIdentity
	DiscSelf
	DiscReadTimeout
	DiscSubprotocolError = 0x10
)

var discReasonToString = [...]string{
	DiscRequested:           "disconnect requested",
	DiscNetworkError:        "network error",
	DiscProtocolError:       "breach of protocol",
	DiscUselessPeer:         "useless peer",
	DiscTooManyPeers:        "too many peers",
	DiscAlreadyConnected:    "already connected",
	DiscIncompatibleVersion: "incompatible p2p protocol version",
	DiscInvalidIdentity:     "invalid node identity",
	DiscQuitting:            "client quitting",
	DiscUnexpectedIdentity:  "unexpected identity",
	DiscSelf:                "connected to self",
	DiscReadTimeout:         "read timeout",
	DiscSubprotocolError:    "subprotocol error",
}

func (d DiscReason) String() string {
	if len(discReasonToString) < int(d) {
		return fmt.Sprintf("unknown disconnect reason %d", d)
	}
	return discReasonToString[d]
}

func (d DiscReason) Error() string {
	return d.String()
}

func DiscReasonForError(err error) DiscReason {
	if reason, ok := err.(DiscReason); ok {
		return reason
	}
	if err == ErrProtocolReturned {
		return DiscQuitting
	}
	peerError, ok := err.(*peerError)
	if ok {
		switch peerError.code {
		case ErrInvalidMsgCode, ErrInvalidMsg:
			return DiscProtocolError
		default:
			return DiscSubprotocolError
		}
	}
	return DiscSubprotocolError
}
