package main

import "syscall"

const SM_CXSCREEN = 0
const SM_CYSCREEN = 1

func getResolution() (int, int) {
	var mod = syscall.NewLazyDLL("user32.dll")
	var proc = mod.NewProc("GetSystemMetrics")
	xr, _, _ := proc.Call(uintptr(SM_CXSCREEN))
	yr, _, _ := proc.Call(uintptr(SM_CYSCREEN))
	return int(xr), int(yr)
}
