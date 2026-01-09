package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ProxyProfile represents a proxy configuration
type ProxyProfile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// ProcessEntry represents a process for the frontend
type ProcessEntry struct {
	PID      uint32 `json:"pid"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Injected bool   `json:"injected"`
}

// AppState represents the current application state
type AppState struct {
	Initialized       bool    `json:"initialized"`
	CurrentProxy      *string `json:"currentProxy"`
	LoggingEnabled    bool    `json:"loggingEnabled"`
	SubprocessEnabled bool    `json:"subprocessEnabled"`
	AutoInjectEnabled bool    `json:"autoInjectEnabled"`
	InjectedCount     int     `json:"injectedCount"`
}

// ProcessInjectedPayload is emitted when a process is injected
type ProcessInjectedPayload struct {
	PID  uint32 `json:"pid"`
	Name string `json:"name"`
}

// App struct
type App struct {
	ctx context.Context
	mu  sync.Mutex

	initialized       bool
	currentProxy      *string
	loggingEnabled    bool
	subprocessEnabled bool
	autoInjectEnabled bool
	watchPatterns     []string
	locale            string
	trayManager       *TrayManager

	// Monitor control
	monitorRunning atomic.Bool
	monitorStop    chan struct{}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{locale: "en"}
}

// SetTrayManager sets the tray manager reference
func (a *App) SetTrayManager(tm *TrayManager) {
	a.trayManager = tm
}

// SetLocale sets the application locale and updates tray menu
func (a *App) SetLocale(locale string) {
	a.mu.Lock()
	a.locale = locale
	a.mu.Unlock()

	// Update tray menu language
	if a.trayManager != nil {
		a.trayManager.UpdateTrayLanguage()
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// InitInjector initializes the DLL and injector
func (a *App) InitInjector() (uint16, error) {
	// Get executable directory for DLL path
	exePath, err := os.Executable()
	if err != nil {
		return 0, err
	}
	exeDir := filepath.Dir(exePath)
	
	// Try multiple DLL locations
	dllPaths := []string{
		filepath.Join(exeDir, "resources", "proxinject_ffi.dll"),
		filepath.Join(exeDir, "proxinject_ffi.dll"),
		// For development - try relative to project
		"resources/proxinject_ffi.dll",
	}
	
	var loadErr error
	for _, dllPath := range dllPaths {
		if _, err := os.Stat(dllPath); err == nil {
			loadErr = ffi.LoadLibrary(dllPath)
			if loadErr == nil {
				break
			}
		}
	}
	
	if !ffi.loaded {
		if loadErr != nil {
			return 0, loadErr
		}
		return 0, os.ErrNotExist
	}
	
	port, err := ffi.Init()
	if err != nil {
		return 0, err
	}
	
	a.mu.Lock()
	a.initialized = true
	a.mu.Unlock()
	
	return port, nil
}

// ShutdownInjector shuts down the injector
func (a *App) ShutdownInjector() error {
	a.StopAutoInject()
	
	err := ffi.Shutdown()
	
	a.mu.Lock()
	a.initialized = false
	a.mu.Unlock()
	
	return err
}

// ListProcesses returns a list of running processes
func (a *App) ListProcesses() ([]ProcessEntry, error) {
	processes, err := ffi.ListProcesses()
	if err != nil {
		return nil, err
	}
	
	result := make([]ProcessEntry, 0, len(processes))
	for _, p := range processes {
		if p.PID == 0 {
			continue
		}
		result = append(result, ProcessEntry{
			PID:      p.PID,
			Name:     p.Name,
			Path:     p.Path,
			Injected: p.Injected,
		})
	}
	
	return result, nil
}

// InjectProcess injects a single process
func (a *App) InjectProcess(pid uint32) (bool, error) {
	return ffi.Inject(pid)
}

// InjectProcesses injects multiple processes
func (a *App) InjectProcesses(pids []uint32) (map[uint32]bool, error) {
	results := make(map[uint32]bool)
	for _, pid := range pids {
		success, _ := ffi.Inject(pid)
		results[pid] = success
	}
	return results, nil
}

// SetProxy sets the proxy configuration
func (a *App) SetProxy(profile ProxyProfile) (bool, error) {
	result, err := ffi.SetProxy(profile.Address, profile.Port, profile.Username, profile.Password)
	if err != nil {
		return false, err
	}
	
	if result {
		a.mu.Lock()
		a.currentProxy = &profile.ID
		a.mu.Unlock()
	}
	
	return result, nil
}

// ClearProxy clears the proxy configuration
func (a *App) ClearProxy() error {
	err := ffi.ClearProxy()
	if err != nil {
		return err
	}
	
	a.mu.Lock()
	a.currentProxy = nil
	a.mu.Unlock()
	
	return nil
}

// SetLogging enables or disables logging
func (a *App) SetLogging(enable bool) error {
	err := ffi.EnableLog(enable)
	if err != nil {
		return err
	}
	
	a.mu.Lock()
	a.loggingEnabled = enable
	a.mu.Unlock()
	
	return nil
}

// SetSubprocess enables or disables subprocess injection
func (a *App) SetSubprocess(enable bool) error {
	err := ffi.EnableSubprocess(enable)
	if err != nil {
		return err
	}
	
	a.mu.Lock()
	a.subprocessEnabled = enable
	a.mu.Unlock()
	
	return nil
}

// StartAutoInject starts auto-injection with patterns
func (a *App) StartAutoInject(patterns string) (bool, error) {
	// Parse patterns
	patternList := make([]string, 0)
	for _, p := range strings.Split(patterns, ";") {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			patternList = append(patternList, p)
		}
	}
	
	a.mu.Lock()
	a.watchPatterns = patternList
	a.autoInjectEnabled = true
	a.mu.Unlock()
	
	// Start monitor if not running
	if !a.monitorRunning.Load() {
		a.monitorRunning.Store(true)
		a.monitorStop = make(chan struct{})
		go a.monitorLoop()
	}
	
	return true, nil
}

// StopAutoInject stops auto-injection
func (a *App) StopAutoInject() error {
	a.monitorRunning.Store(false)
	if a.monitorStop != nil {
		close(a.monitorStop)
		a.monitorStop = nil
	}
	
	a.mu.Lock()
	a.autoInjectEnabled = false
	a.mu.Unlock()
	
	return nil
}

// GetState returns the current application state
func (a *App) GetState() AppState {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	injectedCount, _ := ffi.GetInjectedCount()
	
	return AppState{
		Initialized:       a.initialized,
		CurrentProxy:      a.currentProxy,
		LoggingEnabled:    a.loggingEnabled,
		SubprocessEnabled: a.subprocessEnabled,
		AutoInjectEnabled: a.autoInjectEnabled,
		InjectedCount:     injectedCount,
	}
}

// monitorLoop runs the auto-injection monitor
func (a *App) monitorLoop() {
	knownInjected := make(map[uint32]bool)
	failedAttempts := make(map[uint32]bool)
	pendingInject := make(map[uint32]struct {
		name string
		time time.Time
	})
	
	const injectDelay = 1500 * time.Millisecond
	
	// Initialize with currently injected processes
	if processes, err := ffi.ListProcesses(); err == nil {
		for _, p := range processes {
			if p.Injected {
				knownInjected[p.PID] = true
			}
		}
	}
	
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-a.monitorStop:
			return
		case <-ticker.C:
			if !a.monitorRunning.Load() {
				return
			}
			
			// Get current patterns
			a.mu.Lock()
			patterns := append([]string{}, a.watchPatterns...)
			a.mu.Unlock()
			
			now := time.Now()
			
			// Process pending injections
			for pid, info := range pendingInject {
				if now.Sub(info.time) >= injectDelay {
					delete(pendingInject, pid)
					
					// Verify pattern still matches
					nameLower := strings.ToLower(info.name)
					matches := false
					for _, p := range patterns {
						if strings.Contains(nameLower, p) {
							matches = true
							break
						}
					}
					
					if matches {
						// Try injection with retry
						success := false
						for attempt := 0; attempt < 3; attempt++ {
							if ok, _ := ffi.Inject(pid); ok {
								success = true
								break
							}
							time.Sleep(100 * time.Millisecond)
						}
						
						if success {
							knownInjected[pid] = true
							runtime.EventsEmit(a.ctx, "process-injected", ProcessInjectedPayload{
								PID:  pid,
								Name: info.name,
							})
						} else {
							failedAttempts[pid] = true
							runtime.EventsEmit(a.ctx, "process-injection-failed", map[string]interface{}{
								"pid":      pid,
								"name":     info.name,
								"attempts": 3,
							})
						}
					}
				}
			}
			
			// Scan for new processes
			processes, err := ffi.ListProcesses()
			if err != nil {
				continue
			}
			
			for _, p := range processes {
				if knownInjected[p.PID] || failedAttempts[p.PID] {
					continue
				}
				if _, pending := pendingInject[p.PID]; pending {
					continue
				}
				
				if p.Injected {
					knownInjected[p.PID] = true
					runtime.EventsEmit(a.ctx, "process-injected", ProcessInjectedPayload{
						PID:  p.PID,
						Name: p.Name,
					})
				} else if p.PID != 0 {
					nameLower := strings.ToLower(p.Name)
					for _, pattern := range patterns {
						if strings.Contains(nameLower, pattern) {
							pendingInject[p.PID] = struct {
								name string
								time time.Time
							}{p.Name, now}
							break
						}
					}
				}
			}
		}
	}
}
