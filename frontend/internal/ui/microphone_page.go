package ui

import (
	"strconv"
	"strings"

	"github.com/Hirogava/WindowsAgent/frontend/internal/services"
	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMicrophonePageContainer(page *viewmodel.MicrophonePageViewModel, w fyne.Window) *fyne.Container {
	page.SetLoadMicrophonesBtnAction(func() {
		devices, err := page.ConfigService.ListMicrophones()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		page.SelectWidget.SetOptions(devices)

		cfg, cfgErr := page.ConfigService.LoadMicrophoneConfig()
		if cfgErr == nil {
			if cfg.Device != "" {
				page.SelectWidget.SetSelected(cfg.Device)
			}
			page.TriggerKeyEntry.Set(cfg.TriggerKey)
			page.DurationEntry.Set(strconv.Itoa(cfg.DurationSeconds))
		}

		page.StatusLabel.SetText("Список микрофонов обновлён")
	})

	page.SetSaveMicrophoneBtnAction(func() {
		selected := page.SelectWidget.Selected
		if selected == "" {
			dialog.ShowInformation("Info", "Выберите микрофон", w)
			return
		}

		triggerKey, _ := page.TriggerKeyEntry.Get()
		triggerKey = strings.ToLower(strings.TrimSpace(triggerKey))
		if triggerKey == "" {
			triggerKey = "space"
		}

		durationText, _ := page.DurationEntry.Get()
		duration, convErr := strconv.Atoi(strings.TrimSpace(durationText))
		if convErr != nil || duration <= 0 {
			duration = 5
		}

		err := page.ConfigService.SaveMicrophoneConfig(&services.MicrophoneConfig{
			Device:          selected,
			DurationSeconds: duration,
			TriggerKey:      triggerKey,
		})
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		page.StatusLabel.SetText("Микрофон сохранён")
		dialog.ShowInformation("OK", "Настройка микрофона сохранена", w)
	})

	page.SetRecordAudioSampleBtnAction(func() {
		cfg, cfgErr := page.ConfigService.LoadMicrophoneConfig()
		if cfgErr != nil {
			dialog.ShowError(cfgErr, w)
			return
		}

		path, err := page.MainService.RecordAudioSample(cfg.DurationSeconds)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		page.StatusLabel.SetText("Тестовая запись создана: " + path)
		dialog.ShowInformation("OK", "Тестовая запись завершена", w)
	})

	return container.NewVBox(
		page.SelectLabel,
		page.SelectWidget,
		page.TriggerKeyLabel,
		widget.NewEntryWithData(page.TriggerKeyEntry),
		page.DurationLabel,
		widget.NewEntryWithData(page.DurationEntry),
		widget.NewButton("Обновить микрофоны", func() {
			page.LoadMicrophonesBtn()
		}),
		widget.NewButton("Сохранить выбранный микрофон", func() {
			page.SaveMicrophoneBtn()
		}),
		widget.NewButton("Записать тестовое аудио", func() {
			page.RecordAudioSampleBtn()
		}),
		page.StatusLabel,
	)
}
