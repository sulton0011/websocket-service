package errors

import "errors"

func Wrap(err error, msg string) error {
	if err != nil {
		return errors.New(msg + err.Error())
	}
	return nil
}
