package rainforest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestWriteWispTest(t *testing.T) {
	var buffer bytes.Buffer
	writer := NewWispWriter(&buffer)

	button := "left"
	elementID := 123
	seconds := 1
	hold := false
	visibility := false

	verbs := []Verb{
		{
			Action: "click",
			Button: &button,
			Target: &Noun{
				Type: "ui_element_reference",
				ID:   &elementID,
			},
			Hold:        &hold,
			HoldSeconds: &seconds,
		},
		{
			Action: "observe",
			Object: &Noun{
				Type: "ui_element_reference",
				ID:   &elementID,
			},
			Visibility: &visibility,
		},
	}

	wisp := Wisp{
		Version: "0.0.1",
		Verbs:   verbs,
	}

	wispJson := WispJson{
		TestID: 123,
		Title:  "my_title",
		Wisp:   wisp,
	}

	test := RFTest{
		TestID:  wispJson.TestID,
		Title:   wispJson.Title,
		Wisp:    &wispJson.Wisp,
		HasWisp: true,
	}

	getOutput := func() string {
		writer.WriteWispTest(&test)

		raw, err := ioutil.ReadAll(&buffer)

		if err != nil {
			t.Fatal(err.Error())
		}
		return string(raw)
	}

	output := getOutput()
	marshaledWispJson, _ := json.Marshal(wispJson)
	validJsonString := string(marshaledWispJson)

	if output != validJsonString {
		t.Errorf("Incorrect values for Wisp.\nGot %#v\nWant %#v", output, validJsonString)
	}
}

func TestRead(t *testing.T) {
	button := "left"
	elementID := 123
	seconds := 1
	hold := false
	visibility := false

	verbs := []Verb{
		{
			Action: "click",
			Button: &button,
			Target: &Noun{
				Type: "ui_element_reference",
				ID:   &elementID,
			},
			Hold:        &hold,
			HoldSeconds: &seconds,
		},
		{
			Action: "observe",
			Object: &Noun{
				Type: "ui_element_reference",
				ID:   &elementID,
			},
			Visibility: &visibility,
		},
	}

	wisp := Wisp{
		Version: "0.0.1",
		Verbs:   verbs,
	}

	validWispJson := WispJson{
		TestID: 123,
		Title:  "my_title",
		Wisp:   wisp,
	}

	testJson, _ := json.Marshal(validWispJson)

	r := strings.NewReader(string(testJson))
	reader := NewWispReader(r)
	output, err := reader.Read()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(*output, validWispJson) {
		t.Errorf("Incorrect values for Wisp json.\nGot %#v\nWant %#v", output, validWispJson)
	}
}
