package com

type IUnknownVtbl struct {
	// every COM object starts with these three
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// _QueryInterface2 uintptr
}
