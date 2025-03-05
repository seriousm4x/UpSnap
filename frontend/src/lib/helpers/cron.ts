import { m } from '$lib/paraglide/messages';
import { dateFnsLocale } from '$lib/stores/locale';
import cronParser from 'cron-parser';
import { formatDistanceStrict, type Locale } from 'date-fns';
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
			locale: get(dateFnsLocale) as unknown as Locale,
			addSuffix: true
		});
	} catch {
		return m.settings_invalid_cron();
	}
}
