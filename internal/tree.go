package internal

import "unsafe"

type Art struct{
	count int
	root unsafe.Pointer
}

func NewArt() *Art{
	return &Art{
		count: 0,
		root:  unsafe.Pointer(),
	}

}



