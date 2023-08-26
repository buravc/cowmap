package rwmutexmap_test

import (
	"cowmap/race"
	"cowmap/rwmutexmap"
	"testing"
	"time"
)

func Test_CopyOnWrite(t *testing.T) {
	if !race.Enabled {
		t.Fatalf("race detector is not enabled")
	}
	rwmutexmap := rwmutexmap.New()

	quitSignal := make(chan struct{})

	iterSetRwmutexmap(quitSignal, rwmutexmap)

	go func() {
		for {
			iterSetRwmutexmap(quitSignal, rwmutexmap)
		}
	}()

	for i := 0; i < 10; i++ {
		go func(goroutineNumber int) {
			for {
				select {
				case <-quitSignal:
					return
				default:
					time.Sleep(time.Millisecond * 100)
					_ = rwmutexmap.Get(goroutineNumber).(int)
				}
			}
		}(i)
	}

	time.Sleep(time.Second * 10)
	quitSignal <- struct{}{}
}

func iterSetRwmutexmap(quitSignal <-chan struct{}, rwmutexmap *rwmutexmap.Rwmutexmap) {
	for i := 0; i < 100; i++ {
		select {
		case <-quitSignal:
			return
		default:
			time.Sleep(time.Millisecond)
			rwmutexmap.Set(i, i)
		}
	}
}
