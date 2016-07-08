package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestParseArguments(t *testing.T) {
	var testCases = []struct {
		input        []string
		wantCommands []string
		wantFlags    []string
	}{
		{
			input:        []string{"foo", "bar", "-words=baz"},
			wantCommands: []string{"bar"},
			wantFlags:    []string{"-words=baz"},
		},
		{
			input:        []string{"foo", "-words=baz", "bar"},
			wantCommands: []string{"bar"},
			wantFlags:    []string{"-words=baz"},
		},
		{
			input:        []string{"foo", "-numbers=321", "bar", "-words=baz"},
			wantCommands: []string{"bar"},
			wantFlags:    []string{"-numbers=321", "-words=baz"},
		},
		{
			input:        []string{"foo", "-numbers=321", "-words=baz"},
			wantCommands: nil,
			wantFlags:    []string{"-numbers=321", "-words=baz"},
		},
		{
			input:        []string{"foo"},
			wantCommands: nil,
			wantFlags:    nil,
		},
		{
			input:        []string{"foo", "bar", "wow"},
			wantCommands: []string{"bar", "wow"},
			wantFlags:    nil,
		},
	}
	tempOsArgs := os.Args
	for _, test := range testCases {
		os.Args = test.input
		gotCommands, gotFlags := parseArgs(os.Args)
		if !reflect.DeepEqual(test.wantCommands, gotCommands) {
			t.Errorf("Command incorrect. Expected %v, got %v", test.wantCommands, gotCommands)
		}
		if !reflect.DeepEqual(test.wantFlags, gotFlags) {
			t.Errorf("Flag incorrect. Expected %v, got %v", test.wantFlags, gotFlags)
		}
	}
	os.Args = tempOsArgs
}

func TestApiToken(t *testing.T) {
	realAPIToken := apiToken
	realOSArgs := os.Args
	realCommandLine := flag.CommandLine
	realEnvToken := os.Getenv("RAINFOREST_API_TOKEN")

	defaultOsArgs := []string{"rainforest-cli", "run"}
	testCases := []struct {
		envToken      string
		osArgs        []string
		expectedToken string
	}{
		{
			osArgs:        defaultOsArgs,
			envToken:      "",
			expectedToken: "",
		},
		{
			osArgs:        []string{"rainforest-cli", "run", "-token=flag_token"},
			envToken:      "",
			expectedToken: "flag_token",
		},
		{
			osArgs:        defaultOsArgs,
			envToken:      "env_token",
			expectedToken: "env_token",
		},
		{
			osArgs:        []string{"rainforest-cli", "run", "-token=flag_token"},
			envToken:      "env_token",
			expectedToken: "flag_token",
		},
	}

	for _, test := range testCases {
		apiToken = ""
		os.Setenv("RAINFOREST_API_TOKEN", test.envToken)
		os.Args = test.osArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()

		if apiToken != test.expectedToken {
			t.Logf("os.Args = %v", os.Args)
			t.Logf(`RAINFOREST_API_TOKEN = "%v"`, os.Getenv("RAINFOREST_API_TOKEN"))
			t.Errorf("Wrong flag detected. Expected %v, got %v", test.expectedToken, apiToken)
		}
	}

	apiToken = realAPIToken
	os.Args = realOSArgs
	flag.CommandLine = realCommandLine
	os.Setenv("RAINFOREST_API_TOKEN", realEnvToken)
}
