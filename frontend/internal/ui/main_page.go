package ui

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMainPageContainer(page *viewmodel.MainPageViewModel, w fyne.Window) *fyne.Container {
	page.SetStartServicesBtnAction(func() {
		err := page.MainService.StartAllServices()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		page.StatusLabel.SetText("Сервисы запущены")
	})

	page.SetStopServicesBtnAction(func() {
		err := page.MainService.StopAllServices()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		page.StatusLabel.SetText("Сервисы остановлены")
	})

	page.SetRecordAudioBtnAction(func() {
		path, err := page.MainService.RecordAudioSample(5)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		page.RecordStatusLabel.SetText("Аудио записано: " + path)
		dialog.ShowInformation("OK", "Запись завершена", w)
	})

	return container.NewVBox(
		page.TitleLabel,
		page.StatusLabel,
		widget.NewButton("Запустить сервисы", func() {
			page.StartServicesBtn()
		}),
		widget.NewButton("Остановить сервисы", func() {
			page.StopServicesBtn()
		}),
		widget.NewButton("Записать аудио", func() {
			page.RecordAudioBtn()
		}),
		page.RecordStatusLabel,
	)
}
