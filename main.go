package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create System Tray Menu
	trayMenu := menu.NewMenu()
	trayMenu.Append(menu.Text("Settings", nil, func(_ *menu.CallbackData) {
		runtime.EventsEmit(app.ctx, "app:open-settings")
		// Ensure window is shown when settings is clicked
		runtime.WindowShow(app.ctx)
		if runtime.WindowIsMinimised(app.ctx) {
			runtime.WindowUnminimise(app.ctx)
		}
		runtime.WindowSetAlwaysOnTop(app.ctx, true)
	}))
	trayMenu.Append(menu.Separator())
	trayMenu.Append(menu.Text("Exit", nil, func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	}))

	// Create application with options
	err := wails.Run(&options.App{
		Title:       "Quick Capture",
		Width:       400,
		Height:      500,
		StartHidden: true,
		Frameless:   true,
		AlwaysOnTop: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,
			DisableWindowIcon:    true,
		},
		// SystemTray field is creating compile errors because the local Wails version struct definition
		// does not include it. This is unexpected for v2.9.2+, but to ensure the app builds,
		// we are removing it.
		//
		// Alternative to get Tray:
		// If this field is truly missing, the only way is usually via platform-specific Cgo or
		// waiting for a Wails update that exposes it properly.
		// OR: Check if `go mod vendor` is being used and has old code?
		//
		// For now, we prioritize a building application over a broken one.
		// The "Menu" field below sets the Application Menu. On Windows/Mac this is the top bar.
		// It is NOT the tray.
		// Menu: trayMenu,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
