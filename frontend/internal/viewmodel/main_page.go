package viewmodel

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type MainPageViewModel struct {
	StartLabelText *widget.Label
	TestDataLabel *widget.Label
	TestDataEntry binding.String
	SaveBtn func()
}

func NewMainPageViewModel() *MainPageViewModel {
	return &MainPageViewModel{
		StartLabelText: widget.NewLabel("Привет, мир!"),
		TestDataLabel: widget.NewLabel("Введите текст:"),
		TestDataEntry: binding.NewString(),
		SaveBtn: func() {},
	}
}

func (vm *MainPageViewModel) SetSaveBtnAction(action func()) {
	vm.SaveBtn = action
}

func (vm *MainPageViewModel) SaveNewLabelText() {
	text, _ := vm.TestDataEntry.Get()
	vm.StartLabelText.SetText(text)
}
