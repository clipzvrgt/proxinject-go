// Chinese translations
export default {
  // Navigation
  nav: {
    processes: '进程列表',
    settings: '设置',
    logs: '日志',
    help: '帮助',
  },
  
  // Process List
  process: {
    inject: '注入',
    injected: '已注入',
    ready: '就绪',
    proxy: '代理',
    noProxy: '无',
    logging: '日志记录',
    subprocess: '子进程',
    autoInject: '自动注入(Beta)',
    watchlistFilter: '监控列表',
    searchPlaceholder: '搜索进程...',
    noProcesses: '未找到进程',
    processName: '进程名称',
    pid: 'PID',
    status: '状态',
  },
  
  // Settings
  settings: {
    title: '设置',
    language: '语言',
    proxyProfiles: '代理配置',
    addNew: '新建',
    noProfiles: '暂无代理配置',
    watchlist: '监控列表',
    watchlistDesc: '匹配这些名称的进程将在检测到时自动注入。',
    add: '添加',
    exePlaceholder: '例如: chrome.exe, firefox*',
    // Profile dialog
    editProfile: '编辑配置',
    newProfile: '新建配置',
    profileName: '配置名称',
    address: '地址',
    port: '端口',
    username: '用户名 (可选)',
    password: '密码 (可选)',
    cancel: '取消',
    save: '保存',
  },
  
  // Logs
  logs: {
    title: '连接日志',
    autoScroll: '自动滚动',
    clear: '清空',
    noLogs: '暂无日志',
    enableLogging: '在主页面启用日志记录以查看连接日志。',
  },
  
  // Errors
  error: {
    initFailed: '初始化失败',
    dllMissing: '请确保 DLL 文件存在于 resources 目录',
  },
  
  // Tray
  tray: {
    show: '显示窗口',
    exit: '退出',
  },
  
  // Help
  help: {
    title: '使用方法',
    subtitle: '按照以下步骤使用 Proxinject',
    steps: {
      step1: '设置代理配置',
      step2: '设置监控列表',
      step3: '首页右上角设置代理服务器',
      step4: '勾选监控列表',
      step5: '打开 Antigravity 程序（或者其他满足监控列表的程序）',
      step6: '选中所有监控列表内的进程，并点击注入按钮',
      step7: '等待注入成功，成功加载对话框',
    },
  },
  
  // First Run
  firstRun: {
    title: '欢迎使用 Proxinject',
    message: '是否查看使用方法介绍？',
    dontShowAgain: '不再显示',
    viewGuide: '查看指南',
    skip: '跳过',
  },
};
