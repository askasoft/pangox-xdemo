package server

import (
	"testing"

	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/test/require"
	"github.com/askasoft/pangox/xwa/xtxts"
)

func TestLoadMessages(t *testing.T) {
	tb := tbs.NewTextBundles()
	for _, fs := range xtxts.FSs {
		require.NoError(t, tb.LoadFS(fs, "."))
	}
}
