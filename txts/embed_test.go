package txts

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/askasoft/pango/ini"
)

func TestEmbedFS(t *testing.T) {
	fmt.Println("------------------------")
	err := fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path, d.IsDir())
		if d.IsDir() {
			return nil
		}
		return ini.NewIni().LoadFileFS(FS, path)
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("------------------------")
}
