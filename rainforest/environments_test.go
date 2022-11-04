package rainforest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateTemporaryEnvironment(t *testing.T) {
	testCases := []struct {
		runDescription string
		urlString      string
		envID          int
		expected       EnvironmentParams
	}{
		{
			runDescription: "",
			urlString:      "https://no-name.url",
			envID:          7331,
			expected: EnvironmentParams{
				Name:        "temporary-env-for-custom-url-via-CLI",
				URL:         "https://no-name.url",
				IsTemporary: true,
			},
		},
		{
			runDescription: "my-run-description",
			urlString:      "https://with-a-name.url",
			envID:          7332,
			expected: EnvironmentParams{
				Name:        "my-run-description-temporary-env",
				URL:         "https://with-a-name.url",
				IsTemporary: true,
			},
		},
		{
			runDescription: "My run with a giant description that goes on for over 255 characters, count them if you must. No seriously this is more than that. This won't fit in the environments table's name column, so we'll have to trim some off if we don't want this to loudly blow up.",
			urlString:      "https://with-a-name.url",
			envID:          7332,
			expected: EnvironmentParams{
				Name:        "My run with a giant description that goes on for over 255 characters, count them if you must. No seriously this is more than that. This won't fit in the environments table's name column, so we'll have to trim some off if we don't want this t-temporary-env",
				URL:         "https://with-a-name.url",
				IsTemporary: true,
			},
		},
	}

	for _, testCase := range testCases {
		setup()
		defer cleanup()

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

				if p != testCase.expected {
					t.Errorf("Unexpected request body. Want %v, Got %v", testCase.expected, p)
				}

				resJSON := fmt.Sprintf(`{"id":%v,"name":"%v","temporary":true}`, testCase.envID, p.Name)
				w.Write([]byte(resJSON))
			} else {
				t.Errorf("Unexpected request method: %v", r.Method)
			}
		})

		env, err := client.CreateTemporaryEnvironment(testCase.runDescription, testCase.urlString)
		if err != nil {
			t.Error(err.Error())
		}

		if env.ID != testCase.envID {
			t.Errorf("Correct ID not found in environment. Want %v, Got %v", testCase.envID, env.ID)
		}

		if env.Name != testCase.expected.Name {
			t.Errorf("Name not properly assigned to environment struct. Want %v, Got %v", testCase.expected.Name, env.Name)
		}
	}
}
