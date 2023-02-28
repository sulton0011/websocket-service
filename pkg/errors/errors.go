package errors

import (
	"errors"
	"strings"
	"websocket-service/config"
	"websocket-service/pkg/logger"
)

type Error struct {
	log  logger.LoggerI
	name string
	port string
}

func NewError(log logger.LoggerI, name, port string) *Error {
	return &Error{
		log:  log,
		port: port,
		name: name,
	}
}

func (e *Error) Wrap(err *error, funcName string, req interface{}) {
	if *err == nil {
		return
	}
	*err = Wrap(*err, funcName)

	e.log.Error(msges(config.ErrorModel, e.name),
		logger.Error(*err),
		logger.Any("Service Port", e.port),
		logger.Any("request:", req),
	)

	*err = Wrap(Wrap(*err, e.name), config.ErrorModel)
}

func (e *Error) GetError(data interface{}) string {
	switch s := data.(type) {
	case string:
		errs := strings.Split(s, config.ErrorStyle)
		return errs[len(errs)-1]
	}
	return ""
}

func Wrap(err error, msg string) error {
	if err != nil {
		return errors.New(msg + err.Error())
	}
	return nil
}

func New(msg string) error {
	return errors.New(msg)
}

func WrapCheck(err *error, msg string) {
	if *err == nil {
		return
	}
	er := *err

	*err = errors.New(msg + config.ErrorStyle + er.Error())
}

func msges(msg1, msg2 string) string {
	return msg1 + config.ErrorStyle + msg2
}
