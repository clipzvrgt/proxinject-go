<script setup lang="ts">
import { useInjectorStore } from '../stores/injector';

const store = useInjectorStore();
</script>

<template>
  <footer class="status-bar">
    <div class="status-indicator">
      <span :class="['status-dot', store.injectedCount > 0 ? 'online' : 'offline']"></span>
      <span>Connected: {{ store.injectedCount }} processes</span>
    </div>
    
    <div class="status-right">
      <span v-if="store.selectedProfile">
        Proxy: {{ store.selectedProfile.name }} ({{ store.selectedProfile.address }}:{{ store.selectedProfile.port }})
      </span>
      <span v-else>No proxy configured</span>
      
      <span v-if="store.autoInjectEnabled" class="auto-badge">
        🔄 Auto-Inject Active
      </span>
    </div>
  </footer>
</template>

<style scoped>
.status-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.auto-badge {
  padding: 4px 10px;
  background: rgba(99, 102, 241, 0.2);
  color: var(--accent-secondary);
  border-radius: 12px;
  font-size: 12px;
}
</style>
