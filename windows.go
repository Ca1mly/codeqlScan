//go:build windows

package main

import (
	"github.com/lxn/win"
)

func init() {
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
}