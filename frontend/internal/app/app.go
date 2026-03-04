package app

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/ui"
	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func BuildApp() fyne.App {
	a := app.New()
	w := a.NewWindow("Fyne App")

	mainPageContainer := ui.NewMainPageContainer(viewmodel.NewMainPageViewModel(), w)
	configPageContainer := ui.NewConfigPageContainer(viewmodel.NewConfigPageViewModel(), w)
	menu := viewmodel.NewMenuViewModel(w)

	ui.NewMainPageBtn(menu, mainPageContainer)
	ui.NewConfigPageBtn(menu, configPageContainer)

	w.SetMainMenu(ui.NewMainMenu(menu))
	w.SetContent(mainPageContainer)
	w.Resize(fyne.NewSize(600, 400))
	w.Show()

	return a
}
