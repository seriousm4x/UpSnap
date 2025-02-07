import LL from '$lib/i18n/i18n-svelte';
import { dateFnsLocale } from '$lib/stores/locale';
import cronParser from 'cron-parser';
import { formatDistanceStrict } from 'date-fns';
import { get } from 'svelte/store';

export const cronRegex = new RegExp(/((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){6})/);

export function parseCron(expression: string, now: number) {
	try {
		if (
			expression === '' ||
			!cronRegex.test(expression) ||
			Object.keys(cronParser.parseString(expression).errors).length > 0
		) {
			throw new Error('cron-parser failed to parse string');
		}

		const cron = cronParser.parseExpression(expression, {});

		return formatDistanceStrict(cron.next().toISOString(), new Date(now), {
			locale: get(dateFnsLocale),
			addSuffix: true
		});
	} catch {
		return get(LL).settings.invalid_cron();
	}
}
