package jobs

import (
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox/xwa/xjobs"
)

type JobState = xjobs.JobState
type JobStateLx = xjobs.JobStateLx
type JobStateSx = xjobs.JobStateSx

type JobStateLix struct {
	xjobs.JobStateLix
}

func (jsl *JobStateLix) AddLastIDFilter(sqb *sqlx.Builder, col string) {
	sqb.Gt(col, jsl.LastID)
}

type JobStateSix struct {
	xjobs.JobStateSix
}

func (jss *JobStateSix) AddLastIDFilter(sqb *sqlx.Builder, col string) {
	sqb.Gt(col, jss.LastID)
}

type JobStateLixs struct {
	xjobs.JobStateLixs
}

func (jse *JobStateLixs) AddLastIDFilter(sqb *sqlx.Builder, col string) {
	if len(jse.LastIDs) > 0 {
		sqb.Gt(col, asg.Max(jse.LastIDs))
	}
}
