package config

import (
	"sync"
	"time"
)

var (
	once sync.Once
	mu   sync.Mutex
)

func SetTimeZone(locationName string) error {
	var err error
	once.Do(func() {
		var timeZone *time.Location
		timeZone, err = time.LoadLocation(locationName)
		if err != nil {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		time.Local = timeZone
	})

	return err
}
