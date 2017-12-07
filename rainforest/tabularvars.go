package rainforest

import (
	"errors"
	"strconv"
)

// Generator is a type representing generators which can be used as variables in RF tests.
// They can be builtin or uploaded by customer as a tabular variable.
type Generator struct {
	ID          int               `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"generator_type,omitempty"`
	SingleUse   bool              `json:"single_use,omitempty"`
	Columns     []GeneratorColumn `json:"columns,omitempty"`
	RowCount    int               `json:"row_count,omitempty"`
}

// GetID returns the Generator name
func (g Generator) GetID() string {
	return g.Name
}

// GetDescription returns the Generator's description
func (g Generator) GetDescription() string {
	return g.Description
}

// GeneratorColumn is a type of column in a generator
type GeneratorColumn struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
}

// GeneratorRelatedTests is a type which holds tests where the generator has been used
type GeneratorRelatedTests struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

// GetGenerators fetches a list of all available generators for the account
func (c *Client) GetGenerators() ([]Generator, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "generators", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var generatorResp []Generator
	_, err = c.Do(req, &generatorResp)
	if err != nil {
		return nil, err
	}
	return generatorResp, nil
}

// DeleteGenerator deletes generator with specified ID
func (c *Client) DeleteGenerator(genID int) error {
	// Prepare request
	req, err := c.NewRequest("DELETE", "generators/"+strconv.Itoa(genID), nil)
	if err != nil {
		return err
	}

	// Fire away the request
	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	// And just nil when we're done
	return nil
}

// CreateTabularVar creates new tabular variable on RF and returns Generator associated with it
// columns argument should contain just an array of column names, contents of the generator should
// be filled using AddGeneratorRows.
func (c *Client) CreateTabularVar(name, description string,
	columns []string, singleUse bool) (*Generator, error) {
	//Prepare request
	type genCreateRequest struct {
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		SingleUse   bool     `json:"single_use,omitempty"`
		Columns     []string `json:"columns,omitempty"`
	}
	body := genCreateRequest{name, description, singleUse, columns}
	req, err := c.NewRequest("POST", "generators", body)
	if err != nil {
		return &Generator{}, err
	}

	var out *Generator
	_, err = c.Do(req, &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

// AddGeneratorRows adds rows to the specified tabular variable
// rowData is in a form of [{ 123: "foo", 124: "bar" }, { 123: "baz", 124: "qux" }] where 123 is a column ID
// AddGeneratorRowsFromTable is also provided which accepts different rowData format
func (c *Client) AddGeneratorRows(targetGenerator *Generator, rowData []map[int]string) error {
	//Prepare request
	type batchRowsRequest struct {
		RowData []map[int]string `json:"data,omitempty"`
	}
	body := batchRowsRequest{rowData}
	reqURL := "generators/" + strconv.Itoa(targetGenerator.ID) + "/rows/batch"
	req, err := c.NewRequest("POST", reqURL, body)
	if err != nil {
		return err
	}

	// Send out the batch
	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// AddGeneratorRowsFromTable adds rows to the specified tabular variable
// data should be formatted as follows:
// targetColumns contains names of existing columns to which add data e.g. ["login", "password"]
// rowData contains row data in columns order specified in targetColumns e.g. [["foo", "bar"], ["baz", "qux"]]
func (c *Client) AddGeneratorRowsFromTable(targetGenerator *Generator,
	targetColumns []string, rowData [][]string) error {
	// Quick sanity check of the args
	if len(targetColumns) != len(targetGenerator.Columns) {
		return errors.New("Invalid number of columns for given generator.")
	}

	// Create a mapping between column names and their id for the Generator
	colNameToID := make(map[string]int)
	for _, column := range targetGenerator.Columns {
		colNameToID[column.Name] = column.ID
	}

	// Use created mapping to create an array of column IDs
	targetColumnIDs := make([]int, len(targetColumns))
	for i, colName := range targetColumns {
		id, ok := colNameToID[colName]
		if !ok {
			return errors.New("Invalid column name for given generator.")
		}
		targetColumnIDs[i] = id
	}

	// Use above helper array to create properly formatted row_data array
	formattedRowData := make([]map[int]string, len(rowData))
	for rowIndex, row := range rowData {
		formattedRow := make(map[int]string)
		if len(row) != len(targetColumns) {
			return errors.New("Invalid number of columns in a row" + strconv.Itoa(rowIndex) +
				"for given generator.")
		}
		for colIndex, colValue := range row {
			formattedRow[targetColumnIDs[colIndex]] = colValue
		}
		formattedRowData[rowIndex] = formattedRow
	}

	// Call the function that will add the rows.
	return c.AddGeneratorRows(targetGenerator, formattedRowData)
}
