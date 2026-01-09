// English translations
export default {
  // Navigation
  nav: {
    processes: 'Processes',
    settings: 'Settings',
    logs: 'Logs',
    help: 'Help',
  },
  
  // Process List
  process: {
    inject: 'Inject',
    injected: 'Injected',
    ready: 'Ready',
    proxy: 'Proxy',
    noProxy: 'None',
    logging: 'Logging',
    subprocess: 'Subprocess',
    autoInject: 'Auto-Inject(Beta)',
    watchlistFilter: 'Watchlist',
    searchPlaceholder: 'Search processes...',
    noProcesses: 'No processes found',
    processName: 'Process Name',
    pid: 'PID',
    status: 'Status',
  },
  
  // Settings
  settings: {
    title: 'Settings',
    language: 'Language',
    proxyProfiles: 'Proxy Profiles',
    addNew: 'Add New',
    noProfiles: 'No proxy profiles configured',
    watchlist: 'Executable Watchlist',
    watchlistDesc: 'Processes matching these names will be auto-injected when detected.',
    add: 'Add',
    exePlaceholder: 'e.g., chrome.exe, firefox*',
    // Profile dialog
    editProfile: 'Edit Profile',
    newProfile: 'New Profile',
    profileName: 'Profile Name',
    address: 'Address',
    port: 'Port',
    username: 'Username (optional)',
    password: 'Password (optional)',
    cancel: 'Cancel',
    save: 'Save',
  },
  
  // Logs
  logs: {
    title: 'Connection Logs',
    autoScroll: 'Auto-scroll',
    clear: 'Clear',
    noLogs: 'No logs yet.',
    enableLogging: 'Enable logging in the main page to see connection logs.',
  },
  
  // Errors
  error: {
    initFailed: 'Initialization Failed',
    dllMissing: 'Please ensure the DLL file exists in the resources directory',
  },
  
  // Tray
  tray: {
    show: 'Show Window',
    exit: 'Exit',
  },
  
  // Help
  help: {
    title: 'Getting Started',
    subtitle: 'Follow these steps to use Proxinject',
    steps: {
      step1: 'Configure proxy settings',
      step2: 'Set up the watchlist',
      step3: 'Select proxy server in the top right corner of the home page',
      step4: 'Enable watchlist filter',
      step5: 'Open Antigravity program (or other programs matching the watchlist)',
      step6: 'Select all processes in the watchlist and click the Inject button',
      step7: 'Wait for injection to complete successfully',
    },
  },
  
  // First Run
  firstRun: {
    title: 'Welcome to Proxinject',
    message: 'Would you like to view the getting started guide?',
    dontShowAgain: "Don't show again",
    viewGuide: 'View Guide',
    skip: 'Skip',
  },
};
