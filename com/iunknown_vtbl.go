package com

import "structs"

type IUnknownVtbl struct {
	_ structs.HostLayout
	// every COM object starts with these three
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// _QueryInterface2 uintptr
}
