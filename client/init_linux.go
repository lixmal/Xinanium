package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

func init() {
	// X11 multithreading, linux/X11 only
	C.XInitThreads()
}
