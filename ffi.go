package main

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

// ProcessInfo represents a process entry from the DLL
type ProcessInfo struct {
	PID      uint32
	Name     string
	Path     string
	Injected bool
}

// FFI manages the DLL interface
type FFI struct {
	dll       *syscall.DLL
	mu        sync.Mutex
	loaded    bool
	
	// Function pointers
	procInit             *syscall.Proc
	procShutdown         *syscall.Proc
	procInject           *syscall.Proc
	procSetProxy         *syscall.Proc
	procClearProxy       *syscall.Proc
	procEnableLog        *syscall.Proc
	procEnableSubprocess *syscall.Proc
	procListProcesses    *syscall.Proc
	procGetInjectedCount *syscall.Proc
	procStartMonitor     *syscall.Proc
	procStopMonitor      *syscall.Proc
	procIsMonitoring     *syscall.Proc
}

// C struct for process info (matches proxinject_ffi.dll)
type ffiProcessInfo struct {
	PID      uint32
	Name     [260]byte
	Path     [520]byte
	Injected int32
}

var ffi = &FFI{}

// LoadLibrary loads the DLL from the given path
func (f *FFI) LoadLibrary(dllPath string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if f.loaded {
		return nil
	}
	
	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		return fmt.Errorf("failed to load DLL: %w", err)
	}
	
	f.dll = dll
	
	// Load all function pointers
	f.procInit, _ = dll.FindProc("proxinject_init")
	f.procShutdown, _ = dll.FindProc("proxinject_shutdown")
	f.procInject, _ = dll.FindProc("proxinject_inject")
	f.procSetProxy, _ = dll.FindProc("proxinject_set_proxy")
	f.procClearProxy, _ = dll.FindProc("proxinject_clear_proxy")
	f.procEnableLog, _ = dll.FindProc("proxinject_enable_log")
	f.procEnableSubprocess, _ = dll.FindProc("proxinject_enable_subprocess")
	f.procListProcesses, _ = dll.FindProc("proxinject_list_processes")
	f.procGetInjectedCount, _ = dll.FindProc("proxinject_get_injected_count")
	f.procStartMonitor, _ = dll.FindProc("proxinject_start_monitor")
	f.procStopMonitor, _ = dll.FindProc("proxinject_stop_monitor")
	f.procIsMonitoring, _ = dll.FindProc("proxinject_is_monitoring")
	
	f.loaded = true
	return nil
}

// Init initializes the injector and returns the port
func (f *FFI) Init() (uint16, error) {
	if f.procInit == nil {
		return 0, fmt.Errorf("init function not found")
	}
	
	ret, _, _ := f.procInit.Call()
	port := uint16(ret)
	if port == 0 {
		return 0, fmt.Errorf("failed to initialize")
	}
	return port, nil
}

// Shutdown shuts down the injector
func (f *FFI) Shutdown() error {
	if f.procShutdown == nil {
		return fmt.Errorf("shutdown function not found")
	}
	
	f.procShutdown.Call()
	return nil
}

// Inject injects a process by PID
func (f *FFI) Inject(pid uint32) (bool, error) {
	if f.procInject == nil {
		return false, fmt.Errorf("inject function not found")
	}
	
	ret, _, _ := f.procInject.Call(uintptr(pid))
	return ret != 0, nil
}

// SetProxy sets the proxy configuration
func (f *FFI) SetProxy(address string, port uint16, username, password string) (bool, error) {
	if f.procSetProxy == nil {
		return false, fmt.Errorf("set_proxy function not found")
	}
	
	addrPtr, _ := syscall.BytePtrFromString(address)
	
	var userPtr, passPtr *byte
	if username != "" {
		userPtr, _ = syscall.BytePtrFromString(username)
	}
	if password != "" {
		passPtr, _ = syscall.BytePtrFromString(password)
	}
	
	ret, _, _ := f.procSetProxy.Call(
		uintptr(unsafe.Pointer(addrPtr)),
		uintptr(port),
		uintptr(unsafe.Pointer(userPtr)),
		uintptr(unsafe.Pointer(passPtr)),
	)
	return ret != 0, nil
}

// ClearProxy clears the proxy configuration
func (f *FFI) ClearProxy() error {
	if f.procClearProxy == nil {
		return fmt.Errorf("clear_proxy function not found")
	}
	
	f.procClearProxy.Call()
	return nil
}

// EnableLog enables or disables logging
func (f *FFI) EnableLog(enable bool) error {
	if f.procEnableLog == nil {
		return fmt.Errorf("enable_log function not found")
	}
	
	var val uintptr
	if enable {
		val = 1
	}
	f.procEnableLog.Call(val)
	return nil
}

// EnableSubprocess enables or disables subprocess injection
func (f *FFI) EnableSubprocess(enable bool) error {
	if f.procEnableSubprocess == nil {
		return fmt.Errorf("enable_subprocess function not found")
	}
	
	var val uintptr
	if enable {
		val = 1
	}
	f.procEnableSubprocess.Call(val)
	return nil
}

// ListProcesses returns a list of all processes
func (f *FFI) ListProcesses() ([]ProcessInfo, error) {
	if f.procListProcesses == nil {
		return nil, fmt.Errorf("list_processes function not found")
	}
	
	// First call to get count
	ret, _, _ := f.procListProcesses.Call(0, 0)
	count := int(ret)
	if count <= 0 {
		return []ProcessInfo{}, nil
	}
	
	// Allocate buffer
	buffer := make([]ffiProcessInfo, count)
	
	// Second call to fill buffer
	ret, _, _ = f.procListProcesses.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(count),
	)
	actual := int(ret)
	if actual < count {
		count = actual
	}
	
	// Convert to Go types
	result := make([]ProcessInfo, 0, count)
	for i := 0; i < count; i++ {
		info := &buffer[i]
		if info.PID == 0 {
			continue
		}
		
		// Convert null-terminated strings
		name := bytesToString(info.Name[:])
		path := bytesToString(info.Path[:])
		
		result = append(result, ProcessInfo{
			PID:      info.PID,
			Name:     name,
			Path:     path,
			Injected: info.Injected != 0,
		})
	}
	
	return result, nil
}

// GetInjectedCount returns the count of injected processes
func (f *FFI) GetInjectedCount() (int, error) {
	if f.procGetInjectedCount == nil {
		return 0, fmt.Errorf("get_injected_count function not found")
	}
	
	ret, _, _ := f.procGetInjectedCount.Call()
	return int(ret), nil
}

// StartMonitor starts monitoring for processes matching patterns
func (f *FFI) StartMonitor(patterns string) (bool, error) {
	if f.procStartMonitor == nil {
		return false, fmt.Errorf("start_monitor function not found")
	}
	
	patternsPtr, _ := syscall.BytePtrFromString(patterns)
	ret, _, _ := f.procStartMonitor.Call(uintptr(unsafe.Pointer(patternsPtr)))
	return ret != 0, nil
}

// StopMonitor stops the process monitor
func (f *FFI) StopMonitor() error {
	if f.procStopMonitor == nil {
		return fmt.Errorf("stop_monitor function not found")
	}
	
	f.procStopMonitor.Call()
	return nil
}

// IsMonitoring returns whether monitoring is active
func (f *FFI) IsMonitoring() (bool, error) {
	if f.procIsMonitoring == nil {
		return false, fmt.Errorf("is_monitoring function not found")
	}
	
	ret, _, _ := f.procIsMonitoring.Call()
	return ret != 0, nil
}

// bytesToString converts a null-terminated byte array to string
func bytesToString(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
