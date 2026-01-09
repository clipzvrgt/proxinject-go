<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { InitInjector, ShutdownInjector } from '../wailsjs/go/main/App';
import Sidebar from './components/Sidebar.vue';
import ProcessList from './components/ProcessList.vue';
import Settings from './components/Settings.vue';
import Logs from './components/Logs.vue';
import Help from './components/Help.vue';
import { useInjectorStore } from './stores/injector';
import { useI18n } from './i18n';

const store = useInjectorStore();
const { t } = useI18n();
const currentPage = ref<'main' | 'settings' | 'logs' | 'help'>('main');
const initialized = ref(false);
const error = ref<string | null>(null);

// First run dialog state
const showFirstRunDialog = ref(false);
const dontShowAgain = ref(false);

// Check if this is the first run
function isFirstRun(): boolean {
  return localStorage.getItem('proxinject-first-run-done') !== 'true';
}

// Mark first run as complete
function markFirstRunDone() {
  if (dontShowAgain.value) {
    localStorage.setItem('proxinject-first-run-done', 'true');
  }
}

// Handle view guide button
function handleViewGuide() {
  markFirstRunDone();
  showFirstRunDialog.value = false;
  currentPage.value = 'help';
}

// Handle skip button
function handleSkip() {
  markFirstRunDone();
  showFirstRunDialog.value = false;
}

onMounted(async () => {
  try {
    await InitInjector();
    initialized.value = true;
    await store.initEventListeners();
    
    // Check for first run
    if (isFirstRun()) {
      currentPage.value = 'help';
      showFirstRunDialog.value = true;
    }
  } catch (e) {
    error.value = String(e);
    console.error('Failed to initialize:', e);
  }
});

onUnmounted(async () => {
  try {
    await ShutdownInjector();
  } catch (e) {
    console.error('Failed to shutdown:', e);
  }
});
</script>

<template>
  <div class="app-container">
    <Sidebar v-model:current="currentPage" />
    
    <div class="main-content">
      <template v-if="error">
        <div class="error-container">
          <h2>{{ t('error.initFailed') }}</h2>
          <p>{{ error }}</p>
          <p>{{ t('error.dllMissing') }}</p>
        </div>
      </template>
      
      <template v-else>
        <ProcessList v-if="currentPage === 'main'" />
        <Settings v-else-if="currentPage === 'settings'" />
        <Logs v-else-if="currentPage === 'logs'" />
        <Help v-else-if="currentPage === 'help'" />
      </template>
    </div>
    
    <!-- First Run Welcome Dialog -->
    <div v-if="showFirstRunDialog" class="dialog-overlay">
      <div class="dialog first-run-dialog">
        <div class="dialog-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
        </div>
        <h3>{{ t('firstRun.title') }}</h3>
        <p class="dialog-message">{{ t('firstRun.message') }}</p>
        
        <label class="checkbox-label">
          <input type="checkbox" v-model="dontShowAgain" />
          <span>{{ t('firstRun.dontShowAgain') }}</span>
        </label>
        
        <div class="dialog-actions">
          <button class="btn-secondary" @click="handleSkip">{{ t('firstRun.skip') }}</button>
          <button class="btn-primary" @click="handleViewGuide">{{ t('firstRun.viewGuide') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: var(--error);
}

.error-container h2 {
  margin-bottom: 16px;
}

.error-container p {
  color: var(--text-secondary);
  max-width: 400px;
}

/* Dialog Styles */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.dialog {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px;
  width: 420px;
  max-width: 90%;
}

.first-run-dialog {
  text-align: center;
}

.dialog-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
  border-radius: 50%;
}

.dialog-icon svg {
  width: 32px;
  height: 32px;
  color: white;
}

.dialog h3 {
  margin-bottom: 12px;
  font-size: 20px;
}

.dialog-message {
  color: var(--text-secondary);
  margin-bottom: 20px;
  font-size: 14px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 20px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
}

.checkbox-label input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.dialog-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.btn-secondary {
  padding: 10px 20px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary:hover {
  background: var(--bg-hover);
}
</style>

