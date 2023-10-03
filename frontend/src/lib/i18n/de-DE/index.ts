import de from '../de';
import type { Translation } from '../i18n-types';

const de_DE = {
	...(de as unknown as Translation)
} satisfies Translation;

export default de_DE;
