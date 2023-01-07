package com

import (
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ReflectQueryInterface(self interface{}, method uintptr, interfaceID *windows.GUID, obj interface{}) int32 {
	selfValue := reflect.ValueOf(self).Elem()
	objValue := reflect.ValueOf(obj).Elem()

	hr, _, _ := syscall.SyscallN(
		method,
		selfValue.UnsafeAddr(),
		uintptr(unsafe.Pointer(interfaceID)),
		objValue.Addr().Pointer())

	return int32(hr)
}
