<script setup lang="ts">
import { useI18n } from '../i18n';

const { t } = useI18n();

defineProps<{
  current: 'main' | 'settings' | 'logs' | 'help';
}>();

const emit = defineEmits<{
  'update:current': [value: 'main' | 'settings' | 'logs' | 'help'];
}>();
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-logo">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
      </svg>
      <span>Proxinject</span>
    </div>
    
    <nav>
      <div 
        class="nav-item" 
        :class="{ active: current === 'main' }"
        @click="emit('update:current', 'main')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="7" height="7"/>
          <rect x="14" y="3" width="7" height="7"/>
          <rect x="14" y="14" width="7" height="7"/>
          <rect x="3" y="14" width="7" height="7"/>
        </svg>
        <span>{{ t('nav.processes') }}</span>
      </div>
      
      <div 
        class="nav-item"
        :class="{ active: current === 'settings' }"
        @click="emit('update:current', 'settings')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M12 1v4m0 14v4m11-11h-4M5 12H1m16.95-6.95l-2.83 2.83M8.88 15.12l-2.83 2.83m11.9 0l-2.83-2.83M8.88 8.88L6.05 6.05"/>
        </svg>
        <span>{{ t('nav.settings') }}</span>
      </div>
      
      <div 
        class="nav-item"
        :class="{ active: current === 'logs' }"
        @click="emit('update:current', 'logs')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
        </svg>
        <span>{{ t('nav.logs') }}</span>
      </div>
      
      <div 
        class="nav-item"
        :class="{ active: current === 'help' }"
        @click="emit('update:current', 'help')"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
          <line x1="12" y1="17" x2="12.01" y2="17"/>
        </svg>
        <span>{{ t('nav.help') }}</span>
      </div>
    </nav>
  </aside>
</template>
