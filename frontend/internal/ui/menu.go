package ui

import (
	"github.com/Hirogava/WindowsAgent/frontend/internal/viewmodel"

	"fyne.io/fyne/v2"
)

func NewMainMenu(menuM *viewmodel.MenuViewModel) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("Меню",
			fyne.NewMenuItem("Главная страница", menuM.MainPageBtn),
			fyne.NewMenuItem("Настройки конфигов", menuM.ConfigPageBtn),
		),
	)
}

func NewMainPageBtn(menuM *viewmodel.MenuViewModel, container *fyne.Container) {
	menuM.SetMainPageBtnAction(func() {
		menuM.MainWindow.SetContent(container)
	})
}

func NewConfigPageBtn(menuM *viewmodel.MenuViewModel, container *fyne.Container) {
	menuM.SetConfigPageBtnAction(func() {
		menuM.MainWindow.SetContent(container)
	})
}
