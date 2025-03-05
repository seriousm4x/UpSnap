import { paraglideMiddleware } from '$lib/paraglide/server';
import { localeStore } from '$lib/stores/locale';
import type { Handle } from '@sveltejs/kit';

// creating a handle to use the paraglide middleware
const paraglideHandle: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ locale }) => {
		localeStore.set(locale);
		return resolve(event, {
			transformPageChunk: ({ html }) => html.replace('%lang%', locale)
		});
	});

export const handle: Handle = paraglideHandle;
