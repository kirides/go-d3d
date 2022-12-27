package framelimiter

import "time"

// finer granularity for sleeping
type frameLimiter struct {
	DesiredFps  int
	frameTimeNs int64

	LastFrameTime     time.Time
	LastSleepDuration time.Duration

	DidSleep bool
	DidSpin  bool
}

func New(desiredFps int) *frameLimiter {
	return &frameLimiter{
		DesiredFps:    desiredFps,
		frameTimeNs:   (time.Second / time.Duration(desiredFps)).Nanoseconds(),
		LastFrameTime: time.Now(),
	}
}

func (l *frameLimiter) Wait() {
	l.DidSleep = false
	l.DidSpin = false

	now := time.Now()
	spinWaitUntil := now

	sleepTime := l.frameTimeNs - now.Sub(l.LastFrameTime).Nanoseconds()

	if sleepTime > int64(1*time.Millisecond) {
		if sleepTime < int64(30*time.Millisecond) {
			l.LastSleepDuration = time.Duration(sleepTime / 8)
		} else {
			l.LastSleepDuration = time.Duration(sleepTime / 4 * 3)
		}
		time.Sleep(time.Duration(l.LastSleepDuration))
		l.DidSleep = true

		newNow := time.Now()
		spinWaitUntil = newNow.Add(time.Duration(sleepTime) - newNow.Sub(now))
		now = newNow

		for spinWaitUntil.After(now) {
			now = time.Now()
			// SPIN WAIT
			l.DidSpin = true
		}
	} else {
		l.LastSleepDuration = 0
		spinWaitUntil = now.Add(time.Duration(sleepTime))
		for spinWaitUntil.After(now) {
			now = time.Now()
			// SPIN WAIT
			l.DidSpin = true
		}
	}
	l.LastFrameTime = time.Now()
}
