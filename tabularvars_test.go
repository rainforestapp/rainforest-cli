package main

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"

	"log"

	"strings"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

const (
	fakeCSVPath = "./fake.csv"
)

func TestMin(t *testing.T) {
	var testCases = []struct {
		a        int
		b        int
		expected int
	}{
		{
			a:        1,
			b:        3,
			expected: 1,
		},
		{
			a:        100,
			b:        23,
			expected: 23,
		},
		{
			a:        42,
			b:        42,
			expected: 42,
		},
	}

	for _, testCase := range testCases {
		res := min(testCase.a, testCase.b)

		if res != testCase.expected {
			t.Errorf("Incorrect value for min(%v, %v) = %v, expected: %v",
				testCase.a, testCase.b, res, testCase.expected)
		}
	}
}

type fakeAPI struct {
	getGenerators    func() ([]rainforest.Generator, error)
	deleteGenerator  func(genID int) error
	createTabularVar func(name, description string,
		columns []string, singleUse bool) (*rainforest.Generator, error)
	addGeneratorRowsFromTable func(targetGenerator *rainforest.Generator,
		targetColumns []string, rowData [][]string) error
}

func (f fakeAPI) GetGenerators() ([]rainforest.Generator, error) {
	if f.getGenerators != nil {
		return f.getGenerators()
	}
	return nil, nil
}

func (f fakeAPI) DeleteGenerator(genID int) error {
	if f.deleteGenerator != nil {
		return f.deleteGenerator(genID)
	}
	return nil
}

func (f fakeAPI) CreateTabularVar(name, description string,
	columns []string, singleUse bool) (*rainforest.Generator, error) {
	if f.createTabularVar != nil {
		return f.createTabularVar(name, description, columns, singleUse)
	}
	return &rainforest.Generator{}, nil
}

func (f fakeAPI) AddGeneratorRowsFromTable(targetGenerator *rainforest.Generator,
	targetColumns []string, rowData [][]string) error {
	if f.addGeneratorRowsFromTable != nil {
		return f.addGeneratorRowsFromTable(targetGenerator, targetColumns, rowData)
	}
	return nil
}

func TestRowUploadWorker(t *testing.T) {
	const testBatchesCount = 2
	gen := &rainforest.Generator{ID: 123}
	cols := []string{"test", "columns"}
	inChan := make(chan [][]string, testBatchesCount)
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	inChan <- firstBatch
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	inChan <- secondBatch
	close(inChan)
	errorsChan := make(chan error, testBatchesCount)
	var callCount int
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, gen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, gen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount = callCount + 1
			return nil
		},
	}
	rowUploadWorker(f, gen, cols, inChan, errorsChan)
	if callCount != testBatchesCount {
		t.Errorf("Api called wrong number of times. Called: %v, expected: %v", callCount, testBatchesCount)
	}
	err := <-errorsChan
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	err = <-errorsChan
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
}

func TestRowUploadWorker_Error(t *testing.T) {
	const testBatchesCount = 2
	gen := &rainforest.Generator{ID: 123}
	cols := []string{"test", "columns"}
	inChan := make(chan [][]string, testBatchesCount)
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	inChan <- firstBatch
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	inChan <- secondBatch
	close(inChan)
	errorsChan := make(chan error, testBatchesCount)
	var callCount int
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, gen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, gen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount = callCount + 1
			return errors.New("PICNIC!!1")
		},
	}
	rowUploadWorker(f, gen, cols, inChan, errorsChan)
	if callCount != testBatchesCount {
		t.Errorf("Api called wrong number of times. Called: %v, expected: %v", callCount, testBatchesCount)
	}
	err := <-errorsChan
	if err == nil {
		t.Errorf("Expected error, instead got nil.")
	}
}

func createValidFakeCSV(t *testing.T) {
	f, err := os.Create(fakeCSVPath)
	if err != nil {
		t.Fatal("Couldn't create the test csv file.")
	}
	defer f.Close()

	f.WriteString("test,columns\nfoo,bar\nbaz,wut\nqwe,asd\nzxc,jkl\n")
}

func deleteFakeCSV(t *testing.T) {
	err := os.Remove(fakeCSVPath)
	if err != nil {
		t.Fatal("Couldn't remove the test csv file.")
	}
}

