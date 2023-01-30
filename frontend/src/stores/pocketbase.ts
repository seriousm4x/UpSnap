import { writable } from 'svelte/store';
import PocketBase from 'pocketbase';

// set backend url based on environment
let backend_url = '';
let isDevMode = import.meta.env.DEV;
if (isDevMode) {
    backend_url = 'http://127.0.0.1:8090';
} else {
    backend_url = '/';
}

// get default values for stores
const pb = new PocketBase(backend_url);
pb.autoCancellation(false);
const settingsRes = await pb.collection('settings').getList(1, 1);
let tempDevices = {};
const devicesRes = await pb.collection('devices').getFullList(200, {
    sort: 'name',
    expand: 'ports'
});
devicesRes.forEach((device) => {
    tempDevices[device.id] = device;
});

// export stores
export const pocketbase = writable(pb);
export const settings = writable(settingsRes.items[0]);
export const devices = writable(tempDevices);
