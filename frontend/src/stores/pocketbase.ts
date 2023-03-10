import { writable } from 'svelte/store';
import PocketBase from 'pocketbase';

// set backend url based on environment
let backend_url = import.meta.env.DEV ? 'http://127.0.0.1:8090' : '/';

// get default values for stores
const pb = new PocketBase(backend_url);
pb.autoCancellation(false);

// export stores
export const pocketbase = writable(pb);
export const settings_private = writable({
    interal: '',
    scan_range: ''
});
export const settings_public = writable({
    website_title: '',
    favicon: ''
});
export const devices = writable({});
export const authorizedStore = writable(false);
