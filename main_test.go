package cowmap_test

import (
	"cowmap"
	"cowmap/race"
	"cowmap/rwmutexmap"
	"testing"
	"time"
)

func Test_CopyOnWrite(t *testing.T) {
	if !race.Enabled {
		t.Fatalf("race detector is not enabled")
	}
	cowmap := cowmap.New()

	quitSignal := make(chan struct{})

	iterSetCowmap(quitSignal, cowmap)

	go func() {
		for {
			iterSetCowmap(quitSignal, cowmap)
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
					_ = cowmap.Get(goroutineNumber).(int)
				}
			}
		}(i)
	}

	time.Sleep(time.Second * 10)
	close(quitSignal)
}

func Benchmark_Maps(b *testing.B) {
	b.Run("concurrent read", func(b *testing.B) {
		b.Run("copy on write map", func(b *testing.B) {
			cowmap := cowmap.New()
			const maxval = 10
			for i := 0; i < maxval; i++ {
				cowmap.Set(i, i)
			}

			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					_ = cowmap.Get(0).(int)
				}
			})
		})

		b.Run("rwmutex map", func(b *testing.B) {
			rwmutexmap := rwmutexmap.New()
			const maxval = 10
			for i := 0; i < maxval; i++ {
				rwmutexmap.Set(i, i)
			}

			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					_ = rwmutexmap.Get(0).(int)
				}
			})

		})
	})

	b.Run("single write", func(b *testing.B) {
		b.Run("copy on write map", func(b *testing.B) {
			cowmap := cowmap.New()
			for i := 0; i < b.N; i++ {
				cowmap.Set(i, i)
			}
		})

		b.Run("rwmutex map", func(b *testing.B) {
			rwmutexmap := rwmutexmap.New()
			for i := 0; i < b.N; i++ {
				rwmutexmap.Set(i, i)
			}
		})
	})
}

func iterSetCowmap(quitSignal <-chan struct{}, cowmap *cowmap.Cowmap) {
	for i := 0; i < 100; i++ {
		select {
		case <-quitSignal:
			return
		default:
			time.Sleep(time.Millisecond)
			cowmap.Set(i, i)
		}
	}
}
