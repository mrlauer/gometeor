package meteor

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type RawMessage map[string]json.RawMessage

func _removeTrailingNewline(buf *bytes.Buffer) {
	l := buf.Len()
	if l > 0 && buf.Bytes()[l-1] == '\n' {
		buf.Truncate(l - 1)
	}
}

// ToJSON encodes a struct or map to JSON as the meteor protocol wants it.
func ToJSON(obj interface{}) ([]byte, error) {
	value := reflect.ValueOf(obj)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	bufw := bytes.NewBuffer(nil)

	if value.Kind() == reflect.Struct {
		typ := value.Type()
		nfield := value.NumField()
		bufw.WriteString("{")
		for i := 0; i < nfield; i++ {
			ftyp := typ.Field(i)
			f := value.Field(i)
			// TODO: json tags.
			if i != 0 {
				bufw.WriteString(",")
			}
			name := strings.ToLower(ftyp.Name)
			err := json.NewEncoder(bufw).Encode(name)
			_removeTrailingNewline(bufw)
			if err != nil {
				return nil, err
			}
			bufw.WriteString(":")
			err = json.NewEncoder(bufw).Encode(f.Interface())
			_removeTrailingNewline(bufw)
			if err != nil {
				return nil, err
			}
		}
		bufw.WriteString("}")
	}
	return bufw.Bytes(), nil
}

func (m RawMessage) Decode(obj interface{}) error {
	// obj had better be a pointer to a structure.
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return errors.New("Invalid receiver")
	}
	if value.IsNil() {
		return errors.New("Invalid receiver")
	}
	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("Invalid receiver")
	}

	// For each field in the struct
	typ := elem.Type()
	nfield := elem.NumField()
	for i := 0; i < nfield; i++ {
		ftyp := typ.Field(i)
		f := elem.Field(i)
		// TODO: json tags.
		name := ftyp.Name
		raw, ok := m[name]
		if !ok {
			raw, ok = m[strings.ToLower(name)]
		}
		if ok {
			err := json.Unmarshal(raw, f.Addr().Interface())
			if err != nil {
				// TODO: um, something else.
				return err
			}
		}
	}
	return nil
}

// OK, maybe this belongs in its own file.

// Call the given function with the given arguments.
// fn must be a function that returns a single result,
// and args must be unmarshallable into the argument types of the functions.
func Call(fn interface{}, args []json.RawMessage) (interface{}, error) {
	fnvalue := reflect.ValueOf(fn)
	if fnvalue.Kind() != reflect.Func {
		return nil, errors.New("Non-function passed to call")
	}
	fntyp := fnvalue.Type()
	nin := fntyp.NumIn()
	argvals := make([]reflect.Value, nin)
	if nin != len(args) {
		return nil, errors.New("Wrong number of arguments passed to call")
	}
	if fntyp.NumOut() > 1 {
		return nil, errors.New("Function must return at most one result")
	}
	for i := 0; i < nin; i++ {
		argtyp := fntyp.In(i)
		argval := reflect.New(argtyp)
		err := json.Unmarshal(args[i], argval.Interface())
		if err != nil {
			return nil, errors.New("Could not unmarshall parameters")
		}
		argvals[i] = argval.Elem()
	}
	results := fnvalue.Call(argvals)
	if len(results) > 0 {
		return results[0].Interface(), nil
	}
	return nil, nil
}
