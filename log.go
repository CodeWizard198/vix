package vix

import (
	"fmt"
	"time"
)

var (
	greenBackground = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	redBackground   = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blueBackground  = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	reset           = string([]byte{27, 91, 48, 109})
)

const (
	seg = "|"
)

type logs struct {
	IP     string
	Method string
	Path   string
	Times  string
	Status int
}

func (l *logs) logging() {
	now := time.Now().Format("2006-01-02 15:04:05")
	if l.Status == 200 {
		fmt.Println(greenBackground, l.Status, reset, now, seg, l.IP, seg, l.Times, seg, l.Path, seg, blueBackground, l.Method, reset)
		return
	}
	fmt.Println(redBackground, l.Status, reset, now, seg, l.IP, seg, l.Times, seg, l.Path, seg, blueBackground, l.Method, reset)
}
