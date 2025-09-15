import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
import type { SettingsPublic } from '$lib/types/settings';
import { get } from 'svelte/store';

export async function GET() {
	const manifest = {
		name: 'UpSnap',
		short_name: 'UpSnap',
		description: 'A simple wake on lan web app written with SvelteKit, Go and PocketBase.',
		theme_color: '#55BCD9',
		background_color: '#55BCD9',
		display: 'standalone',
		scope: '/',
		start_url: '/',
		icons: [{}]
	};
	const defaultIcons = [
		{
			src: '/icon_192.png',
			sizes: '192x192',
			type: 'image/png',
			purpose: 'any'
		},
		{
			src: '/icon_512.png',
			sizes: '512x512',
			type: 'image/png',
			purpose: 'any'
		},
		{
			src: '/maskable_192.png',
			sizes: '192x192',
			type: 'image/png',
			purpose: 'maskable'
		},
		{
			src: '/maskable_512.png',
			sizes: '512x512',
			type: 'image/png',
			purpose: 'maskable'
		},
		{
			src: '/gopher.svg',
			type: 'image/svg+xml',
			sizes: 'any'
		}
	];
	const headers = {
		'Content-Type': 'application/manifest+json',
		'Cache-Control': 'no-cache, no-store'
	};

	// get settings from backend
	const pb = get(pocketbase);
	const res = await pb
		.collection('settings_public')
		.getFirstListItem('')
		.catch(() => {
			return new Response(
				JSON.stringify({
					...manifest,
					icons: defaultIcons
				}),
				{
					headers: headers
				}
			);
		});
	const settings = res as SettingsPublic;

	// use custom website_title
	if (settings.website_title) {
		manifest.name = settings.website_title;
		manifest.short_name = settings.website_title;
	}

	// use custom icon
	if (settings.id && settings.favicon) {
		manifest.icons = [
			{
				src: `${backendUrl}api/files/settings_public/${settings.id}/${settings?.favicon}`,
				purpose: 'any'
			}
		];
	} else {
		manifest.icons = defaultIcons;
	}

	return new Response(JSON.stringify(manifest), {
		headers: headers
	});
}
