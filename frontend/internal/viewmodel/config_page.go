package viewmodel

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/services"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ConfigPageViewModel struct {
	LoadConfigsButton func()
	SelectLabel       *widget.Label
	SelectWidget      *widget.Select
	EntryLabel        *widget.Label
	ConfigEntry       binding.String
	SaveConfigsButton func()

	ConfigService *services.ConfigService
}

func NewConfigPageViewModel() *ConfigPageViewModel {
	return &ConfigPageViewModel{
		LoadConfigsButton: func() {},
		SelectLabel:       widget.NewLabel("Выберите конфиг для изменения:"),
		SelectWidget:      widget.NewSelect([]string{}, func(string) {}),
		EntryLabel:        widget.NewLabel("Введите новые конфиги"),
		ConfigEntry:       binding.NewString(),
		SaveConfigsButton: func() {},
		ConfigService:     services.NewConfigService(),
	}
}

func (vm *ConfigPageViewModel) SetLoadConfigsButtonAction(action func()) {
	vm.LoadConfigsButton = action
}

func (vm *ConfigPageViewModel) SetSelectWidgetOptions(options []string, onSelect func(string)) {
	vm.SelectWidget.Options = options
	vm.SelectWidget.OnChanged = onSelect
}

func (vm *ConfigPageViewModel) SetSaveConfigsButtonAction(action func()) {
	vm.SaveConfigsButton = action
}
