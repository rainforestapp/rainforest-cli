package rainforest

// EnvironmentParams are the parameters used to create a new Environment
type EnvironmentParams struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	IsTemporary bool   `json:"is_temporary"`
}

// Environment represents an environment in Rainforest
type Environment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsTemporary bool   `json:"is_temporary"`
}

// CreateTemporaryEnvironment creates a new temporary environment and returns the
// Environment.
func (c *Client) CreateTemporaryEnvironment(urlString string) (*Environment, error) {
	body := EnvironmentParams{
		Name:        "temporary-env-for-custom-url-via-CLI",
		URL:         urlString,
		IsTemporary: true,
	}
	req, err := c.NewRequest("POST", "environments", &body)
	if err != nil {
		return nil, err
	}

	var env Environment
	_, err = c.Do(req, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
