import es from '../es';
import type { Translation } from '../i18n-types';

const es_ES = {
	...(es as unknown as Translation)
} satisfies Translation;

export default es_ES;
