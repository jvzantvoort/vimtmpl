package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// runSelfAsSubprocess re-executes the current test binary, running only
// TestMainHelperProcess, which in turn invokes main() with os.Args rebuilt
// from VIMTMPL_MAIN_ARGS. This lets us exercise main()'s os.Exit paths
// without terminating the real test process.
func runSelfAsSubprocess(t *testing.T, home string, args []string) (string, error) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestMainHelperProcess$")
	cmd.Env = append(os.Environ(),
		"VIMTMPL_BE_MAIN=1",
		"VIMTMPL_MAIN_ARGS="+strings.Join(args, "\x1f"),
		"HOME="+home,
	)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func asExitError(err error) (int, bool) {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return 0, false
	}
	return exitErr.ExitCode(), true
}
