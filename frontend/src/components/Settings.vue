<script setup lang="ts">
import { ref } from 'vue';
import { useInjectorStore, type ProxyProfile } from '../stores/injector';
import { useI18n, type Locale } from '../i18n';

const store = useInjectorStore();
const { t, locale, availableLocales } = useI18n();

// Profile form
const showProfileForm = ref(false);
const editingProfile = ref<ProxyProfile | null>(null);
const profileForm = ref({
  name: '',
  address: '',
  port: 1080,
  username: '',
  password: '',
});

// Watchlist form
const newExe = ref('');

const localeLabels: Record<Locale, string> = {
  en: 'English',
  zh: '中文',
};

function openProfileForm(profile?: ProxyProfile) {
  if (profile) {
    editingProfile.value = profile;
    profileForm.value = { ...profile, username: profile.username || '', password: profile.password || '' };
  } else {
    editingProfile.value = null;
    profileForm.value = { name: '', address: '', port: 1080, username: '', password: '' };
  }
  showProfileForm.value = true;
}

function saveProfile() {
  if (!profileForm.value.name || !profileForm.value.address) return;
  
  const profile = {
    name: profileForm.value.name,
    address: profileForm.value.address,
    port: profileForm.value.port,
    username: profileForm.value.username || undefined,
    password: profileForm.value.password || undefined,
  };
  
  if (editingProfile.value) {
    store.updateProfile(editingProfile.value.id, profile);
  } else {
    store.addProfile(profile);
  }
  
  showProfileForm.value = false;
}

function addExe() {
  if (newExe.value.trim()) {
    store.addToWatchlist(newExe.value.trim());
    newExe.value = '';
  }
}
</script>

<template>
  <div class="settings-page">
    <header class="header">
      <h2>{{ t('settings.title') }}</h2>
    </header>
    
    <div class="settings-content">
      <!-- Language -->
      <section class="settings-section">
        <div class="section-header">
          <h3>{{ t('settings.language') }}</h3>
        </div>
        <div class="input-group">
          <select v-model="locale">
            <option v-for="loc in availableLocales" :key="loc" :value="loc">
              {{ localeLabels[loc] }}
            </option>
          </select>
        </div>
      </section>
      
      <!-- Proxy Profiles -->
      <section class="settings-section">
        <div class="section-header">
          <h3>{{ t('settings.proxyProfiles') }}</h3>
          <button class="btn-add" @click="openProfileForm()">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            {{ t('settings.addNew') }}
          </button>
        </div>
        
        <div class="profile-list">
          <div 
            v-for="profile in store.config.profiles" 
            :key="profile.id"
            class="profile-item"
          >
            <div class="profile-info">
              <span class="profile-name">📡 {{ profile.name }}</span>
              <span class="profile-address">
                {{ profile.address }}:{{ profile.port }}
                <template v-if="profile.username"> • user: {{ profile.username }}</template>
              </span>
            </div>
            <div class="profile-actions">
              <button class="btn-icon" @click="openProfileForm(profile)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="btn-icon danger" @click="store.deleteProfile(profile.id)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>
          </div>
          
          <div v-if="store.config.profiles.length === 0" class="empty-state">
            {{ t('settings.noProfiles') }}
          </div>
        </div>
      </section>
      
      <!-- Watchlist -->
      <section class="settings-section">
        <div class="section-header">
          <h3>{{ t('settings.watchlist') }}</h3>
        </div>
        <p class="section-desc">
          {{ t('settings.watchlistDesc') }}
        </p>
        
        <div class="watchlist">
          <div 
            v-for="exe in store.config.watchlist" 
            :key="exe"
            class="watchlist-item"
          >
            <span>{{ exe }}</span>
            <button class="btn-icon danger" @click="store.removeFromWatchlist(exe)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          
          <div class="add-exe-row">
            <input 
              v-model="newExe"
              type="text" 
              :placeholder="t('settings.exePlaceholder')"
              @keyup.enter="addExe"
            />
            <button class="btn-add" @click="addExe">{{ t('settings.add') }}</button>
          </div>
        </div>
      </section>
    </div>
    
    <!-- Profile Dialog -->
    <div v-if="showProfileForm" class="dialog-overlay">
      <div class="dialog">
        <h3>{{ editingProfile ? t('settings.editProfile') : t('settings.newProfile') }}</h3>
        
        <div class="input-group">
          <label>{{ t('settings.profileName') }}</label>
          <input v-model="profileForm.name" type="text" placeholder="My Proxy" />
        </div>
        
        <div class="input-row">
          <div class="input-group" style="flex: 2">
            <label>{{ t('settings.address') }}</label>
            <input v-model="profileForm.address" type="text" placeholder="127.0.0.1" />
          </div>
          <div class="input-group" style="flex: 1">
            <label>{{ t('settings.port') }}</label>
            <input v-model.number="profileForm.port" type="number" />
          </div>
        </div>
        
        <div class="input-row">
          <div class="input-group">
            <label>{{ t('settings.username') }}</label>
            <input v-model="profileForm.username" type="text" placeholder="username" />
          </div>
          <div class="input-group">
            <label>{{ t('settings.password') }}</label>
            <input v-model="profileForm.password" type="password" placeholder="password" />
          </div>
        </div>
        
        <div class="dialog-actions">
          <button class="btn-secondary" @click="showProfileForm = false">{{ t('settings.cancel') }}</button>
          <button class="btn-primary" @click="saveProfile">{{ t('settings.save') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.settings-content {
  flex: 1;
  overflow: auto;
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h3 {
  font-size: 16px;
}

.section-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.profile-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.empty-state {
  color: var(--text-muted);
  text-align: center;
  padding: 24px;
}

.add-exe-row {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.add-exe-row input {
  flex: 1;
  padding: 10px 14px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
}

/* Dialog */
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

.dialog h3 {
  margin-bottom: 20px;
}

.input-row {
  display: flex;
  gap: 12px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.btn-secondary {
  padding: 10px 20px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
}

.btn-secondary:hover {
  background: var(--bg-hover);
}
</style>
