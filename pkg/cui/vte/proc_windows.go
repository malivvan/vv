//go:build windows

package vte

import "syscall"

var sysProcAttr = &syscall.SysProcAttr{}
