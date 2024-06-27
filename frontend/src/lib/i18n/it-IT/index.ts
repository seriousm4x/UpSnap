import it from '../it';
import type { Translation } from '../i18n-types';

const it_IT = {
	...(it as unknown as Translation)
} satisfies Translation;

export default it_IT;
