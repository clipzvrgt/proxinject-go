import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import * as App from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

export interface ProcessEntry {
  pid: number;
  name: string;
  path: string;
  injected: boolean;
}

export interface ProxyProfile {
  id: string;
  name: string;
  address: string;
  port: number;
  username?: string;
  password?: string;
}

export interface AppConfig {
  profiles: ProxyProfile[];
  watchlist: string[];
  selectedProfileId: string | null;
}

const DEFAULT_CONFIG: AppConfig = {
  profiles: [
    { id: '1', name: 'Default Proxy', address: '127.0.0.1', port: 1080 }
  ],
  watchlist: [],
  selectedProfileId: null,
};

export const useInjectorStore = defineStore('injector', () => {
  // State
  const processes = ref<ProcessEntry[]>([]);
  const selectedPids = ref<Set<number>>(new Set());
  const injectedCount = ref(0);
  const loggingEnabled = ref(false);
  const subprocessEnabled = ref(false);
  const autoInjectEnabled = ref(false);
  const searchQuery = ref('');
  const watchlistFilterEnabled = ref(false);
  const config = ref<AppConfig>(loadConfig());
  const logs = ref<string[]>([]);
  const hasInjectedOnce = ref<boolean>(loadHasInjectedOnce());
  
  // Helper function to match process name against watchlist pattern
  function matchesWatchlistPattern(processName: string, pattern: string): boolean {
    // Support wildcard patterns like firefox*, *chrome*, *.exe
    const regexPattern = pattern
      .replace(/[.+^${}()|[\]\\]/g, '\\$&') // Escape special regex chars except *
      .replace(/\*/g, '.*'); // Replace * with .*
    const regex = new RegExp(`^${regexPattern}$`, 'i');
    return regex.test(processName);
  }
  
  // Check if a process matches any watchlist pattern
  function matchesWatchlist(processName: string): boolean {
    return config.value.watchlist.some(pattern => 
      matchesWatchlistPattern(processName, pattern)
    );
  }
  
  // Load config from localStorage
  function loadConfig(): AppConfig {
    try {
      const saved = localStorage.getItem('proxinject-config');
      if (saved) {
        return { ...DEFAULT_CONFIG, ...JSON.parse(saved) };
      }
    } catch (e) {
      console.error('Failed to load config:', e);
    }
    return { ...DEFAULT_CONFIG };
  }
  
  // Save config to localStorage
  function saveConfig() {
    localStorage.setItem('proxinject-config', JSON.stringify(config.value));
  }
  
  // Load hasInjectedOnce from localStorage
  function loadHasInjectedOnce(): boolean {
    return localStorage.getItem('proxinject-has-injected-once') === 'true';
  }
  
  // Save hasInjectedOnce to localStorage
  function saveHasInjectedOnce(value: boolean) {
    localStorage.setItem('proxinject-has-injected-once', value.toString());
  }
  
  // Computed
  const filteredProcesses = computed(() => {
    const query = searchQuery.value.toLowerCase();
    let result = processes.value;
    
    // Apply search query filter (only match process name)
    if (query) {
      result = result.filter(p => p.name.toLowerCase().includes(query));
    }
    
    // Apply watchlist filter if enabled
    if (watchlistFilterEnabled.value) {
      result = result.filter(p => matchesWatchlist(p.name));
    }
    
    return result;
  });
  
  const selectedProfile = computed(() => {
    return config.value.profiles.find(p => p.id === config.value.selectedProfileId);
  });
  
  // Actions
  async function refreshProcesses() {
    try {
      const result = await App.ListProcesses();
      processes.value = result as ProcessEntry[];
      injectedCount.value = processes.value.filter(p => p.injected).length;
    } catch (e) {
      console.error('Failed to refresh processes:', e);
    }
  }
  
  async function injectSelected() {
    if (selectedPids.value.size === 0) return;
    
    try {
      const pids = Array.from(selectedPids.value);
      await App.InjectProcesses(pids);
      selectedPids.value.clear();
      await refreshProcesses();
      
      // First injection: activate auto-inject feature
      if (!hasInjectedOnce.value) {
        hasInjectedOnce.value = true;
        saveHasInjectedOnce(true);
        
        // Auto-enable auto-inject if watchlist is configured
        if (config.value.watchlist.length > 0 && !autoInjectEnabled.value) {
          const patterns = config.value.watchlist.join(';');
          const result = await App.StartAutoInject(patterns);
          autoInjectEnabled.value = result;
        }
      }
    } catch (e) {
      console.error('Failed to inject:', e);
    }
  }
  
  async function setProxy(profile: ProxyProfile) {
    try {
      await App.SetProxy(profile);
      config.value.selectedProfileId = profile.id;
      saveConfig();
    } catch (e) {
      console.error('Failed to set proxy:', e);
    }
  }
  
  async function setLogging(enable: boolean) {
    try {
      await App.SetLogging(enable);
      loggingEnabled.value = enable;
    } catch (e) {
      console.error('Failed to set logging:', e);
    }
  }
  
  async function setSubprocess(enable: boolean) {
    try {
      await App.SetSubprocess(enable);
      subprocessEnabled.value = enable;
    } catch (e) {
      console.error('Failed to set subprocess:', e);
    }
  }
  
  async function toggleAutoInject() {
    try {
      if (autoInjectEnabled.value) {
        await App.StopAutoInject();
        autoInjectEnabled.value = false;
      } else {
        const patterns = config.value.watchlist.join(';');
        if (patterns) {
          const result = await App.StartAutoInject(patterns);
          autoInjectEnabled.value = result;
        }
      }
    } catch (e) {
      console.error('Failed to toggle auto-inject:', e);
    }
  }
  
  // Profile management
  function addProfile(profile: Omit<ProxyProfile, 'id'>) {
    config.value.profiles.push({
      ...profile,
      id: Date.now().toString(),
    });
    saveConfig();
  }
  
  function updateProfile(id: string, updates: Partial<ProxyProfile>) {
    const idx = config.value.profiles.findIndex(p => p.id === id);
    if (idx !== -1) {
      config.value.profiles[idx] = { ...config.value.profiles[idx], ...updates };
      saveConfig();
    }
  }
  
  function deleteProfile(id: string) {
    config.value.profiles = config.value.profiles.filter(p => p.id !== id);
    if (config.value.selectedProfileId === id) {
      config.value.selectedProfileId = null;
    }
    saveConfig();
  }
  
  // Watchlist management
  function addToWatchlist(exe: string) {
    if (!config.value.watchlist.includes(exe)) {
      config.value.watchlist.push(exe);
      saveConfig();
    }
  }
  
  function removeFromWatchlist(exe: string) {
    config.value.watchlist = config.value.watchlist.filter(e => e !== exe);
    saveConfig();
  }
  
  // Add log entry
  function addLog(entry: string) {
    logs.value.push(`[${new Date().toLocaleTimeString()}] ${entry}`);
    if (logs.value.length > 1000) {
      logs.value = logs.value.slice(-500);
    }
  }
  
  // Event listeners for auto-inject notifications
  let eventListenersInitialized = false;
  
  async function initEventListeners() {
    if (eventListenersInitialized) return;
    
    // Listen for successful injections
    EventsOn('process-injected', (data: { pid: number; name: string }) => {
      console.log('Process injected:', data);
      refreshProcesses();
      addLog(`Auto-injected: ${data.name} (PID: ${data.pid})`);
    });
    
    // Listen for injection failures
    EventsOn('process-injection-failed', (data: { pid: number; name: string; attempts: number }) => {
      console.log('Process injection failed:', data);
      refreshProcesses();
      addLog(`Injection failed: ${data.name} (PID: ${data.pid}) after ${data.attempts} attempts`);
    });
    
    eventListenersInitialized = true;
  }
  
  function cleanupEventListeners() {
    EventsOff('process-injected');
    EventsOff('process-injection-failed');
    eventListenersInitialized = false;
  }
  
  return {
    // State
    processes,
    selectedPids,
    injectedCount,
    loggingEnabled,
    subprocessEnabled,
    autoInjectEnabled,
    hasInjectedOnce,
    searchQuery,
    watchlistFilterEnabled,
    config,
    logs,
    
    // Computed
    filteredProcesses,
    selectedProfile,
    
    // Actions
    refreshProcesses,
    injectSelected,
    setProxy,
    setLogging,
    setSubprocess,
    toggleAutoInject,
    addProfile,
    updateProfile,
    deleteProfile,
    addToWatchlist,
    removeFromWatchlist,
    addLog,
    initEventListeners,
    cleanupEventListeners,
  };
});
