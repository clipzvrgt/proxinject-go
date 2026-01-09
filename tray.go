package main

import (
	_ "embed"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var iconData []byte

// TrayManager manages the system tray icon and menu
type TrayManager struct {
	app   *App
	mShow *systray.MenuItem
	mQuit *systray.MenuItem
}

// NewTrayManager creates a new tray manager
func NewTrayManager(app *App) *TrayManager {
	return &TrayManager{app: app}
}

// getTrayTexts returns localized tray texts based on app locale setting
func (t *TrayManager) getTrayTexts() (tooltip, showTitle, showTooltip, exitTitle, exitTooltip string) {
	locale := "en"
	if t.app != nil {
		t.app.mu.Lock()
		locale = t.app.locale
		t.app.mu.Unlock()
	}

	if locale == "zh" {
		return "Proxinject - 代理注入工具",
			"显示窗口", "显示主窗口",
			"退出", "退出程序"
	}
	return "Proxinject - Proxy Injection Tool",
		"Show Window", "Show the main window",
		"Exit", "Exit the application"
}

// SetupTray initializes the system tray
func (t *TrayManager) SetupTray() {
	tooltip, showTitle, showTooltip, exitTitle, exitTooltip := t.getTrayTexts()

	systray.SetIcon(iconData)
	systray.SetTitle("Proxinject")
	systray.SetTooltip(tooltip)

	// Add menu items
	t.mShow = systray.AddMenuItem(showTitle, showTooltip)
	systray.AddSeparator()
	t.mQuit = systray.AddMenuItem(exitTitle, exitTooltip)

	// Handle menu item clicks
	t.mShow.Click(func() {
		t.showWindow()
	})

	t.mQuit.Click(func() {
		t.quitApp()
	})

	// Handle tray icon clicks
	systray.SetOnClick(func(menu systray.IMenu) {
		t.showWindow()
	})

	systray.SetOnDClick(func(menu systray.IMenu) {
		t.showWindow()
	})

	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})
}

// UpdateTrayLanguage updates the tray menu text based on current locale
func (t *TrayManager) UpdateTrayLanguage() {
	tooltip, showTitle, showTooltip, exitTitle, exitTooltip := t.getTrayTexts()

	systray.SetTooltip(tooltip)
	if t.mShow != nil {
		t.mShow.SetTitle(showTitle)
		t.mShow.SetTooltip(showTooltip)
	}
	if t.mQuit != nil {
		t.mQuit.SetTitle(exitTitle)
		t.mQuit.SetTooltip(exitTooltip)
	}
}

// showWindow shows the main application window
func (t *TrayManager) showWindow() {
	if t.app != nil && t.app.ctx != nil {
		runtime.WindowShow(t.app.ctx)
	}
}

// quitApp quits the application
func (t *TrayManager) quitApp() {
	systray.Quit()
	if t.app != nil && t.app.ctx != nil {
		runtime.Quit(t.app.ctx)
	}
}
