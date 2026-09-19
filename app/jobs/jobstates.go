package jobs

import (
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox/xwa/xjobs"
)

func maxFindTargetsCount(limit int) int {
	if limit > 0 {
		return min(limit, MaxFindTargetsCount)
	}
	return MaxFindTargetsCount
}

type JobState = xjobs.JobState
type JobStateLx = xjobs.JobStateLx
type JobStateSx = xjobs.JobStateSx

type JobStateLix struct {
	xjobs.JobStateLix
}

func (jsl *JobStateLix) MaxFindTargetsCount() int {
	return maxFindTargetsCount(jsl.Limit)
}

func (jsl *JobStateLix) AddLastIDFilter(sqb *sqlx.Builder, col string) {
	sqb.Gt(col, jsl.GetLastID())
}

type JobStateSix struct {
	xjobs.JobStateSix
}

func (jss *JobStateSix) MaxFindTargetsCount() int {
	return maxFindTargetsCount(jss.Limit)
}

func (jss *JobStateSix) AddLastIDFilter(sqb *sqlx.Builder, col string) {
	sqb.Gt(col, jss.GetLastID())
}

type JobStateLixs struct {
	xjobs.JobStateLixs
}

func (jse *JobStateLixs) MaxFindTargetsCount() int {
	return maxFindTargetsCount(jse.Limit)
}

func (jse *JobStateLixs) AddLastIDFilter(sqb *sqlx.Builder, col string) {
	sqb.Gt(col, jse.GetLastID())
}

type IJobStater = xjobs.IJobStater

func InitState(js IJobStater, limit int) error {
	return xjobs.InitState(js, limit)
}
