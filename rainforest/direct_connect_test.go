package rainforest

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestDirectConnect(t *testing.T) {
	mockResponse := `{"id":"my-tunnel-id","server_public_key":"server public key","server_port":12345,"server_endpoint":"server.example.com"}`
	testCases := []struct {
		id        string
		publicKey string
		expected  DirectConnectParams
	}{
		{
			id:        "my-tunnel-id",
			publicKey: "my-public-key",
			expected: DirectConnectParams{
				ID:        "my-tunnel-id",
				PublicKey: "my-public-key",
			},
		},
	}

	for _, testCase := range testCases {
		setup()
		defer cleanup()

		mux.HandleFunc("/direct_connect", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("Error parsing request body: %v", err.Error())
				}
				p := DirectConnectParams{}
				err = json.Unmarshal(body, &p)
				if err != nil {
					t.Errorf("Error unmarshalling request body: %v", err.Error())
				}

				if p != testCase.expected {
					t.Errorf("Unexpected request body. Want %v, Got %v", testCase.expected, p)
				}

				_, err = w.Write([]byte(mockResponse))
				if err != nil {
					t.Errorf("Error writing mock response: %v", err.Error())
				}
			} else {
				t.Errorf("Unexpected request method: %v", r.Method)
			}
		})

		serverInfo, err := client.SetupDirectConnectTunnel(testCase.id, testCase.publicKey)
		if err != nil {
			t.Error(err.Error())
		}

		if serverInfo.ID != testCase.id {
			t.Errorf("Correct ID not found in response. Want %v, Got %v", serverInfo.ID, testCase.id)
		}

		if serverInfo.ServerPublicKey != "server public key" {
			t.Errorf("Server public key not properly assigned to direct connect struct. Want %v, Got %v", serverInfo.ServerPublicKey, "server public key")
		}
		if serverInfo.ServerEndpoint != "server.example.com" {
			t.Errorf("Server endpoint not properly assigned to direct connect struct. Want %v, Got %v", serverInfo.ServerPublicKey, "server.example.com")
		}
		if serverInfo.ServerPort != 12345 {
			t.Errorf("Server port not properly assigned to direct connect struct. Want %v, Got %v", serverInfo.ServerPublicKey, 12345)
		}
	}
}
