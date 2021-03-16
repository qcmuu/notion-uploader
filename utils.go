package main

import (
	"fmt"
	"time"
)

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

// buildOp creates an Operation for this block
func buildOp(blockID, command string, path []string, args interface{}) *Operation {
	return &Operation{
		Point: Pointer{
			ID:    blockID,
			Table: "block",
		},
		Path:    path,
		Command: command,
		Args:    args,
	}
}

// DotTicker starts a infinity dot animation and returns a chan to stop.
func DotTicker() *chan struct{} {
	tick := time.NewTicker(time.Second)
	end := make(chan struct{})
	go func() {
		for {
			select {
			case <-tick.C:
				fmt.Printf(".")
			case <-end:
				return
			}
		}
	}()
	return &end
}
