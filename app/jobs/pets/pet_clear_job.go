package pets

import (
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/jobs"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xjm"
)

func init() {
	jobs.RegisterJobRun(jobs.JobNamePetClear, NewPetClearJob)
	jobs.RegisterJobArg(jobs.JobNamePetClear, NewPetClearArg)
}

type PetClearArg struct {
	jobs.ChainArg

	ResetSequence bool `json:"reset_sequence" form:"reset_sequence"`
}

func NewPetClearArg(tt *tenant.Tenant) jobs.IArg {
	pca := &PetClearArg{}
	pca.ResetSequence = true
	return pca
}

func (pca *PetClearArg) Bind(c *xin.Context) error {
	return c.Bind(pca)
}

type PetClearJob struct {
	*jobs.JobRunner[PetClearArg]

	jobs.JobState
}

func NewPetClearJob(tt *tenant.Tenant, job *xjm.Job) jobs.IJobRunner {
	pc := &PetClearJob{}

	pc.JobRunner = jobs.NewJobRunner[PetClearArg](tt, job)

	pc.ChainArg = pc.Arg.ChainArg
	return pc
}

func (pc *PetClearJob) Run() error {
	tt := pc.Tenant
	db := app.SDB()

	pc.Logger.Infof("Delete Pet Files ...")

	sfs := tt.SFS(db)
	cnt, err := sfs.DeleteTagged(models.TagPetFile)
	if err != nil {
		return err
	}
	pc.Logger.Infof("%d Pet Files Deleted.", cnt)

	pc.Logger.Info("Delete Pets ...")
	cnt, err = tt.DeletePets(db)
	if err != nil {
		return err
	}
	pc.Logger.Infof("%d Pets Deleted.", cnt)

	pc.Success = int(cnt)
	if err = pc.SetState(&pc.JobState); err != nil {
		return err
	}

	if pc.Arg.ResetSequence {
		pc.Logger.Info("Pets Sequence Resetted.")
		err = tt.ResetPetsAutoIncrement(db)
		if err != nil {
			return err
		}
	}

	return nil
}
