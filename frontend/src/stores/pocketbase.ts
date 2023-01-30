import { writable } from 'svelte/store';
import PocketBase from 'pocketbase';

let backend_url = '';
let isDevMode = import.meta.env.DEV;
if (isDevMode) {
    backend_url = 'http://localhost:8090';
} else {
    backend_url = '/';
}
const pb = new PocketBase(backend_url);
pb.autoCancellation(false);

export const pocketbase = writable(pb);
