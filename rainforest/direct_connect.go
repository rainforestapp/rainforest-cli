package rainforest

// DirectConnectParams are the parameters the Rainforest Direct Connect API expects
type DirectConnectParams struct {
	PublicKey string `json:"public_key"`
	ID        string `json:"id"`
}

// DirectConnectConnection is the response from Rainforest's Direct Connect API telling us how to setup the wireguard tunnel
type DirectConnectConnection struct {
	ID              string `json:"id"`
	ServerPublicKey string `json:"server_public_key"`
	ServerPort      int    `json:"server_port"`
	ServerEndpoint  string `json:"server_endpoint"`
}

// SetupDirectConnectTunnel tells Rainforest to set up a Direct Connect tunnel for us
func (c *Client) SetupDirectConnectTunnel(id string, publicKey string) (*DirectConnectConnection, error) {

	body := DirectConnectParams{
		ID:        id,
		PublicKey: publicKey,
	}
	req, err := c.NewRequest("POST", "direct_connect", &body)
	if err != nil {
		return nil, err
	}

	var connection DirectConnectConnection
	_, err = c.Do(req, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}
