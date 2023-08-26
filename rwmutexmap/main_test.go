package rwmutexmap_test

import (
	"cowmap/race"
	"cowmap/rwmutexmap"
	"fmt"
	"testing"
	"time"
)

func Test_CopyOnWrite(t *testing.T) {
	if !race.Enabled {
		t.Fatalf("race detector is not enabled")
	}
	t.Log("Running copy on write test")
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
					fmt.Printf("%d: %d\n", goroutineNumber, rwmutexmap.Get(goroutineNumber).(int))
				}
			}
		}(i)
	}

	time.Sleep(time.Second * 10)
	fmt.Println("finalizing test")
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
