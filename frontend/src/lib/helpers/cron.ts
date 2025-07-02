import { m } from '$lib/paraglide/messages';
import { dateFnsLocale } from '$lib/stores/locale';
import { backendUrl } from '$lib/stores/pocketbase';
import cronParser from 'cron-parser';
import { formatDate, type Locale } from 'date-fns';
import { get } from 'svelte/store';

export function nextCronDate(expression: string) {
	try {
		const cron = cronParser.parse(expression);
		const cronStr = cron.next().toISOString();
		if (!cronStr) {
			throw new Error(m.settings_invalid_cron());
		}
		return formatDate(cronStr, 'PPpp', {
			locale: get(dateFnsLocale) as unknown as Locale
		});
	} catch {
		return m.settings_invalid_cron();
	}
}

export function parseCron(expression: string): Promise<boolean> {
	return fetch(backendUrl + 'api/upsnap/validate-cron', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			cron: expression
		})
	})
		.then((response) => {
			return response.ok;
		})
		.catch(() => {
			return false;
		});
}
