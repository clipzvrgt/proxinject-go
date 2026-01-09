<script setup lang="ts">
import { onMounted, onUnmounted, computed, watch } from 'vue';
import { useInjectorStore } from '../stores/injector';
import { useI18n } from '../i18n';

const store = useInjectorStore();
const { t } = useI18n();

let refreshInterval: number;

function startRefreshInterval() {
  clearInterval(refreshInterval);
  // Faster refresh when auto-inject is active
  const interval = store.autoInjectEnabled ? 1500 : 2500;
  refreshInterval = window.setInterval(() => {
    store.refreshProcesses();
  }, interval);
}

onMounted(() => {
  store.refreshProcesses();
  startRefreshInterval();
});

// Restart interval when auto-inject state changes
watch(() => store.autoInjectEnabled, () => {
  startRefreshInterval();
});

onUnmounted(() => {
  clearInterval(refreshInterval);
});

function toggleSelect(pid: number) {
  if (store.selectedPids.has(pid)) {
    store.selectedPids.delete(pid);
  } else {
    store.selectedPids.add(pid);
  }
}

function toggleSelectAll() {
  if (store.selectedPids.size === store.filteredProcesses.length) {
    store.selectedPids.clear();
  } else {
    store.filteredProcesses.forEach(p => store.selectedPids.add(p.pid));
  }
}

const canInject = computed(() => store.selectedPids.size > 0);

async function handleInject() {
  if (store.selectedProfile) {
    await store.setProxy(store.selectedProfile);
  }
  await store.injectSelected();
}
</script>

<template>
  <div class="process-page">
    <!-- Top Action Bar -->
    <div class="top-bar">
      <!-- Left: Inject Button -->
      <button 
        class="btn-primary inject-btn" 
        :disabled="!canInject"
        @click="handleInject"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M12 5v14M5 12h14"/>
        </svg>
        <span>{{ t('process.inject') }}</span>
        <span class="count-badge" v-if="store.selectedPids.size > 0">{{ store.selectedPids.size }}</span>
      </button>
      
      <!-- Center: Status Cards -->
      <div class="status-cards">
        <div class="stat-card">
          <div class="stat-icon injected">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
          </div>
          <div class="stat-content">
            <span class="stat-value">{{ store.injectedCount }}</span>
            <span class="stat-label">{{ t('process.injected') }}</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon proxy">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="2" y1="12" x2="22" y2="12"/>
              <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
            </svg>
          </div>
          <div class="stat-content">
            <span class="stat-value" v-if="store.selectedProfile">
              {{ store.selectedProfile.address }}
            </span>
            <span class="stat-value no-proxy" v-else>{{ t('process.noProxy') }}</span>
            <span class="stat-label">{{ t('process.proxy') }}</span>
          </div>
        </div>
      </div>
      
      <!-- Right: Proxy Selector -->
      <div class="proxy-selector">
        <select v-model="store.config.selectedProfileId">
          <option :value="null">{{ t('process.noProxy') }}</option>
          <option 
            v-for="profile in store.config.profiles" 
            :key="profile.id"
            :value="profile.id"
          >
            {{ profile.name }}
          </option>
        </select>
      </div>
    </div>
    
    <!-- Options Bar -->
    <div class="options-bar">
      <div class="toggle-pills">
        <label class="toggle-pill" :class="{ active: store.loggingEnabled }">
          <input 
            type="checkbox" 
            :checked="store.loggingEnabled"
            @change="store.setLogging(!store.loggingEnabled)"
          />
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
          </svg>
          <span>{{ t('process.logging') }}</span>
        </label>
        
        <label class="toggle-pill" :class="{ active: store.subprocessEnabled }">
          <input 
            type="checkbox" 
            :checked="store.subprocessEnabled"
            @change="store.setSubprocess(!store.subprocessEnabled)"
          />
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="16 18 22 12 16 6"/>
            <polyline points="8 6 2 12 8 18"/>
          </svg>
          <span>{{ t('process.subprocess') }}</span>
        </label>
        
        <label 
          class="toggle-pill" 
          :class="{ active: store.autoInjectEnabled, disabled: !store.hasInjectedOnce || store.config.watchlist.length === 0 }"
        >
          <input 
            type="checkbox" 
            :checked="store.autoInjectEnabled"
            :disabled="!store.hasInjectedOnce || store.config.watchlist.length === 0"
            @change="store.toggleAutoInject"
          />
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.21 15.89A10 10 0 1 1 8 2.83"/>
            <path d="M22 12A10 10 0 0 0 12 2v10z"/>
          </svg>
          <span>{{ t('process.autoInject') }}</span>
        </label>
        
        <label 
          class="toggle-pill" 
          :class="{ active: store.watchlistFilterEnabled, disabled: store.config.watchlist.length === 0 }"
        >
          <input 
            type="checkbox" 
            v-model="store.watchlistFilterEnabled"
            :disabled="store.config.watchlist.length === 0"
          />
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 3H2l8 9.46V19l4 2v-8.54L22 3z"/>
          </svg>
          <span>{{ t('process.watchlistFilter') }}</span>
        </label>
      </div>
      
      <!-- Search Box -->
      <div class="search-box">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input 
          type="text" 
          v-model="store.searchQuery"
          :placeholder="t('process.searchPlaceholder')" 
        />
      </div>
    </div>
    
    <!-- Process Table -->
    <div class="process-table">
      <table>
        <thead>
          <tr>
            <th style="width: 50px">
              <input 
                type="checkbox" 
                :checked="store.selectedPids.size === store.filteredProcesses.length && store.filteredProcesses.length > 0"
                @change="toggleSelectAll"
              />
            </th>
            <th>{{ t('process.processName') }}</th>
            <th style="width: 100px">{{ t('process.pid') }}</th>
            <th style="width: 120px">{{ t('process.status') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="process in store.filteredProcesses" 
            :key="process.pid"
            class="fade-in"
            :class="{ selected: store.selectedPids.has(process.pid) }"
          >
            <td>
              <input 
                type="checkbox" 
                :checked="store.selectedPids.has(process.pid)"
                @change="toggleSelect(process.pid)"
              />
            </td>
            <td class="process-name">
              <span class="name-text">{{ process.name }}</span>
            </td>
            <td class="pid-cell">{{ process.pid }}</td>
            <td>
              <span :class="['badge', process.injected ? 'badge-injected' : 'badge-ready']">
                {{ process.injected ? t('process.injected') : t('process.ready') }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="store.filteredProcesses.length === 0" class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="48" height="48">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <p>{{ t('process.noProcesses') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.process-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-primary);
}

/* ===================== */
/* Top Bar */
/* ===================== */

.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
  border-bottom: 1px solid var(--border);
  gap: 16px;
  flex-wrap: wrap;
}

@media (max-width: 900px) {
  .top-bar {
    justify-content: flex-start;
  }
}

.inject-btn {
  min-width: 140px;
  padding: 14px 24px;
  font-size: 15px;
  gap: 12px;
}

.inject-btn svg {
  width: 18px;
  height: 18px;
}

.count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 6px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 11px;
  font-size: 12px;
  font-weight: 700;
}

/* Status Cards */
.status-cards {
  display: flex;
  gap: 12px;
  flex: 1 1 auto;
  justify-content: center;
  flex-wrap: wrap;
  min-width: 0;
}

@media (max-width: 700px) {
  .status-cards {
    display: none;
  }
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: var(--bg-glass);
  backdrop-filter: var(--backdrop-blur);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  min-width: 120px;
  flex-shrink: 1;
}

.stat-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
}

