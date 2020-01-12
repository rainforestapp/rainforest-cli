package rainforest

import (
	"strconv"
)

// EnvironmentParams are the parameters used to create a new Environment
type EnvironmentParams struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Environment represents an environment in Rainforest
type Environment struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	SiteEnvironments []SiteEnvironment `json:"site_environments"`
}

// CreateTemporaryEnvironment creates a new temporary environment and returns the
// Environment.
func (c *Client) CreateTemporaryEnvironment(urlString string) (*Environment, error) {
	body := EnvironmentParams{
		Name: "temporary-env-for-custom-url-via-CLI",
		URL:  urlString,
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

// SetSiteEnvironments uses a map of site id to urls
func (c *Client) SetSiteEnvironments(env *Environment, siteURLMap *map[int]string) (*Environment, error) {
	siteEnvs := env.SiteEnvironments
	for siteID, URL := range *siteURLMap {
		for i, env := range siteEnvs {
			if env.SiteID == siteID {
				err := c.setSiteEnvironmentURL(siteEnvs[i].ID, URL)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return env, nil
}

// DeleteEnvironment deletes an environment with a specified ID
func (c *Client) DeleteEnvironment(environmentID int) error {
	// Prepare request
	req, err := c.NewRequest("DELETE", "environments/"+strconv.Itoa(environmentID), nil)
	if err != nil {
		return err
	}

	// Send request and process response
	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
