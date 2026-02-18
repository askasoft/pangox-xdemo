package jobs

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xfs"
	"github.com/askasoft/pangox/xjm"
	"github.com/askasoft/pangox/xwa/xjobs"
)

var (
	ErrItemSkip = errors.New("item skip")
)

type IArg interface {
	Bind(c *xin.Context) error
}

type JobArgCreater func(*tenant.Tenant) IArg

var jobArgCreators = map[string]JobArgCreater{}

func RegisterJobArg(name string, jac JobArgCreater) {
	jobArgCreators[name] = jac
}

func CreateJobArg(tt *tenant.Tenant, name string) (IArg, error) {
	jac, ok := jobArgCreators[name]
	if !ok {
		return nil, fmt.Errorf("missing job argument %q", name)
	}

	return jac(tt), nil
}

type FileArg struct {
	File string `json:"file,omitempty" form:"-"`
}

func (fa *FileArg) GetFile() string {
	return fa.File
}

func (fa *FileArg) SetFile(tt *tenant.Tenant, mfh *multipart.FileHeader) error {
	fid := app.MakeFileID(models.TagJobFile, mfh.Filename)
	tfs := tt.FS()
	if _, err := xfs.SaveUploadedFile(tfs, fid, mfh, models.TagJobFile); err != nil {
		return err
	}

	fa.File = fid
	return nil
}

type CsvFileArg struct {
	FileArg
}

func (cfa *CsvFileArg) BindFile(c *xin.Context) error {
	mfh, err := c.FormFile("file")
	if err != nil {
		return tbs.Error(c.Locale, "csv.error.required")
	}

	tt := tenant.FromCtx(c)
	if err = cfa.SetFile(tt, mfh); err != nil {
		return tbs.Errorf(c.Locale, "csv.error.read", err)
	}
	return nil
}

type normalizer interface {
	Normalize()
}

func ArgBind(c *xin.Context, a any) error {
	err := c.Bind(a)

	if nm, ok := a.(normalizer); ok {
		nm.Normalize()
	}

	return err
}

type IJobRunner = xjobs.IJobRunner

type JobRunCreator func(*tenant.Tenant, *xjm.Job) IJobRunner

var jobRunCreators = map[string]JobRunCreator{}

func RegisterJobRun(name string, jrc JobRunCreator) {
	jobRunCreators[name] = jrc
}

type FailedItem = xjobs.FailedItem

type JobContext = xjobs.JobContext

type JobRunner[A any] struct {
	*xjobs.JobRunner

	Tenant *tenant.Tenant
	Logger log.Logger
	Arg    A
}

func NewJobRunner[A any](tt *tenant.Tenant, job *xjm.Job) *JobRunner[A] {
	job.RID = app.Sequencer().NextID().Int64()

	jr := &JobRunner[A]{
		JobRunner: xjobs.NewJobRunner(job, tt.JC(), tt.JM(), tt.Logger("JOB")),
		Tenant:    tt,
	}

	xjm.MustDecode(job.Param, &jr.Arg)

	jr.Log().SetProp("VERSION", app.Version())
	jr.Log().SetProp("REVISION", app.Revision())
	jr.Log().SetProp("TENANT", string(tt.Schema))
	jr.Logger = jr.Log().GetLogger("JOB")

	jr.JobChainContinue = jr.jobChainContinue

	return jr
}

func (jr *JobRunner[A]) jobChainContinue(next *xjobs.JobRunState) error {
	return JobChainAppendJob(jr.Tenant, next.Name, jr.Locale(), jr.ChainID(), jr.ChainSeq+1, jr.ChainData)
}
