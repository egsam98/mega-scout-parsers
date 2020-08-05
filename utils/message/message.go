package message

import errors2 "github.com/go-errors/errors"

type Message struct {
	Data  interface{}
	Error *errors2.Error
}

func (m Message) IsError() bool {
	return m.Error != nil
}

func Error(err *errors2.Error) Message {
	return Message{
		Data:  nil,
		Error: err,
	}
}

func Ok(data interface{}) Message {
	return Message{
		Data:  data,
		Error: nil,
	}
}

func Nil() Message {
	return Message{}
}
