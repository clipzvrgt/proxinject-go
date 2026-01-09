package main

import (
	"context"
	"embed"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
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

	// Create tray manager
	trayManager := NewTrayManager(app)
	app.SetTrayManager(trayManager)

	// Setup systray with external loop (for integration with Wails)
	systrayStart, systrayEnd := systray.RunWithExternalLoop(func() {
		trayManager.SetupTray()
	}, func() {
		// Cleanup on systray exit
	})

	// Start systray
	systrayStart()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Proxinject",
		Width:     960,
		Height:    640,
		MinWidth:  800,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 15, A: 255}, // Dark theme background
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			app.ShutdownInjector()
			systrayEnd()
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			// Hide window instead of closing (minimize to tray)
			runtime.WindowHide(ctx)
			return true // Prevent the window from closing
		},
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
