import { writable } from 'svelte/store';
import type { SettingsPublic, SettingsPrivate } from '$lib/types/settings';

export const settingsPub = writable<SettingsPublic>();
export const settingsPriv = writable<SettingsPrivate>();
