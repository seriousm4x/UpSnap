import { writable } from 'svelte/store';

export const theme = writable();
if (typeof window !== 'undefined')
    [theme.set(localStorage.getItem('theme') ? localStorage.getItem('theme') : 'auto')];
