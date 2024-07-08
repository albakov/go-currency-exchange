package util

import (
	"log"
)

// LogError wrapper for log.Printf
func LogError(fileName, operation string, err error) {
	log.Printf("%v -> %v error: %v", fileName, operation, err)
}
