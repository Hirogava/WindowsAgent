package viewmodel

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/services"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type MicrophonePageViewModel struct {
	SelectLabel          *widget.Label
	SelectWidget         *widget.Select
	TriggerKeyLabel      *widget.Label
	TriggerKeyEntry      binding.String
	DurationLabel        *widget.Label
	DurationEntry        binding.String
	StatusLabel          *widget.Label
	LoadMicrophonesBtn   func()
	SaveMicrophoneBtn    func()
	RecordAudioSampleBtn func()

	ConfigService *services.ConfigService
	MainService   *services.MainService
}

func NewMicrophonePageViewModel(mainService *services.MainService) *MicrophonePageViewModel {
	return &MicrophonePageViewModel{
		SelectLabel:          widget.NewLabel("Выберите микрофон:"),
		SelectWidget:         widget.NewSelect([]string{}, func(string) {}),
		TriggerKeyLabel:      widget.NewLabel("Клавиша для старта записи (например: space, enter, a):"),
		TriggerKeyEntry:      binding.NewString(),
		DurationLabel:        widget.NewLabel("Длительность записи (сек):"),
		DurationEntry:        binding.NewString(),
		StatusLabel:          widget.NewLabel(""),
		LoadMicrophonesBtn:   func() {},
		SaveMicrophoneBtn:    func() {},
		RecordAudioSampleBtn: func() {},
		ConfigService:        services.NewConfigService(),
		MainService:          mainService,
	}
}

func (vm *MicrophonePageViewModel) SetLoadMicrophonesBtnAction(action func()) {
	vm.LoadMicrophonesBtn = action
}

func (vm *MicrophonePageViewModel) SetSaveMicrophoneBtnAction(action func()) {
	vm.SaveMicrophoneBtn = action
}

func (vm *MicrophonePageViewModel) SetRecordAudioSampleBtnAction(action func()) {
	vm.RecordAudioSampleBtn = action
}
