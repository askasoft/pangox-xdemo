package super

import (
	"bytes"
	"context"
	"strings"

	"github.com/askasoft/pango/osu"
	"github.com/askasoft/pango/str"
)

func shellExecCommand(ctx context.Context, command string) (code int, output string, err error) {
	stdin := strings.NewReader(command + "\r\n")
	stdout := &bytes.Buffer{}

	code, err = osu.ExecCommand(ctx, stdin, stdout, stdout, "cmd")

	output = str.Strip(stdout.String())
	return
}
