package server

import (
	"testing"

	"github.com/askasoft/pango/test/require"
	"github.com/askasoft/pangox/xwa/xtpls"
)

func TestLoadTemplates(t *testing.T) {
	ht := xtpls.NewHTMLTemplates()
	for _, fs := range xtpls.FSs {
		require.NoError(t, ht.LoadFS(fs, "."))
	}
}
