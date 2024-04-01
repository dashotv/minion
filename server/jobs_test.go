package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_CleanupTime(t *testing.T) {
	finished := time.Now().Add(-time.Hour * time.Duration(2))
	fmt.Printf("Time: %v\n", finished)
}
