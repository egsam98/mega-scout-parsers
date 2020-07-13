package utils

import (
	"reflect"
)

type Set map[interface{}]interface{}

func NewSet() Set {
	return Set{}
}

func (s Set) Slice() []interface{} {
	values := make([]interface{}, 0, len(s))
	for _, v := range s {
		values = append(values, v)
	}
	return values
}

func (s Set) Add(data interface{}, keys ...string) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		hash := ""
		for _, key := range keys {
			hash += val.FieldByName(key).String() + "_"
		}
		s[hash] = data
		return
	}
	s[data] = data
}

func (s Set) Remove(data interface{}) {
	delete(s, data)
}

func (s Set) RemoveByStructFields(keys ...string) {
	hash := ""
	for _, key := range keys {
		hash += key + "_"
	}
	delete(s, hash)
}

func (s Set) Size() int {
	return len(s)
}
