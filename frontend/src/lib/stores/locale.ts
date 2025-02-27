import { locale } from '$lib/i18n/i18n-svelte';
import type { Locale } from 'date-fns';
import { id } from 'date-fns/locale';
import { de } from 'date-fns/locale/de';
import { enUS } from 'date-fns/locale/en-US';
import { es } from 'date-fns/locale/es';
import { fr } from 'date-fns/locale/fr';
import { it } from 'date-fns/locale/it';
import { ja } from 'date-fns/locale/ja';
import { nl } from 'date-fns/locale/nl';
import { pl } from 'date-fns/locale/pl';
import { pt } from 'date-fns/locale/pt';
import { zhCN } from 'date-fns/locale/zh-CN';
import { zhTW } from 'date-fns/locale/zh-TW';

import { writable, type Writable } from 'svelte/store';

// dynamically import locales
// it's ugly but the only working solution i found in ~ 3 hours
export const dateFnsLocale: Writable<Locale> = writable(enUS);

locale.subscribe((l) => {
	switch (l) {
		case 'de':
		case 'de-DE':
			dateFnsLocale.set(de);
			break;
		case 'en':
		case 'en-US':
			dateFnsLocale.set(enUS);
			break;
		case 'es':
		case 'es-ES':
			dateFnsLocale.set(es);
			break;
		case 'fr':
		case 'fr-FR':
			dateFnsLocale.set(fr);
			break;
		case 'id':
		case 'id-ID':
			dateFnsLocale.set(id);
			break;
		case 'it':
		case 'it-IT':
			dateFnsLocale.set(it);
			break;
		case 'ja':
		case 'ja-JP':
			dateFnsLocale.set(ja);
			break;
		case 'nl':
		case 'nl-NL':
			dateFnsLocale.set(nl);
			break;
		case 'pl':
		case 'pl-PL':
			dateFnsLocale.set(pl);
			break;
		case 'pt':
		case 'pt-PT':
			dateFnsLocale.set(pt);
			break;
		case 'zh':
		case 'zh-CN':
			dateFnsLocale.set(zhCN);
			break;
		case 'zh-TW':
			dateFnsLocale.set(zhTW);
			break;
		default:
			dateFnsLocale.set(enUS);
			break;
	}
});
