package rainforest

import (
	"path/filepath"
	"github.com/urfave/cli"
)

// getUploadEndpoint gets a presigned s3 url for a customer to upload to
func (c *Client) GetUploadEndpoint(fileName string) (*string, error) {
	var env string

	name := filepath.Base(fileName)

	req, err := c.NewRequest("GET", "uploads?file_name=" + name, nil)
	if err != nil {
		return &env, err
	}

	_, err = c.Do(req, &env)
	if err != nil {
		return &env, err
	}

	return &env, nil
}
