import nl from '../nl';
import type { Translation } from '../i18n-types';

const nl_NL = {
	...(nl as unknown as Translation)
} satisfies Translation;

export default nl_NL;
