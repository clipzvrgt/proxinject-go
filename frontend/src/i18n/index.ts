import { ref, computed } from 'vue';
import en from './en';
import zh from './zh';

export type Locale = 'en' | 'zh';

type TranslationValue = string | Record<string, unknown>;

const messages: Record<Locale, typeof en> = { en, zh };

// Current locale, persisted to localStorage
const currentLocale = ref<Locale>(loadLocale());

// Sync initial locale to backend
syncLocaleToBackend(currentLocale.value);

function loadLocale(): Locale {
  const saved = localStorage.getItem('proxinject-locale');
  if (saved === 'en' || saved === 'zh') {
    return saved;
  }
  // Auto-detect from browser
  const browserLang = navigator.language.toLowerCase();
  return browserLang.startsWith('zh') ? 'zh' : 'en';
}

function saveLocale(locale: Locale) {
  localStorage.setItem('proxinject-locale', locale);
}

// Sync locale to Go backend for tray menu
function syncLocaleToBackend(locale: Locale) {
  try {
    // Call Go backend SetLocale
    if (window.go?.main?.App?.SetLocale) {
      window.go.main.App.SetLocale(locale);
    }
  } catch (e) {
    console.error('Failed to sync locale to backend:', e);
  }
}

export function useI18n() {
  const locale = computed({
    get: () => currentLocale.value,
    set: (val: Locale) => {
      currentLocale.value = val;
      saveLocale(val);
      syncLocaleToBackend(val);
    },
  });
  
  function t(key: string): string {
    const keys = key.split('.');
    let value: TranslationValue = messages[currentLocale.value];
    
    for (const k of keys) {
      if (value && typeof value === 'object' && k in value) {
        value = (value as Record<string, TranslationValue>)[k];
      } else {
        // Fallback to English
        value = messages.en;
        for (const fallbackKey of keys) {
          if (value && typeof value === 'object' && fallbackKey in value) {
            value = (value as Record<string, TranslationValue>)[fallbackKey];
          } else {
            return key; // Return key if not found
          }
        }
        break;
      }
    }
    
    return typeof value === 'string' ? value : key;
  }
  
  return {
    locale,
    t,
    availableLocales: ['en', 'zh'] as const,
  };
}

