package files

import (
	"mime/multipart"
	"net/http"

	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xfs"
)

func Upload(c *xin.Context) {
	UploadFile(c, models.TagTmpFile)
}

func Uploads(c *xin.Context) {
	UploadFiles(c, models.TagTmpFile)
}

func SaveUploadedFile(c *xin.Context, mfh *multipart.FileHeader, tag string) (*xfs.File, error) {
	fid := app.MakeFileID(tag, mfh.Filename)

	tt := tenant.FromCtx(c)
	tfs := tt.FS()
	return xfs.SaveUploadedFile(tfs, fid, mfh, tag)
}

func UploadFile(c *xin.Context, tag string) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fi, err := SaveUploadedFile(c, file, tag)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fr := &xfs.FileResult{File: fi}
	c.JSON(http.StatusOK, fr)
}

func UploadFiles(c *xin.Context, tag string) {
	files, err := c.FormFiles("files")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result := &xfs.FilesResult{}
	for _, file := range files {
		fi, err := SaveUploadedFile(c, file, tag)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		result.Files = append(result.Files, fi)
	}

	c.JSON(http.StatusOK, result)
}
