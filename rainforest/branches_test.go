package rainforest

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestCreateBranch(t *testing.T) {
	testCases := []struct {
		branch       Branch
		expectedBody string
	}{
		{
			branch:       Branch{Name: "branch"},
			expectedBody: `{"name":"branch"}`,
		},
	}

	for _, testCase := range testCases {
		setup()
		defer cleanup()

		const requestMethod = "POST"
		mux.HandleFunc("/branches", func(w http.ResponseWriter, request *http.Request) {
			if request.Method != requestMethod {
				t.Errorf("Request method = %v, want %v", request.Method, requestMethod)
			}

			buffer := new(bytes.Buffer)
			buffer.ReadFrom(request.Body)
			response := strings.TrimSpace(buffer.String())
			if response != testCase.expectedBody {
				t.Errorf("Request body = %v, want %v", response, testCase.expectedBody)
			}
		})

		err := client.CreateBranch(&testCase.branch)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}
