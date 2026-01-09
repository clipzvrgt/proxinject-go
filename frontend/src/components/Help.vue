<script setup lang="ts">
import { useI18n } from '../i18n';

const { t } = useI18n();
</script>

<template>
  <div class="help-page">
    <header class="header">
      <h2>{{ t('help.title') }}</h2>
    </header>
    
    <div class="help-content">
      <p class="help-subtitle">{{ t('help.subtitle') }}</p>
      
      <div class="steps-container">
        <div class="step-item">
          <div class="step-number">1</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step1') }}</span>
          </div>
        </div>
        
        <div class="step-item">
          <div class="step-number">2</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step2') }}</span>
          </div>
        </div>
        
        <div class="step-item">
          <div class="step-number">3</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step3') }}</span>
          </div>
        </div>
        
        <div class="step-item">
          <div class="step-number">4</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step4') }}</span>
          </div>
        </div>
        
        <div class="step-item">
          <div class="step-number">5</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step5') }}</span>
          </div>
        </div>
        
        <div class="step-item">
          <div class="step-number">6</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step6') }}</span>
          </div>
        </div>
        
        <div class="step-item">
          <div class="step-number">7</div>
          <div class="step-content">
            <span class="step-text">{{ t('help.steps.step7') }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.help-page {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.help-content {
  flex: 1;
  overflow: auto;
  padding: 24px;
}

.help-subtitle {
  font-size: 15px;
  color: var(--text-secondary);
  margin-bottom: 32px;
}

.steps-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 600px;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  transition: all 0.2s ease;
}

.step-item:hover {
  border-color: var(--accent-primary);
  background: var(--bg-tertiary);
}

.step-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
  border-radius: 50%;
  color: white;
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.step-content {
  flex: 1;
  display: flex;
  align-items: center;
  min-height: 32px;
}

.step-text {
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.5;
}
</style>
