package gometeor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type RMStruct1 struct {
	IntField    int
	StringField string
	FloatField  float64
	SliceField  []string
}

func TestRawMessage(t *testing.T) {
	js := `{"intfield" : 42, "stringfield" : "ohai!", "floatfield":137, "slicefield" : ["a", "x", "q"]}`
	var rawmsg RawMessage
	err := json.Unmarshal([]byte(js), &rawmsg)
	if err != nil {
		t.Errorf("Could not unmarshal to RawMessage: %v", err)
	}
	var rms RMStruct1
	err = rawmsg.Decode(&rms)
	if err != nil {
		t.Errorf("Could not unmarshal RawMessage to object %v", err)
	} else {
		if rms.IntField != 42 ||
			rms.StringField != "ohai!" ||
			rms.FloatField != 137.0 ||
			!reflect.DeepEqual(rms.SliceField, []string{"a", "x", "q"}) {
			t.Errorf("Decoded object was %v", rms)
		}

	}

	buf, err := ToJSON(rms)
	if err != nil {
		t.Errorf("Error in ToJSON: %v", err)
	}

	spaces := regexp.MustCompile(`\s+`)
	nospaces := string(spaces.ReplaceAll(buf, nil))
	if nospaces != strings.Replace(js, " ", "", -1) {
		t.Errorf("ToJSON failure: %s", nospaces)
	}
}

func TestCall(t *testing.T) {
	fn := func(a string, b int, c []string) string {
		return fmt.Sprintf("%s %d: %v", a, b, c)
	}
	args := []json.RawMessage{
		json.RawMessage(`"There are"`),
		json.RawMessage(`3`),
		json.RawMessage(`["a","b","c"]`),
	}
	result, err := Call(fn, args)
	if err != nil {
		t.Errorf("Call returned error %v", err)
	}
	str, ok := result.(string)
	if !ok {
		t.Errorf("Call result was a %T", str)
	}
	if str != `There are 3: [a b c]` {
		t.Errorf("Call result was %T %v", result, result)
	}
	fn2 := func(a interface {
		Foo()
	},) {
		a.Foo()
	}
	args2 := []json.RawMessage{
		json.RawMessage(`"An argument"`),
	}
	result, err = Call(fn2, args2)
	if err == nil {
		t.Errorf("Bad call returned no error")
	}
	args2 = []json.RawMessage{
		json.RawMessage(`null`),
	}
	result, err = Call(fn2, args2)
	if err == nil {
		t.Errorf("Bad call returned no error")
	}
}
