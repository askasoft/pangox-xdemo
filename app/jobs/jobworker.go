package jobs

import (
	"context"

	"github.com/askasoft/pangox/xwa/xjobs"
)

type JobWorker[R any] = xjobs.JobWorker[R]

func StreamRun[T any](ctx context.Context, sr xjobs.IStreamRun[T]) error {
	return xjobs.StreamRun(ctx, sr)
}

func SubmitRun[T any, R any](ctx context.Context, sr xjobs.ISubmitRun[T, R]) error {
	return xjobs.SubmitRun(ctx, sr)
}

func StreamOrSubmitRun[T any, R any](ssr xjobs.IStreamSubmitRun[T, R]) error {
	return xjobs.StreamOrSubmitRun(ssr)
}
