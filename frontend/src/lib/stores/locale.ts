import { getLocale } from '$lib/paraglide/runtime';
import type { Locale } from 'date-fns';
import { id } from 'date-fns/locale';
import { bg } from 'date-fns/locale/bg';
import { de } from 'date-fns/locale/de';
import { enUS } from 'date-fns/locale/en-US';
import { es } from 'date-fns/locale/es';
import { fr } from 'date-fns/locale/fr';
import { it } from 'date-fns/locale/it';
import { ja } from 'date-fns/locale/ja';
import { ko } from 'date-fns/locale/ko';
import { nl } from 'date-fns/locale/nl';
import { nb } from 'date-fns/locale/nb';
import { pl } from 'date-fns/locale/pl';
import { pt } from 'date-fns/locale/pt';
import { zhCN } from 'date-fns/locale/zh-CN';
import { zhTW } from 'date-fns/locale/zh-TW';

import { writable, type Writable } from 'svelte/store';

export const localeStore = writable(getLocale());
export const dateFnsLocale: Writable<Locale> = writable(enUS);

localeStore.subscribe((l: string) => {
	switch (l) {
		case 'bg-BG':
			dateFnsLocale.set(bg);
			break;
		case 'de-DE':
			dateFnsLocale.set(de);
			break;
		case 'en-US':
			dateFnsLocale.set(enUS);
			break;
		case 'es-ES':
			dateFnsLocale.set(es);
			break;
		case 'fr-FR':
			dateFnsLocale.set(fr);
			break;
		case 'id-ID':
			dateFnsLocale.set(id);
			break;
		case 'it-IT':
			dateFnsLocale.set(it);
			break;
		case 'ja-JP':
			dateFnsLocale.set(ja);
			break;
		case 'ko-KR':
			dateFnsLocale.set(ko);
			break;
		case 'nl-NL':
			dateFnsLocale.set(nl);
			break;
		case 'nb-NO':
			dateFnsLocale.set(nb);
			break;
		case 'pl-PL':
			dateFnsLocale.set(pl);
			break;
		case 'pt':
			dateFnsLocale.set(pt);
			break;
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
