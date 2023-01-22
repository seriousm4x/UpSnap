import { writable } from 'svelte/store';
import PocketBase from 'pocketbase';

export let pocketbase = writable(new PocketBase('http://127.0.0.1:8090'));
