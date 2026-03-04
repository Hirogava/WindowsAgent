package viewmodel

import "fyne.io/fyne/v2"

type MenuViewModel struct {
	MainPageBtn   func()
	ConfigPageBtn func()

	MainWindow fyne.Window
}

func NewMenuViewModel(window fyne.Window) *MenuViewModel {
	return &MenuViewModel{
		MainPageBtn:   func() {},
		ConfigPageBtn: func() {},
		MainWindow:    window,
	}
}

func (vm *MenuViewModel) SetMainPageBtnAction(action func()) {
	vm.MainPageBtn = action
}

func (vm *MenuViewModel) SetConfigPageBtnAction(action func()) {
	vm.ConfigPageBtn = action
}