.stat-icon svg {
  width: 20px;
  height: 20px;
}

.stat-icon.injected {
  background: var(--success-bg);
  color: var(--success);
}

.stat-icon.proxy {
  background: rgba(124, 58, 237, 0.15);
  color: var(--accent-secondary);
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-value.no-proxy {
  color: var(--text-muted);
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Proxy Selector */
.proxy-selector {
  flex-shrink: 0;
}

.proxy-selector select {
  padding: 10px 14px;
  padding-right: 36px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 13px;
  min-width: 140px;
  max-width: 200px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.proxy-selector select:hover {
  border-color: var(--border-light);
}

.proxy-selector select:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

/* ===================== */
/* Options Bar */
/* ===================== */

.options-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border);
  gap: 12px 20px;
  flex-wrap: wrap;
}

/* Toggle Pills */
.toggle-pills {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  flex: 1 1 auto;
  min-width: 0;
}

.toggle-pill {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.25s ease;
  user-select: none;
  white-space: nowrap;
}

.toggle-pill input {
  display: none;
}

.toggle-pill svg {
  width: 16px;
  height: 16px;
  opacity: 0.7;
  transition: all 0.25s ease;
}

.toggle-pill:hover {
  border-color: var(--border-light);
  color: var(--text-primary);
}

.toggle-pill:hover svg {
  opacity: 1;
}

.toggle-pill.active {
  background: rgba(124, 58, 237, 0.15);
  border-color: var(--accent-primary);
  color: var(--accent-secondary);
}

.toggle-pill.active svg {
  opacity: 1;
  color: var(--accent-secondary);
}

.toggle-pill.disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Search Box in Options */
.options-bar .search-box {
  flex: 1 1 200px;
  min-width: 150px;
  max-width: 320px;
}

/* ===================== */
/* Process Table */
/* ===================== */

.process-table {
  flex: 1;
  overflow: auto;
  padding: 0;
}

.process-table table {
  width: 100%;
}

.process-table thead {
  background: var(--bg-secondary);
}

.process-table th {
  padding: 16px 20px;
  vertical-align: middle;
}

.process-table td {
  padding: 16px 20px;
  vertical-align: middle;
}

.process-name {
  font-weight: 500;
}

.name-text {
  font-weight: 500;
}

.pid-cell {
  font-family: 'JetBrains Mono', 'Consolas', monospace;
  font-size: 13px;
  color: var(--text-secondary);
}

tbody tr.selected {
  background: rgba(124, 58, 237, 0.08);
}

tbody tr.selected td {
  border-bottom-color: rgba(124, 58, 237, 0.2);
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  gap: 16px;
  color: var(--text-muted);
}

.empty-state svg {
  opacity: 0.3;
}

.empty-state p {
  font-size: 15px;
}
</style>
