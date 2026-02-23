package tasks

import (
	"time"

	"github.com/askasoft/pango/sch"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pangox-xdemo/app/jobs"
	"github.com/askasoft/pangox-xdemo/app/jobs/pets"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func JobSchedule() {
	_ = tenant.Iterate(func(tt *tenant.Tenant) error {
		return startScheduleJobChain(tt, "schedule_pets_reset", jobs.JobChainPetReset, pets.PetResetJobChainStart)
	})
}

func startScheduleJobChain(tt *tenant.Tenant, key, jcname string, fn func(tt *tenant.Tenant) error) error {
	expr := tt.SV(key)
	if expr == "" {
		return nil
	}

	periodic, err := sch.ParsePeriodic(expr)
	if err != nil {
		tt.Logger("JOB").Errorf("Invalid setting %q: %v", key, err)
		return nil
	}

	cexpr := periodic.Cron()
	cron, err := sch.ParseCron(cexpr)
	if err != nil {
		tt.Logger("JOB").Errorf("Invalid cron expression %q: %v", cexpr, err)
		return nil
	}

	tjc := tt.JC()
	jc, err := tjc.FindJobChain(jcname, false)
	if err != nil {
		return err
	}

	now := time.Now()
	stm := tmu.TruncateMinutes(now).Add(-time.Minute)

	if jc != nil && jc.CreatedAt.After(stm) {
		stm = jc.CreatedAt
	}

	jtm := cron.Next(stm)
	if jtm.Before(now) {
		return fn(tt)
	}
	return nil
}
