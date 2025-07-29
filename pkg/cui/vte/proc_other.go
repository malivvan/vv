//go:build !windows

package vte

import "syscall"

var sysProcAttr = &syscall.SysProcAttr{
	Setsid:  true,
	Setctty: true,
	Ctty:    1,
}
