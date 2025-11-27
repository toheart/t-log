package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Quick Capture",
		Width:  400,
		Height: 500,
		// StartHidden: true, // For MVP we might want to start visible to debug, but spec says hidden.
		// Keeping it visible for first run to ensure it works, then user can hide.
		// Actually spec FR-001 says start quietly. Let's set StartHidden: true
		StartHidden: true,
		Frameless:   true,
		AlwaysOnTop: true, // Will be managed by hotkey, but good default for popup
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Transparent background for Acrylic effect
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,
			DisableWindowIcon:    true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
