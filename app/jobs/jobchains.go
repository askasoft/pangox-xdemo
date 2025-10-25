package jobs

import (
	"fmt"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox/xjm"
	"github.com/askasoft/pangox/xwa/xjobs"
)

const (
	JobChainPetReset = "PetReset"
)

var (
	JobChainStartAuditLogs = map[string]string{
		JobChainPetReset: models.AL_PETS_RESET_START,
	}

	JobChainCancelAuditLogs = map[string]string{
		JobChainPetReset: models.AL_PETS_RESET_CANCEL,
	}
)

type IChainArg = xjobs.IChainArg
type ChainArg = xjobs.ChainArg
type JobRunState = xjobs.JobRunState

func JobChainDecodeStates(state string) []*JobRunState {
	return xjobs.JobChainDecodeStates(state)
}

func JobChainEncodeStates(states []*JobRunState) string {
	return xjobs.JobChainEncodeStates(states)
}

func JobChainInitStates(jns ...string) []*JobRunState {
	return xjobs.JobChainInitStates(jns...)
}

func JobChainAbort(xjc xjm.JobChainer, tjm xjm.JobManager, jc *xjm.JobChain, reason string) error {
	return xjobs.JobChainAbort(xjc, tjm, jc, reason)
}

func JobChainCancel(xjc xjm.JobChainer, tjm xjm.JobManager, jc *xjm.JobChain, reason string) error {
	return xjobs.JobChainCancel(xjc, tjm, jc, reason)
}

func JobFindAndAbortChain(xjc xjm.JobChainer, cid, jid int64, jname, reason string) error {
	return xjobs.JobFindAndAbortChain(xjc, cid, jid, jname, reason)
}

func JobFindAndCancelChain(xjc xjm.JobChainer, cid, jid int64, jname, reason string) error {
	return xjobs.JobFindAndCancelChain(xjc, cid, jid, jname, reason)
}

func JobChainStart(tt *tenant.Tenant, chainName string, states []*JobRunState, jobName, jobLocale string, jobParam IChainArg) (cid int64, err error) {
	state := JobChainEncodeStates(states)

	err = app.SDB().Transaction(func(tx *sqlx.Tx) error {
		sjc := tt.SJC(tx)
		cid, err = sjc.CreateJobChain(chainName, state)
		if err != nil {
			return err
		}

		_, cdt := jobParam.GetChain()
		jobParam.SetChain(0, cdt)
		jParam := xjm.MustEncode(jobParam)

		sjm := tt.SJM(tx)
		_, err = sjm.AppendJob(cid, jobName, jobLocale, jParam)

		return err
	})
	if err == nil {
		go StartJobs(tt) //nolint: errcheck
	}

	return
}

func JobChainInitAndStart(tt *tenant.Tenant, cn string, jns ...string) error {
	states := JobChainInitStates(jns...)

	arg, err := CreateJobArg(tt, jns[0])
	if err != nil {
		tt.Logger("JOB").Error("Failed to create JobArg for %q: %v", jns[0], err)
		return err
	}

	if _, ok := arg.(IChainArg); !ok {
		err = fmt.Errorf("invalid chain job %q argument: %T", jns[0], arg)
		tt.Logger("JOB").Error(err)
		return err
	}

	cid, err := JobChainStart(tt, cn, states, jns[0], asg.First(app.Locales()), arg.(IChainArg))
	if err != nil {
		tt.Logger("JOB").Errorf("Failed to start JobChain %q: %v", cn, err)
		return err
	}

	tt.Logger("JOB").Infof("Start JobChain %q: #%d", cn, cid)
	return nil
}

func JobChainAppendJob(tt *tenant.Tenant, name, locale string, cid int64, csq int, cdt bool) error {
	tjm := tt.JM()

	arg, err := CreateJobArg(tt, name)
	if err != nil {
		return err
	}

	if ica, ok := arg.(IChainArg); ok {
		ica.SetChain(csq, cdt)
	} else {
		return fmt.Errorf("invalid chain job %q", name)
	}

	param := xjm.MustEncode(arg)
	if _, err := tjm.AppendJob(cid, name, locale, param); err != nil {
		return err
	}

	go StartJobs(tt) //nolint: errcheck

	return nil
}
