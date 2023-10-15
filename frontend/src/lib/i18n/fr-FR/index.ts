import fr from '../fr';
import type { Translation } from '../i18n-types';

const fr_FR = {
	...(fr as unknown as Translation)
} satisfies Translation;

export default fr_FR;
