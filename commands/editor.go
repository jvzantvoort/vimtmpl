// Package commands provides editor integration and command execution utilities.
package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Editor represents an editor command executor for git and other tools.
type Editor struct {
	Path       *Path
	Cwd        string
	Command    string
	CommandMap map[string]string
}

// Prefix returns a formatted string with the caller function name for logging.
func (g Editor) Prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("Editor.%s", elements[len(elements)-1])
}

// Execute runs the editor command with the provided arguments.
// Returns slices of stdout and stderr output, and an error if execution fails.
func (g Editor) Execute(args ...string) ([]string, []string, error) {
	log_prefix := g.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	stdout_list := []string{}
	stderr_list := []string{}
	cmnd := []string{}

	cmnd = append(cmnd, args...)

	log.Debugf("%s: command %s %s", log_prefix, g.Command, strings.Join(cmnd, " "))

	cmd := exec.Command(g.Command, cmnd...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("stdout pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Errorf("stderr pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	cmd.Dir = g.Cwd
	err = cmd.Start()
	if err != nil {
		log.Errorf("Start failed, %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		msg := scanner.Text()
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		stderr_list = append(stderr_list, msg)
	}

	eerror := cmd.Wait()
	if eerror != nil {
		log.Errorf("command failed, %v", eerror)
	}
	return stdout_list, stderr_list, eerror
}

// NewEditor create a new git object
func NewEditor() *Editor {
	retv := &Editor{}

	log_prefix := retv.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	retv.Path = NewPath("PATH")

	retv.CommandMap = map[string]string{
		"windows": "vim.exe",
		"linux":   "vim",
		"default": "vim",
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s: %s", log_prefix, err)
	} else {
		retv.Cwd = dir
	}

	if result, err := retv.Path.LookupPlatform(retv.CommandMap); err == nil {
		retv.Command = result
	}

	return retv
}

func (g Editor) Edit(args ...string) ([]string, []string, error) {
	log_prefix := g.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	arglist := []string{}
	arglist = append(arglist, args...)

	return g.Execute(arglist...)
}
