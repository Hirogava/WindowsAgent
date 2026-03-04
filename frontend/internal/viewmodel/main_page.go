package viewmodel

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/services"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type MainPageViewModel struct {
	TitleLabel        *widget.Label
	StatusLabel       *widget.Label
	RecordStatusLabel *widget.Label
	TestDataEntry     binding.String
	StartServicesBtn  func()
	StopServicesBtn   func()
	RecordAudioBtn    func()

	MainService *services.MainService
}

func NewMainPageViewModel(mainService *services.MainService) *MainPageViewModel {
	return &MainPageViewModel{
		TitleLabel:        widget.NewLabel("Управление агентом"),
		StatusLabel:       widget.NewLabel("Сервисы остановлены"),
		RecordStatusLabel: widget.NewLabel(""),
		TestDataEntry:     binding.NewString(),
		StartServicesBtn:  func() {},
		StopServicesBtn:   func() {},
		RecordAudioBtn:    func() {},
		MainService:       mainService,
	}
}

func (vm *MainPageViewModel) SetStartServicesBtnAction(action func()) {
	vm.StartServicesBtn = action
}

func (vm *MainPageViewModel) SetStopServicesBtnAction(action func()) {
	vm.StopServicesBtn = action
}

func (vm *MainPageViewModel) SetRecordAudioBtnAction(action func()) {
	vm.RecordAudioBtn = action
}
