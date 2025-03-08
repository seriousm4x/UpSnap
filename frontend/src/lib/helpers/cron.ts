import { m } from '$lib/paraglide/messages';
import { dateFnsLocale } from '$lib/stores/locale';
import { backendUrl } from '$lib/stores/pocketbase';
import cronParser from 'cron-parser';
import { formatDate, type Locale } from 'date-fns';
import { get } from 'svelte/store';

export const cronRegex = new RegExp(/((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,6})/);

export function nextCronDate(expression: string) {
	try {
		if (
			expression === '' ||
			!cronRegex.test(expression) ||
			Object.keys(cronParser.parseString(expression).errors).length > 0
		) {
			throw new Error('cron-parser failed to parse string');
		}

		const cron = cronParser.parseExpression(expression, {});

		return formatDate(cron.next().toISOString(), 'PPpp', {
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
