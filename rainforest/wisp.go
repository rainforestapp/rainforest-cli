package rainforest

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

// WispWriter writes test's id, title and wisp to a given file in json format.
type WispWriter struct {
	w *bufio.Writer
}

// WispReader reads a test's id, title and wisp from a formatted json file.
type WispReader struct {
	r *bufio.Reader
}

type WispJson struct {
	TestID int    `json:"id"`
	Title  string `json:"title"`
	Wisp   Wisp   `json:"wisp"`
}

// NewWispWriter returns Wisp writer based on passed io.Writer
func NewWispWriter(w io.Writer) *WispWriter {
	return &WispWriter{w: bufio.NewWriter(w)}
}

// NewWispReader returns reader based on passed io.Reader
func NewWispReader(r io.Reader) *WispReader {
	return &WispReader{r: bufio.NewReader(r)}
}

// WriteRFMLTest writes a given RFTest to its writer in the given RFML version.
func (r *WispWriter) WriteWispTest(test *RFTest) error {
	writer := r.w
	wispJson := WispJson{
		TestID: test.TestID,
		Title:  test.Title,
		Wisp:   *test.Wisp,
	}

	body, err := json.Marshal(wispJson)

	if err != nil {
		return err
	}

	_, err = writer.WriteString(string(body))

	// Writes buffered data to the underlying io.Writer
	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

// Read parses whole wisp json file
// and returns resulting RFTest
func (r *WispReader) Read() (*WispJson, error) {
	var wispJson WispJson
	err := json.NewDecoder(r.r).Decode(&wispJson)

	if err != nil {
		log.Printf("ERROR PARSING JSON: %v\n\n", err.Error())
		return nil, err
	}

	return &wispJson, nil
}
