package tgbot

import (
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type errorType int

const (
	errorTypeEmptyResponse errorType = iota
	errorTypeStringResponse
)

type tgBotError interface {
	errorType() errorType
}

var (
	_ tgBotError = (*MessageError)(nil)
	_ tgBotError = (*ExceptionError)(nil)
)

type MessageError struct {
	message          string
	replyToMessageID int
	editMessageID    int
}

func NewMessageError(message string) MessageError {
	return MessageError{
		message: message,
	}
}

func (e MessageError) errorType() errorType {
	return errorTypeStringResponse
}

func (e MessageError) Error() string {
	return e.message
}

func (e MessageError) WithReply(message *tgbotapi.Message) MessageError {
	e.replyToMessageID = message.MessageID
	return e
}

func (e MessageError) WithEdit(message *tgbotapi.Message) MessageError {
	e.editMessageID = message.MessageID
	return e
}

type ExceptionError struct {
	err              error
	message          string
	replyToMessageID int
	editMessageID    int
	callFrameSkip    int
	callFrame        *runtime.Frame
}

func NewExceptionError(err error) ExceptionError {
	e := ExceptionError{
		err:           err,
		callFrameSkip: 1,
	}

	pc, file, line, _ := runtime.Caller(e.callFrameSkip)
	funcDetails := runtime.FuncForPC(pc)
	var funcName string
	if funcDetails != nil {
		funcName = funcDetails.Name()
	}

	e.callFrame = &runtime.Frame{
		File:     file,
		Line:     line,
		Function: funcName,
	}

	return e
}

func (e ExceptionError) errorType() errorType {
	return errorTypeEmptyResponse
}

func (e ExceptionError) Error() string {
	return e.err.Error()
}

func (e ExceptionError) WithMessage(message string) ExceptionError {
	e.message = message
	return e
}

func (e ExceptionError) WithReply(message *tgbotapi.Message) ExceptionError {
	e.replyToMessageID = message.MessageID
	return e
}

func (e ExceptionError) WithEdit(message *tgbotapi.Message) ExceptionError {
	e.editMessageID = message.MessageID
	return e
}
