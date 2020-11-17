package trackip

import "context"

type Usecase interface {
	TrackContainers(ctx context.Context, loopIntervalSecond int64) error
}
