package ui

import (
	"encoding/json"
	"fmt"

	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewConfigPageContainer(page *viewmodel.ConfigPageViewModel, w fyne.Window) *fyne.Container {
	NewLoadConfigsButton(page, w)
	NewSelectWidget(page)
	NewSaveConfigsButton(page, w)

	container := container.NewVBox(
		widget.NewButton("Загрузить конфиги", func() { page.LoadConfigsButton() }),
		page.SelectLabel,
		page.SelectWidget,
		page.EntryLabel,
		widget.NewEntryWithData(page.ConfigEntry),
		widget.NewButton("Сохранить конфиг", func() { page.SaveConfigsButton() }),
	)

	return container
}

func NewLoadConfigsButton(page *viewmodel.ConfigPageViewModel, w fyne.Window) {
	page.SetLoadConfigsButtonAction(func() {
		files, err := page.ConfigService.LoadJsonFileNames()
		if err != nil || len(files) == 0 {
			dialog.ShowInformation("Info", "Ошибка загрузки конфиг файлов:", w)
		}

		page.SelectWidget.SetOptions(files)
		page.SelectWidget.Refresh()

		fmt.Println(files)
	})
}

func NewSelectWidget(page *viewmodel.ConfigPageViewModel) {
	page.SetSelectWidgetOptions([]string{}, func(selected string) {
		data, err := page.ConfigService.ReadConfigFile(selected)
		if err != nil {
			fmt.Println(err)
			return
		}

		var cfg map[string]interface{}
		json.Unmarshal(data, &cfg)
		cfgString := page.ConfigService.MapToJSONString(cfg)

		page.ConfigEntry.Set(fmt.Sprint(cfgString))
	})
}

func NewSaveConfigsButton(page *viewmodel.ConfigPageViewModel, w fyne.Window) {
	page.SetSaveConfigsButtonAction(func() {
		var parsed map[string]interface{}

		cfg, _ := page.ConfigEntry.Get()

		err := json.Unmarshal([]byte(cfg), &parsed)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		bytes, _ := json.MarshalIndent(parsed, "", "  ")

		currentFile := page.SelectWidget.Selected

		err = page.ConfigService.WriteConfigFile(currentFile, bytes)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("OK", "Файл сохранён", w)
	})
}
