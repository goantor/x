package x

import (
	"time"
)

type UseTimer struct {
	Start time.Time
	Used  time.Duration
}

func NewUseTimer() *UseTimer {
	return &UseTimer{Start: time.Now()}
}

func (t *UseTimer) Record() {
	t.Start = time.Now()
}

func (t *UseTimer) TakeUsed() time.Duration {
	used := time.Now().Sub(t.Start)
	t.Used += used
	return t.Used
}
