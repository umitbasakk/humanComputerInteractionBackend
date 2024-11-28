package model

type ErrorCode int32

const (
	MustbeVerified ErrorCode = 1 + iota
	ErrorVerifySystem
	ErrorLoginSystem
	ErrorRegisterSystem
	SendedMessage
	Authorized
	NoError
)

type MessageHandler struct {
	Message string
	ErrCode ErrorCode
	Data    interface{}
}