func TestUploadTabularVar(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableOverwrite := false
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:   123,
		Name: variableName,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}

	// Capture output
	var outBuffer bytes.Buffer
	log.SetOutput(&outBuffer)
	defer log.SetOutput(os.Stdout)

	err := uploadTabularVar(f, fakeCSVPath, variableName, variableOverwrite, variableSingleUse)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 1; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 2; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}

	if !strings.Contains(outBuffer.String(), "uploaded") {
		t.Errorf("Expected progress updates logged out (containing 'uploaded'), got %v", outBuffer.String())
	}
}

func TestUploadTabularVar_Exists_NoOverwrite(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableOverwrite := false
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:   123,
		Name: variableName,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
				fakeNewGen,
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}
	err := uploadTabularVar(f, fakeCSVPath, variableName, variableOverwrite, variableSingleUse)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if expected := 1; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
}

func TestUploadTabularVar_Exists_Overwrite(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableDescription := "testVar description"
	variableOverwrite := true
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:          123,
		Name:        variableName,
		Description: variableDescription,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
				fakeNewGen,
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if variableDescription != description {
				t.Errorf("Incorrect value of description passed to newTabVar. Got: %v, expected: %v", description, variableDescription)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}

	// Capture output
	var outBuffer bytes.Buffer
	log.SetOutput(&outBuffer)
	defer log.SetOutput(os.Stdout)

	err := uploadTabularVar(f, fakeCSVPath, variableName, variableOverwrite, variableSingleUse)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 1; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 2; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}

	if !strings.Contains(outBuffer.String(), "uploaded") {
		t.Errorf("Expected progress updates logged out (containing 'uploaded'), got %v", outBuffer.String())
	}
	if !strings.Contains(outBuffer.String(), "overwriting") {
		t.Errorf("Expected information about overwriting variable, got %v", outBuffer.String())
	}
}

func TestCSVUpload(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableOverwrite := true
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:   123,
		Name: variableName,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
				fakeNewGen,
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}

	// Fake cli args
	fakeContext := newFakeContext(map[string]interface{}{
		"name":               variableName,
		"overwrite-variable": variableOverwrite,
		"single-use":         variableSingleUse,
	}, cli.Args{fakeCSVPath})

	// Capture output
	var outBuffer bytes.Buffer
	log.SetOutput(&outBuffer)
	defer log.SetOutput(os.Stdout)

	err := csvUpload(fakeContext, f)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 1; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 2; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}

	if !strings.Contains(outBuffer.String(), "uploaded") {
		t.Errorf("Expected progress updates logged out (containing 'uploaded'), got %v", outBuffer.String())
	}
	if !strings.Contains(outBuffer.String(), "overwriting") {
		t.Errorf("Expected information about overwriting variable, got %v", outBuffer.String())
	}
}

func TestCSVUpload_MissingName(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableOverwrite := true
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:   123,
		Name: variableName,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
				fakeNewGen,
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}

	// Fake cli args
	fakeContext := newFakeContext(map[string]interface{}{
		"overwrite-variable": variableOverwrite,
		"single-use":         variableSingleUse,
	}, cli.Args{fakeCSVPath})

	// Capture output
	var outBuffer bytes.Buffer
	log.SetOutput(&outBuffer)
	defer log.SetOutput(os.Stdout)

	err := csvUpload(fakeContext, f)
	if err == nil {
		t.Errorf("Expected error, instead got nil.")
	}
	if expected := 0; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
}

func TestPreRunCSVUpload(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableOverwrite := true
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:   123,
		Name: variableName,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
				fakeNewGen,
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}

	// Fake cli args
	fakeContext := newFakeContext(map[string]interface{}{
		"import-variable-csv-file": fakeCSVPath,
		"import-variable-name":     variableName,
		"overwrite-variable":       variableOverwrite,
		"single-use":               variableSingleUse,
	}, cli.Args{fakeCSVPath})

	// Capture output
	var outBuffer bytes.Buffer
	log.SetOutput(&outBuffer)
	defer log.SetOutput(os.Stdout)

	err := preRunCSVUpload(fakeContext, f)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 1; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 1; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 2; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}

	if !strings.Contains(outBuffer.String(), "uploaded") {
		t.Errorf("Expected progress updates logged out (containing 'uploaded'), got %v", outBuffer.String())
	}
	if !strings.Contains(outBuffer.String(), "overwriting") {
		t.Errorf("Expected information about overwriting variable, got %v", outBuffer.String())
	}
}

