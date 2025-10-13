import { getLocale } from '$lib/paraglide/runtime';
import type { Locale } from 'date-fns';
import { bg } from 'date-fns/locale/bg';
import { cs } from 'date-fns/locale/cs';
import { de } from 'date-fns/locale/de';
import { enUS } from 'date-fns/locale/en-US';
import { es } from 'date-fns/locale/es';
import { fr } from 'date-fns/locale/fr';
import { hi } from 'date-fns/locale/hi';
import { id } from 'date-fns/locale/id';
import { it } from 'date-fns/locale/it';
import { ja } from 'date-fns/locale/ja';
import { ko } from 'date-fns/locale/ko';
import { nb } from 'date-fns/locale/nb';
import { nl } from 'date-fns/locale/nl';
import { pl } from 'date-fns/locale/pl';
import { pt } from 'date-fns/locale/pt';
import { ru } from 'date-fns/locale/ru';
import { uk } from 'date-fns/locale/uk';
import { vi } from 'date-fns/locale/vi';
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
		case 'cs-CZ':
			dateFnsLocale.set(cs);
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
		case 'hi-IN':
			dateFnsLocale.set(hi);
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
		case 'ru-RU':
			dateFnsLocale.set(ru);
			break;
		case 'uk-UA':
			dateFnsLocale.set(uk);
			break;
		case 'vi-VN':
			dateFnsLocale.set(vi);
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
