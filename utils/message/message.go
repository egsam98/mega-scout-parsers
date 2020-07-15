package message

type Message struct {
	Data  interface{}
	Error error
}

func (m Message) IsError() bool {
	return m.Error != nil
}

func Error(err error) Message {
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
