import pl from '../pl';
import type { Translation } from '../i18n-types';

const pl_PL = {
	...(pl as unknown as Translation)
} satisfies Translation;

export default pl_PL;
