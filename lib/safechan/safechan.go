package safechan

import "fmt"

// Write safely writes to a channel, recovering from any panics.
func Write[T any](channel chan<- T, value T) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
		}
	}()
	channel <- value
	return
}
