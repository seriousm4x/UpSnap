import { getLocale } from '$lib/paraglide/runtime';
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

export const localeStore = writable(getLocale());
export const dateFnsLocale: Writable<Locale> = writable(enUS);

localeStore.subscribe((l: string) => {
	switch (l) {
		case 'de':
			dateFnsLocale.set(de);
			break;
		case 'en':
			dateFnsLocale.set(enUS);
			break;
		case 'es':
			dateFnsLocale.set(es);
			break;
		case 'fr':
			dateFnsLocale.set(fr);
			break;
		case 'id':
			dateFnsLocale.set(id);
			break;
		case 'it':
			dateFnsLocale.set(it);
			break;
		case 'ja':
			dateFnsLocale.set(ja);
			break;
		case 'nl':
			dateFnsLocale.set(nl);
			break;
		case 'pl':
			dateFnsLocale.set(pl);
			break;
		case 'pt':
			dateFnsLocale.set(pt);
			break;
		case 'zh':
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
