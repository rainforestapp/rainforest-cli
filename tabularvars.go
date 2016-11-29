package main

import (
	"math"
	"os"

	"encoding/csv"

	"errors"

	"strings"

	"log"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// tabularVariablesAPI is part of the API connected to the tabular variables
type tabularVariablesAPI interface {
	GetGenerators() ([]rainforest.Generator, error)
	DeleteGenerator(genID int) error
	CreateTabularVar(name, description string,
		columns []string, singleUse bool) (rainforest.Generator, error)
	AddGeneratorRowsFromTable(targetGenerator rainforest.Generator,
		targetColumns []string, rowData [][]string) error
}

// uploadTabularVar takes a path to csv file and creates tabular variable generator from it.
func uploadTabularVar(api tabularVariablesAPI, pathToCSV, name string, overwrite, singleUse bool) error {
	// Open up the CSV file and parse it, return early with an error if we fail to get to the file
	f, err := os.Open(pathToCSV)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Check if the variable exists in RF
	var existingGenID int
	generators, err := api.GetGenerators()
	if err != nil {
		return err
	}

	for _, gen := range generators {
		if gen.Name == name {
			existingGenID = gen.ID
		}
	}

	if existingGenID != 0 {
		if overwrite {
			// if variable exists and we want to override it with new one delete it here
			log.Printf("Tabular var %v exists, overwriting it with new data.\n", name)
			err := api.DeleteGenerator(existingGenID)
			if err != nil {
				return err
			}
		} else {
			// if variable exists but we didn't specify to override it then return with error
			return errors.New("Tabular variable: " + name +
				" already exists, use different name or choose an option to override it")
		}
	}

	// prepare input data
	columnNames, rows := records[0], records[1:]
	parsedColumnNames := make([]string, len(columnNames))
	for i, colName := range columnNames {
		formattedColName := strings.TrimSpace(strings.ToLower(colName))
		parsedColName := strings.Replace(formattedColName, " ", "_", -1)
		parsedColumnNames[i] = parsedColName
	}

	// create new generator for the tabular variable
	description := "Variable " + name + " uploded through cli client."
	newGenerator, err := api.CreateTabularVar(name, description, parsedColumnNames, singleUse)
	if err != nil {
		return err
	}

	// batch the rows and put them into a channel
	numOfBatches := int(math.Ceil(float64(len(rows)) / float64(tabularBatchSize)))
	rowsToUpload := make(chan [][]string, numOfBatches)
	for i := 0; i < len(rows); i += tabularBatchSize {
		batch := rows[i:min(i+tabularBatchSize, len(rows))]
		rowsToUpload <- batch
	}
	close(rowsToUpload)

	// chan to gather errors from workers
	errors := make(chan error, numOfBatches)

	log.Println("Beginning batch upload of csv file...")

	// spawn workers to upload the rows
	for i := 0; i < tabularConcurrency; i++ {
		go rowUploadWorker(api, newGenerator, parsedColumnNames, rowsToUpload, errors)
	}

	for i := 0; i < numOfBatches; i++ {
		if err := <-errors; err != nil {
			return err
		}
		log.Printf("Tabular variable '%v' batch %v of %v uploaded.", name, i+1, numOfBatches)
	}

	return nil
}

// Well... yeah...
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// rowUploadWorker is a helper worker which reads batch of rows to upload from rows chan
// and pushes potential errors through errorsChan
func rowUploadWorker(api tabularVariablesAPI, generator rainforest.Generator,
	columns []string, rowsChan <-chan [][]string, errorsChan chan<- error) {
	for rows := range rowsChan {
		error := api.AddGeneratorRowsFromTable(generator, columns, rows)
		errorsChan <- error
	}
}

// csvUpload is a wrapper around uploadTabularVar to function with csv-upload cli command
func csvUpload(c cliContext, api tabularVariablesAPI) error {
	// Get the csv file path either from the option or command argument
	filePath := c.Args().First()
	if filePath == "" {
		filePath = c.String("csv-file")
	}
	if filePath == "" {
		return cli.NewExitError("CSV filename not specified", 1)
	}
	name := c.String("name")
	if name == "" {
		return cli.NewExitError("Tabular variable name not specified", 1)
	}
	overwrite := c.Bool("overwrite-variable")
	singleUse := c.Bool("single-use")

	err := uploadTabularVar(api, filePath, name, overwrite, singleUse)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

// preRunCSVUpload is a wrapper around uploadTabularVar to be ran before starting a new run
func preRunCSVUpload(c cliContext, api tabularVariablesAPI) error {
	// Get the csv file path either and skip uploading if it's not present
	filePath := c.String("import-variable-csv-file")
	if filePath == "" {
		return nil
	}
	name := c.String("import-variable-name")
	if name == "" {
		return errors.New("Tabular variable name not specified")
	}
	overwrite := c.Bool("overwrite-variable")
	singleUse := c.Bool("single-use")

	return uploadTabularVar(api, filePath, name, overwrite, singleUse)
}
