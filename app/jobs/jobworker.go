package jobs

import (
	"github.com/askasoft/pangox/xwa/xjobs"
)

type JobWorker[R any] = xjobs.JobWorker[R]

func StreamRun[T any](sr xjobs.IStreamRun[T]) error {
	return xjobs.StreamRun(sr)
}

func SubmitRun[T any, R any](sr xjobs.ISubmitRun[T, R]) error {
	return xjobs.SubmitRun(sr)
}

func StreamOrSubmitRun[T any, R any](ssr xjobs.IStreamSubmitRun[T, R]) error {
	return xjobs.StreamOrSubmitRun(ssr)
}
