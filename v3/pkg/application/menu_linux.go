//go:build linux

package application

import "unsafe"

type linuxMenu struct {
	menu   *Menu
	native unsafe.Pointer
}

func newMenuImpl(menu *Menu) *linuxMenu {
	result := &linuxMenu{
		menu: menu,
	}
	return result
}

func (l *linuxMenu) update() {
	// FIXME: This is a noop right now
}
