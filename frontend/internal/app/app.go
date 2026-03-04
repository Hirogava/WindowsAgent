package app

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/services"
	"github.com/Hirogava/WindowsAgent/frontend/internal/ui"
	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func BuildApp() fyne.App {
	a := app.New()
	w := a.NewWindow("Fyne App")
	mainService := services.NewMainService()

	mainPageContainer := ui.NewMainPageContainer(viewmodel.NewMainPageViewModel(mainService), w)
	configPageContainer := ui.NewConfigPageContainer(viewmodel.NewConfigPageViewModel(), w)
	microphonePageContainer := ui.NewMicrophonePageContainer(viewmodel.NewMicrophonePageViewModel(mainService), w)
	menu := viewmodel.NewMenuViewModel(w)

	ui.NewMainPageBtn(menu, mainPageContainer)
	ui.NewConfigPageBtn(menu, configPageContainer)
	ui.NewMicrophonePageBtn(menu, microphonePageContainer)

	w.SetMainMenu(ui.NewMainMenu(menu))
	w.SetContent(mainPageContainer)
	w.Resize(fyne.NewSize(600, 400))
	w.Show()

	return a
}
