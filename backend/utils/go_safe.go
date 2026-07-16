package utils

import "log"

func GoSafe(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] Protected goroutine failed: %v", r)
			}
		}()
		fn()
	}()
}
