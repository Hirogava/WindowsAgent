package main

import "github.com/Hirogava/WindowsAgent/frontend/internal/app"

func main() {
	a := app.BuildApp()
	a.Run()
}