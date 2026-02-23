package core

type CloseCode int

const (
	CloseNormalClosure   CloseCode = 1000
	CloseGoingAway       CloseCode = 1001
	CloseAbnormalClosure CloseCode = 1006
	CloseInternalError   CloseCode = 1011
)

type MessageType int

const (
	MessageText   MessageType = 1
	MessageBinary MessageType = 2
	MessageClose  MessageType = 8
	MessagePing   MessageType = 9
	MessagePong   MessageType = 10
)
