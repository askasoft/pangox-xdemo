package super

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	"github.com/askasoft/pango/str"
)

func shellExecCommand(ctx context.Context, command string) (code int, output string) {
	stdout := &strings.Builder{}

	cmd := exec.CommandContext(ctx, "cmd.exe")
	cmd.Stdin = strings.NewReader(command)
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
