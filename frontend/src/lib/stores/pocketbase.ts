import { resolve } from '$app/paths';
import type { Permission } from '$lib/types/permission';
import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';

// set backend url based on environment
export const backendUrl = import.meta.env.DEV ? 'http://127.0.0.1:8090/' : resolve('/');

// connect to backend
const pb = new PocketBase(backendUrl);
pb.autoCancellation(false);

// export stores
export const pocketbase = writable(pb);
export const permission = writable({} as Permission);
