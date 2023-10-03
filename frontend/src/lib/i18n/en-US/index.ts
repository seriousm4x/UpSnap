import en from '../en';
import type { Translation } from '../i18n-types';

const en_US = {
	...(en as unknown as Translation)
} satisfies Translation;

export default en_US;
