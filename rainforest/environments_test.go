package rainforest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateTemporaryEnvironment(t *testing.T) {
	setup()
	defer cleanup()

	expectedID := 7331

	mux.HandleFunc("/environments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Error parsing request body: %v", err.Error())
			}
			p := EnvironmentParams{}
			err = json.Unmarshal(body, &p)
			if err != nil {
				t.Errorf("Error unmarshalling request body: %v", err.Error())
			}

			resJSON := fmt.Sprintf(`{"id":%v,"name":"%v"}`, expectedID, p.Name)
			w.Write([]byte(resJSON))
		} else {
			t.Errorf("Unexpected request method: %v", r.Method)
		}
	})

	env, err := client.CreateTemporaryEnvironment("https://www.rainforestqa.com/")
	if err != nil {
		t.Error(err.Error())
	}

	if env.ID != expectedID {
		t.Errorf("Correct ID not found in environment. Want %v, Got %v", expectedID, env.ID)
	}

	if env.Name == "" {
		t.Error("Name not properly assigned to environment struct. Got empty string.")
	}
}
