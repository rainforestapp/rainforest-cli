package main

import (
	"math"
	"os"

	"encoding/csv"

	"errors"

	"strings"

	"github.com/urfave/cli"
)

// uploadTabularVar takes a path to csv file and creates tabular variable generator from it.
func uploadTabularVar(pathToCSV, name string, overwrite, singleUse bool) error {
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
	numOfBatches := int(math.Ceil(float64(len(rows)) / tabularBatchSize))
	rowsToUpload := make(chan [][]string, numOfBatches)
	for i := 0; i < len(rows); i += tabularBatchSize {
		batch := rows[i:min(i+tabularBatchSize, len(rows))]
		rowsToUpload <- batch
	}
	close(rowsToUpload)

	// spawn workers to upload the rows
	for i := 0; i < tabularConcurency; i++ {

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
//
// func rowUploadWorker(generator rainforest.Generator, columns []string, rows <-chan []string)

func csvUpload(c *cli.Context) error {
	// Get the csv file path either from the option or command argument
	return nil
}
