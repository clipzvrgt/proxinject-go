<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useInjectorStore } from '../stores/injector';
import { useI18n } from '../i18n';

const store = useInjectorStore();
const { t } = useI18n();
const logContainer = ref<HTMLElement | null>(null);
const autoScroll = ref(true);

function scrollToBottom() {
  if (autoScroll.value && logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight;
  }
}

onMounted(() => {
  scrollToBottom();
});

function clearLogs() {
  store.logs = [];
}
</script>

<template>
  <div class="logs-page">
    <header class="header">
      <h2>{{ t('logs.title') }}</h2>
      <div class="header-actions">
        <label class="auto-scroll">
          <input type="checkbox" v-model="autoScroll" /> {{ t('logs.autoScroll') }}
        </label>
        <button class="btn-add" @click="clearLogs">{{ t('logs.clear') }}</button>
      </div>
    </header>
    
    <div ref="logContainer" class="log-container">
      <div v-if="store.logs.length === 0" class="empty-logs">
        <p>{{ t('logs.noLogs') }}</p>
        <p>{{ t('logs.enableLogging') }}</p>
      </div>
      
      <div 
        v-for="(log, i) in store.logs" 
        :key="i" 
        class="log-entry"
      >
        {{ log }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.logs-page {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.auto-scroll {
  font-size: 13px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.empty-logs {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-muted);
  text-align: center;
}

.empty-logs p {
  margin: 4px 0;
}
</style>
