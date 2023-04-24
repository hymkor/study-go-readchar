//go:build !windows
// +build !windows

package main

func nop() {}

func enable(handle int) (func(), error) {
	return nop, nil
}
