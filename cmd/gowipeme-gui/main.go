package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	"github.com/mat/gowipeme/internal/gui"
)

func main() {
	// Create an instance of the app structure
	app := gui.NewApp()

	// Create application with options
	// Note: Assets are handled by Wails via wails.json, not embedded here
	err := wails.Run(&options.App{
		Title:            "goWipeMe",
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
