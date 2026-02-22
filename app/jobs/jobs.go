package jobs

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xjm"
	"github.com/askasoft/pangox/xwa/xjobs"
)

const (
	JobNameUserCsvImport = "UserCsvImport"
	JobNamePetClear      = "PetClear"
	JobNamePetCatGen     = "PetCatGen"
	JobNamePetDogGen     = "PetDogGen"
)

var (
	JobStartAuditLogs = map[string]string{
		JobNameUserCsvImport: models.AL_USERS_IMPORT_START,
		JobNamePetClear:      models.AL_PETS_CLEAR_START,
		JobNamePetCatGen:     models.AL_PETS_CAT_CREATE_START,
		JobNamePetDogGen:     models.AL_PETS_DOG_CREATE_START,
	}

	JobCancelAuditLogs = map[string]string{
		JobNameUserCsvImport: models.AL_USERS_IMPORT_CANCEL,
		JobNamePetClear:      models.AL_PETS_CLEAR_CANCEL,
		JobNamePetCatGen:     models.AL_PETS_CAT_CREATE_CANCEL,
		JobNamePetDogGen:     models.AL_PETS_DOG_CREATE_CANCEL,
	}
)

var (
	ErrJobOverflow = errors.New("job overflow")
)

var (
	ttJobRuns = xjobs.NewJobsMap()
	ttJobLock sync.Mutex
)

// Starts iterate tenants to start jobs
func Starts() {
	mar := ini.GetInt("job", "maxTotalRunnings", 10)

	if mar-ttJobRuns.Total() > 0 {
		err := tenant.Iterate(StartJobs)
		if err != nil && !errors.Is(err, ErrJobOverflow) {
			log.Errorf("jobs.Starts(): %v", err)
		}

		// sleep 1s to let all job go-routine start
		time.AfterFunc(time.Second, ttJobRuns.Clean)
	}

	log.Info(Stats())
}

// StartJobs start tenant jobs
func StartJobs(tt *tenant.Tenant) error {
	ttJobLock.Lock()
	defer ttJobLock.Unlock()

	mar := ini.GetInt("job", "maxTotalRunnings", 10)
	mtr := ini.GetInt("job", "maxTenantRunnings", 10)

	a := mar - ttJobRuns.Total()
	if a <= 0 {
		return ErrJobOverflow
	}

	c := mtr - ttJobRuns.Count(string(tt.Schema))
	if c <= 0 {
		return nil
	}

	if c > a {
		c = a
	}

	return tt.JM().StartJobs(c, func(job *xjm.Job) {
		go runJob(tt, job)
	})
}

func runJob(tt *tenant.Tenant, job *xjm.Job) {
	logger := tt.Logger("JOB")

	jrc, ok := jobRunCreators[job.Name]
	if !ok {
		logger.Errorf("No Job Runner Creator %q", job.Name)
		return
	}

	run := jrc(tt, job)

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Job %s#%d panic: %v", job.Name, job.ID, err)
		}

		log.Info(Stats())
	}()

	logger.Debugf("Start job %s#%d", job.Name, job.ID)

	ttJobRuns.AddJob(string(tt.Schema), job)

	defer ttJobRuns.DelJob(string(tt.Schema), job)

	xjobs.RunJob(run)
}

func Stats() string {
	total, stats := ttJobRuns.Stats()
	return fmt.Sprintf("INSTANCE ID: 0x%04x, JOB RUNNING: %d\n%s", app.InstanceID(), total, stats)
}

// ------------------------------------
func ReappendJobs() {
	before := time.Now().Add(-1 * ini.GetDuration("job", "reappendBefore", time.Minute*30))

	_ = tenant.Iterate(func(tt *tenant.Tenant) error {
		tjm := tt.JM()
		cnt, err := tjm.ReappendJobs(before)
		if err != nil {
			tt.Logger("JOB").Errorf("Failed to ReappendJobs(%q, %q): %v", string(tt.Schema), before.Format(time.RFC3339), err)
		} else if cnt > 0 {
			tt.Logger("JOB").Infof("ReappendJobs(%q, %q): %d", string(tt.Schema), before.Format(time.RFC3339), cnt)
		}
		return err
	})
}

// ------------------------------------
// CleanOutdatedJobs iterate schemas to clean outdated jobs
func CleanOutdatedJobs() {
	before := time.Now().Add(-1 * ini.GetDuration("job", "outdatedBefore", time.Hour*24*10))

	_ = tenant.Iterate(func(tt *tenant.Tenant) error {
		return app.SDB().Transaction(func(tx *sqlx.Tx) error {
			logger := tt.Logger("JOB")

			sfs := tt.SFS(tx)
			cnt, err := sfs.DeleteTaggedBefore(models.TagJobFile, before)
			if err != nil {
				return err
			}
			if cnt > 0 {
				logger.Infof("CleanOutdatedJobFiles(%q): %d", before.Format(time.RFC3339), cnt)
			}

			sjm := tt.SJM(tx)
			cnt, _, err = sjm.CleanOutdatedJobs(before)
			if err != nil {
				return err
			}
			if cnt > 0 {
				logger.Infof("CleanOutdatedJobs(%q): %d", before.Format(time.RFC3339), cnt)
			}

			xjc := tt.JC()
			cnt, err = xjc.CleanOutdatedJobChains(before)
			if err != nil {
				return err
			}
			if cnt > 0 {
				logger.Infof("CleanOutdatedJobChains(%q): %d", before.Format(time.RFC3339), cnt)
			}

			return nil
		})
	})
}
