package super

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"syscall"

	"github.com/askasoft/pango/str"
)

func shellExecCommand(ctx context.Context, command string) (code int, output string) {
	stdout := &strings.Builder{}

	command = str.RemoveByte(command, '\r')
	exe := "sh"
	arg := []string{"-e", "-x", "-c", command}

	cmd := exec.CommandContext(ctx, exe, arg...)

	// Set Process Group ID so all child processes share the same group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Tell Go to kill the entire process group immediately on timeout.
	// This prevents Go from hanging on open stdout/stderr pipes.
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			// Negative PID targets the entire process group
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}

	cmd.Stdout = stdout
	cmd.Stderr = stdout

	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			code = ee.ExitCode()
		}
		output = err.Error()
		return
	}

	output = str.Strip(stdout.String())
	return
}
