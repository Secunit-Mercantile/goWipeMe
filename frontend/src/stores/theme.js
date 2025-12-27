import { writable } from 'svelte/store';

// Always use refined theme
export const theme = writable('refined');

export function setTheme(newTheme) {
  // Theme is locked to refined
  theme.set('refined');
}
