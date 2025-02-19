import id from '../id';
import type { Translation } from '../i18n-types';

const id_ID = {
	...(id as unknown as Translation)
} satisfies Translation;

export default id_ID;
