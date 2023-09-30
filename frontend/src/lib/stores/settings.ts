import type { SettingsPrivate, SettingsPublic } from '$lib/types/settings';
import { writable } from 'svelte/store';

export const settingsPub = writable<SettingsPublic>();
export const settingsPriv = writable<SettingsPrivate>();
