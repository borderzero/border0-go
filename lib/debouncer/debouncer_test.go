package debouncer

import (
	"sync"
	"testing"
	"time"
)

func TestDebouncerCloseRace(t *testing.T) {

	d := New(func(int) {}, WithDebounceTime(1*time.Millisecond))
	var wg sync.WaitGroup

	for i := range 1000 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.Do(i)
		}(i)
	}

	time.Sleep(2 * time.Millisecond)
	d.Close()
	wg.Wait()
}
