package main

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func TestMain(t *testing.T) {
	commands := []string{"run", "rerun", "new", "validate", "upload", "rm", "download", "csv-upload", "mobile-upload", "report", "sites", "environments", "folders", "filters", "browsers", "features", "run-groups", "update"}

	for _, command := range commands {
		if os.Getenv("TEST_EXIT") == "1" {
			os.Args = []string{"./rainforest", command, "--not-real-flag"}
			main()
			return
		}

		cmd := exec.Command(os.Args[0], "-test.run=Main")
		cmd.Env = append(os.Environ(), "TEST_EXIT=1")
		err := cmd.Run()

		if err == nil {
			t.Error("Expected exit error was not received")
		} else if e, ok := err.(*exec.ExitError); !ok || e.Success() {
			t.Errorf("Unexpected error. Expected an exit error with non-zero status. Got %#v", e.Error())
		}
	}
}

func TestEmptyToken(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		os.Args = []string{"./rainforest", "run", "--token", ""}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=EmptyToken")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()

	if err == nil {
		t.Error("Expected exit error was not received")
	} else if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Errorf("Unexpected error. Expected an exit error with non-zero status. Got %#v", e.Error())
	}
}

func TestShuffleFlags(t *testing.T) {
	var testCases = []struct {
		testArgs []string
		want     []string
	}{
		{
			testArgs: []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag", "--site", "123"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123", "--debug"},
			want:     []string{"./rainforest", "--token", "foobar", "--debug", "run", "--tags", "tag,bag", "--site", "123"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123", "--run-group-id", "255"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag", "--site", "123", "--run-group-id", "255"},
		},
		{
			testArgs: []string{"./rainforest", "--skip-update", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123"},
			want:     []string{"./rainforest", "--skip-update", "--token", "foobar", "run", "--tags", "tag,bag", "--site", "123"},
		},
		{
			testArgs: []string{"./rainforest", "run", "-f", "foo.rfml", "bar.rfml", "--token", "foobar"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "-f", "foo.rfml", "bar.rfml"},
		},
		{
			testArgs: []string{"./rainforest", "run", "-f", "foo.rfml"},
			want:     []string{"./rainforest", "run", "-f", "foo.rfml"},
		},
		{
			testArgs: []string{"./rainforest", "run", "-f", "foo.rfml", "--disable-telemetry"},
			want:     []string{"./rainforest", "--disable-telemetry", "run", "-f", "foo.rfml"},
		},
	}

	for _, tCase := range testCases {
		got := shuffleFlags(tCase.testArgs)
		if !reflect.DeepEqual(tCase.want, got) {
			t.Errorf("shuffleFlags returned %+v, want %+v", got, tCase.want)
		}
	}
}

func TestUserAgent(t *testing.T) {
	os.Unsetenv("ORB_VERSION")
	os.Unsetenv("GH_ACTION_VERSION")
	os.Args = []string{"./rainforest"}
	main()

	if api == nil {
		t.Error("Expected api to be set")
	}

	userAgent := "rainforest-cli/" + version
	if api.UserAgent != userAgent {
		t.Errorf("main() didn't set proper UserAgent %+v, want %+v", api.UserAgent, userAgent)
	}
}

func TestSendTelemetry(t *testing.T) {
	os.Args = []string{"./rainforest"}
	main()

	if api == nil {
		t.Error("Expected api to be set")
	}

	if api.SendTelemetry != true {
		t.Errorf("main() didn't default SendTelemetry - got %+v, want true", api.SendTelemetry)
	}
}

func TestDisableTelemetry(t *testing.T) {
	os.Args = []string{"./rainforest", "--disable-telemetry"}
	main()

	if api == nil {
		t.Error("Expected api to be set")
	}

	if api.SendTelemetry != false {
		t.Errorf("main() didn't disable SendTelemetry - got %+v, want false", api.SendTelemetry)
	}
}

func TestUserAgentWithOrb(t *testing.T) {
	os.Unsetenv("GH_ACTION_VERSION")
	os.Setenv("ORB_VERSION", "1.3.1")
	os.Args = []string{"./rainforest"}

	main()

	if api == nil {
		t.Error("Expected api to be set")
	}

	userAgent := "rainforest-cli/" + version + " rainforest-orb/1.3.1"
	if api.UserAgent != userAgent {
		t.Errorf("main() with orb didn't set proper UserAgent %+v, want %+v", api.UserAgent, userAgent)
	}
}

func TestUserAgentWithGitHubAction(t *testing.T) {
	os.Unsetenv("ORB_VERSION")
	os.Setenv("GH_ACTION_VERSION", "0.4.2")
	os.Args = []string{"./rainforest"}

	main()

	if api == nil {
		t.Error("Expected api to be set")
	}

	userAgent := "rainforest-cli/" + version + " rainforest-gh-action/0.4.2"
	if api.UserAgent != userAgent {
		t.Errorf("main() with GH action didn't set proper UserAgent %+v, want %+v", api.UserAgent, userAgent)
	}
}

var errStub = errors.New("STUB")

// fakeContext is a helper for testing the cli interfacing functions
type fakeContext struct {
	mappings map[string]interface{}
	args     cli.Args
}

func (f fakeContext) String(s string) string {
	val, ok := f.mappings[s].(string)

	if ok {
		return val
	}
	return ""
}

func (f fakeContext) GlobalString(s string) string {
	return f.String(s)
}

func (f fakeContext) StringSlice(s string) []string {
	val, ok := f.mappings[s].([]string)

	if ok {
		return val
	}
	return []string{}
}

func (f fakeContext) GlobalStringSlice(s string) []string {
	return f.StringSlice(s)
}

func (f fakeContext) Bool(s string) bool {
	val, ok := f.mappings[s].(bool)

	if ok {
		return val
	}
	return false
}

func (f fakeContext) GlobalBool(s string) bool {
	return f.Bool(s)
}

func (f fakeContext) Int(s string) int {
	val, ok := f.mappings[s].(int)

	if ok {
		return val
	}
	return 0
}

func (f fakeContext) GlobalInt(s string) int {
	return f.GlobalInt(s)
}

func (f fakeContext) Uint(s string) uint {
	val, ok := f.mappings[s].(uint)

	if ok {
		return val
	}
	return 0
}

func (f fakeContext) GlobalUint(s string) uint {
	return f.GlobalUint(s)
}

func (f fakeContext) Args() cli.Args {
	return f.args
}

func newFakeContext(mappings map[string]interface{}, args cli.Args) *fakeContext {
	return &fakeContext{mappings, args}
}
