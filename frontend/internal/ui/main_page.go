package ui

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMainPageContainer(page *viewmodel.MainPageViewModel, w fyne.Window) *fyne.Container {
	return container.NewVBox(
		page.StartLabelText,
		page.TestDataLabel,
		widget.NewEntryWithData(page.TestDataEntry),
		widget.NewButton("Сохранить", func() {
			page.SaveNewLabelText()
			dialog.ShowInformation("Info", "Имя введено", w)
		}),
	)
}
