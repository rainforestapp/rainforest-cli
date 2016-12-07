package rainforest

import "strconv"

// UploadedFile represents a file that has been uploaded to Rainforest
type UploadedFile struct {
	ID        int    `json:"id"`
	Signature string `json:"signature"`
	Digest    string `json:"digest"`
}

// GetUploadedFiles returns information for all all files uploaded to the
// given test before.
func (c *Client) GetUploadedFiles(fileID int) ([]UploadedFile, error) {
	req, err := c.NewRequest("GET", "tests/"+strconv.Itoa(fileID)+"/files", nil)
	if err != nil {
		return nil, err
	}

	var fileResp []UploadedFile
	_, err = c.Do(req, &fileResp)
	if err != nil {
		return nil, err
	}
	return fileResp, nil
}
