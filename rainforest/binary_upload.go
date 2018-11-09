package rainforest

import (
	"os"
	"path/filepath"
)

// GetS3BinaryUrl gets a presigned s3 url for a customer to upload to
func (c *Client) UploadBinary(file *os.File) error {
	fileName := filepath.Base(file.Name())
	url, err := c.getS3BinaryUrl(fileName)
	if err != nil {
		return err
	}

	return postBinaryToS3(file, url)
}

func (c *Client) getS3BinaryUrl(fileName string) (string, error) {
	req, err := c.NewRequest("GET", "uploads?file_name="+fileName, nil)
	if err != nil {
		return "", err
	}

	var env string
	_, err = c.Do(req, &env)
	if err != nil {
		return "", err
	}

	return env, nil
}

func postBinaryToS3(file *os.File, url string) error {
	// TODO
	return nil
}