func TestPreRunCSVUpload_MissingName(t *testing.T) {
	// Setup testing  stuff
	createValidFakeCSV(t)
	defer deleteFakeCSV(t)
	tabularBatchSize = 2
	variableName := "testVar"
	variableOverwrite := true
	variableSingleUse := false
	// Data from csv
	cols := []string{"test", "columns"}
	firstBatch := [][]string{{"foo", "bar"}, {"baz", "wut"}}
	secondBatch := [][]string{{"qwe", "asd"}, {"zxc", "jkl"}}
	// Fake responses Data
	fakeNewGen := rainforest.Generator{
		ID:   123,
		Name: variableName,
		Columns: []rainforest.GeneratorColumn{
			{
				ID:   456,
				Name: "test",
			},
			{
				ID:   789,
				Name: "columns",
			},
		},
	}
	callCount := make(map[string]int)
	f := fakeAPI{
		addGeneratorRowsFromTable: func(targetGenerator *rainforest.Generator,
			targetColumns []string, rowData [][]string) error {
			if !reflect.DeepEqual(targetGenerator, &fakeNewGen) {
				t.Errorf("Incorrect value of targetGenerator passed. Got: %v, expected: %v", targetGenerator, fakeNewGen)
			}
			if !reflect.DeepEqual(targetColumns, cols) {
				t.Errorf("Incorrect value of targetColumns passed. Got: %v, expected: %v", targetColumns, cols)
			}
			if callCount["addGeneratorRowsFromTable"] == 0 && !reflect.DeepEqual(rowData, firstBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, firstBatch)
			}
			if callCount["addGeneratorRowsFromTable"] == 1 && !reflect.DeepEqual(rowData, secondBatch) {
				t.Errorf("Incorrect value of rowData passed. Got: %v, expected: %v", rowData, secondBatch)
			}
			callCount["addGeneratorRowsFromTable"] = callCount["addGeneratorRowsFromTable"] + 1
			return nil
		},
		getGenerators: func() ([]rainforest.Generator, error) {
			callCount["getGenerators"] = callCount["getGenerators"] + 1
			return []rainforest.Generator{
				{
					ID:   42,
					Name: "Kappa",
				},
				{
					ID:   1337,
					Name: "EleGiggle",
				},
				fakeNewGen,
			}, nil
		},
		deleteGenerator: func(genID int) error {
			if genID != fakeNewGen.ID {
				t.Errorf("Incorrect value of genID passed to delete. Got: %v, expected: %v", genID, fakeNewGen.ID)
			}
			callCount["deleteGenerator"] = callCount["deleteGenerator"] + 1
			return nil
		},
		createTabularVar: func(name, description string, columns []string,
			singleUse bool) (*rainforest.Generator, error) {
			if variableName != name {
				t.Errorf("Incorrect value of name passed to newTabVar. Got: %v, expected: %v", name, variableName)
			}
			if !reflect.DeepEqual(columns, cols) {
				t.Errorf("Incorrect value of columns passed to newTabVar. Got: %v, expected: %v", columns, cols)
			}
			if variableSingleUse != singleUse {
				t.Errorf("Incorrect value of singleUse passed to newTabVar. Got: %v, expected: %v", singleUse, variableSingleUse)
			}
			callCount["createTabularVar"] = callCount["createTabularVar"] + 1
			return &fakeNewGen, nil
		},
	}

	// Fake cli args
	fakeContext := newFakeContext(map[string]interface{}{
		"overwrite-variable": variableOverwrite,
		"single-use":         variableSingleUse,
	}, cli.Args{fakeCSVPath})

	// Capture output
	var outBuffer bytes.Buffer
	log.SetOutput(&outBuffer)
	defer log.SetOutput(os.Stdout)

	err := preRunCSVUpload(fakeContext, f)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 0; callCount["getGenerators"] != expected {
		t.Errorf("api.getGenerators called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["deleteGenerator"] != expected {
		t.Errorf("api.deleteGenerator called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["createTabularVar"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
	if expected := 0; callCount["addGeneratorRowsFromTable"] != expected {
		t.Errorf("api.createTabularVar called invalid number of times: %v, expected %v", callCount["getGenerators"], expected)
	}
}
