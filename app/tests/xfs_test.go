package tests

import (
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/askasoft/gogormx/xfs/gormfs"
	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/test/require"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pangox-xdemo/app/schema"
	"github.com/askasoft/pangox-xdemo/app/server/gormdb"
	"github.com/askasoft/pangox/xfs"
	"github.com/askasoft/pangox/xwa"
	"github.com/askasoft/pangox/xwa/xsqls"
)

var sm schema.Schema

func testInit(t *testing.T) {
	if !bol.Atob(os.Getenv("XTEST"), true) {
		t.Skip()
	}

	if sm != "" {
		return
	}

	xwa.SetDirConfig("../../conf/")

	require.NoError(t, xwa.InitConfigs())

	require.NoError(t, xsqls.OpenDatabase())

	if ok, err := schema.ExistsSchema("xtest"); err != nil {
		t.Fatal(err)
	} else {
		if ok {
			require.NoError(t, schema.DeleteSchema("xtest"))
		}
		require.NoError(t, schema.CreateSchema("xtest", "test schema"))
		require.NoError(t, schema.Schema("xtest").InitSchema())
	}

	sm = "xtest"
}

func TestXFS_Sqlx(t *testing.T) {
	testInit(t)

	xfs := sm.FS()

	testXFS(t, xfs)
}

func TestXFS_Gorm(t *testing.T) {
	testInit(t)

	gdb, err := gormdb.OpenDatabase(string(sm))
	require.NoError(t, err)

	xfs := gormfs.FS(gdb, sm.TableFiles())

	testXFS(t, xfs)
}

func testXFS(t *testing.T, xfs xfs.XFS) {
	require.NoError(t, xfs.Truncate())

	tmz, _ := time.Parse(time.RFC3339, "2020-01-02T03:04:05Z")
	tml := tmz.Local()

	tags := []string{"a", "b", "c", "d"}
	for i, tag := range tags {
		n := num.Itoa(i + 1)

		fid := "/" + tag + "/" + n
		name := n + ".test"
		time := tml.Add(tmu.Day)
		data := []byte(tag + n)
		ftag := tag

		// create file
		fw, err := xfs.SaveFile(fid, name, time, data, ftag)
		require.NoError(t, err)

		ff, err := xfs.FindFile(fid)
		require.NoError(t, err)

		fw.Data = nil
		require.Equal(t, fw, ff)

		// read file
		fd, err := xfs.ReadFile(fid)
		require.NoError(t, err)
		require.Equal(t, data, fd)

		// copy file with new tag
		cTag := tag + tag
		cFID := "/" + cTag + "/" + n
		require.NoError(t, xfs.CopyFile(fid, cFID, cTag))

		fcw := *fw
		fcw.ID = cFID
		fcw.Tag = cTag

		ffc, err := xfs.FindFile(cFID)
		require.NoError(t, err)
		require.Equal(t, &fcw, ffc)

		// copy file without tag
		cFID2 := cFID + "2"
		require.NoError(t, xfs.CopyFile(fid, cFID2))

		fcw2 := *fw
		fcw2.ID = cFID2

		ffc2, err := xfs.FindFile(cFID2)
		require.NoError(t, err)
		require.Equal(t, &fcw2, ffc2)

		// move file without tag
		mFID := cFID2 + ".move"
		require.NoError(t, xfs.MoveFile(cFID2, mFID))

		fmw := fcw2
		fmw.ID = mFID

		ffm, err := xfs.FindFile(mFID)
		require.NoError(t, err)
		require.Equal(t, &fmw, ffm)

		// move file with new tag
		mFID2 := mFID + ".again"
		require.NoError(t, xfs.MoveFile(mFID, mFID2, "m"))

		fmw2 := fmw
		fmw2.ID = mFID2
		fmw2.Tag = "m"

		ffm2, err := xfs.FindFile(mFID2)
		require.NoError(t, err)
		require.Equal(t, &fmw2, ffm2)

		// delete file
		require.NoError(t, xfs.DeleteFile(mFID2))
	}

	cnt, err := xfs.DeleteTagged("a")
	require.NoError(t, err)
	require.Equal(t, int64(1), cnt)

	cnt, err = xfs.DeletePrefix("/b/")
	require.NoError(t, err)
	require.Equal(t, int64(1), cnt)

	cnt, err = xfs.DeleteTaggedBefore("c", time.Now())
	require.NoError(t, err)
	require.Equal(t, int64(1), cnt)

	cnt, err = xfs.DeletePrefixBefore("/d/", time.Now())
	require.NoError(t, err)
	require.Equal(t, int64(1), cnt)

	cnt, err = xfs.DeleteBefore(time.Now())
	require.NoError(t, err)
	require.Equal(t, int64(len(tags)), cnt)

	cnt, err = xfs.DeleteAll()
	require.NoError(t, err)
	require.Equal(t, int64(0), cnt)
}
