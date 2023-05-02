package rainforest

import "fmt"

// EnvironmentParams are the parameters used to create a new Environment
type EnvironmentParams struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	IsTemporary    bool   `json:"is_temporary"`
	Webhook        string `json:"webhook,omitempty"`
	WebhookEnabled bool   `json:"webhook_enabled,omitempty"`
}

// Environment represents an environment in Rainforest
type Environment struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	IsTemporary    bool   `json:"is_temporary"`
	Webhook        string `json:"webhook"`
	WebhookEnabled bool   `json:"webhook_enabled"`
}

// CreateTemporaryEnvironment creates a new temporary environment and returns the
// Environment.
func (c *Client) CreateTemporaryEnvironment(runDescription string, urlString string, webhook string) (*Environment, error) {
	name := "temporary-env-for-custom-url-via-CLI"
	if runDescription != "" {
		if len(runDescription) > 241 {
			runDescription = runDescription[:241]
		}
		name = fmt.Sprintf("%v-temporary-env", runDescription)
	}

	body := EnvironmentParams{
		Name:           name,
		URL:            urlString,
		IsTemporary:    true,
		Webhook:        webhook,
		WebhookEnabled: webhook != "",
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
