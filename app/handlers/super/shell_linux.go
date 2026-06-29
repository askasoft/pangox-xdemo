package super

import (
	"bytes"
	"context"

	"github.com/askasoft/pango/osu"
	"github.com/askasoft/pango/str"
)

func shellExecCommand(ctx context.Context, command string) (code int, output string, err error) {
	stdout := &bytes.Buffer{}

	code, err = osu.ExecCommand(ctx, nil, stdout, stdout, "sh", "-e", "-x", "-c", str.RemoveByte(command, '\r'))

	output = str.Strip(stdout.String())
	return
}
