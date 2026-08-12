package generate

import (
	"os"
	"os/exec"
	"syscall"
)

// editorName returns the user's preferred editor, taken from the EDITOR
// environment variable and falling back to vim if it is unset.
func editorName() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return editor
}

// EditorCommand returns a command that opens path in the user's preferred
// editor. The editor is taken from the EDITOR environment variable, falling
// back to vim if it is unset.
func EditorCommand(path string) *exec.Cmd {
	return exec.Command(editorName(), path)
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

// ExecEditor replaces the current process image with the user's preferred
// editor opened on path, using execve so the editor takes over the process
// outright and control never returns to the caller. It only returns if the
// replacement fails to happen at all (editor not found, or process
// replacement unsupported on this platform, e.g. Windows) — the caller
// should fall back to OpenEditor in that case.
func ExecEditor(path string) error {
	editor := editorName()

	binary, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	return syscall.Exec(binary, []string{editor, path}, os.Environ())
}
