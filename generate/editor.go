package generate

import (
	"os"
	"os/exec"
)

// EditorCommand returns a command that opens path in the user's preferred
// editor. The editor is taken from the EDITOR environment variable, falling
// back to vim if it is unset.
func EditorCommand(path string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return exec.Command(editor, path)
}

// OpenEditor opens path in the user's preferred editor, connecting it to the
// current process's standard streams, and waits for it to exit.
func OpenEditor(path string) error {
	cmd := EditorCommand(path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
