package util

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Retrieve the full path to the current executable using the Windows API.
func Executable() (string, error) {
	l, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return "", err
	}
	defer l.Release()
	p, err := l.FindProc("GetModuleFileNameW")
	if err != nil {
		return "", err
	}
	b := make([]uint16, syscall.MAX_PATH)
	ret, _, err := p.Call(
		0,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
	)
	if ret != 0 {
		return string(utf16.Decode(b[:ret])), nil
	} else {
		return "", err
	}
}
