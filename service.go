// +build windows

package main

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func exePath() (string, error) {
	if l, err := syscall.LoadDLL("kernel32.dll"); err == nil {
		defer l.Release()
		if p, err := l.FindProc("GetModuleFileNameW"); err == nil {
			b := make([]uint16, syscall.MAX_PATH)
			if ret, _, err := p.Call(
				0,
				uintptr(unsafe.Pointer(&b[0])),
				uintptr(len(b)),
			); ret != 0 {
				return string(utf16.Decode(b[:ret])), nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}
