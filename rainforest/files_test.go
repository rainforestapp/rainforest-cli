package rainforest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetUploadedFiles(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	files := []UploadedFile{
		{ID: 123, Signature: "file_sig1", Digest: "digest1"},
		{ID: 456, Signature: "file_sig2", Digest: "digest2"},
	}

	mux.HandleFunc("/tests/rfml_ids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		enc := json.NewEncoder(w)
		enc.Encode(files)
	})

	out, _ := client.GetRFMLIDs()

	if !reflect.DeepEqual(files, out) {
		t.Errorf("Response expected = %v, actual %v", files, out)
	}
}
