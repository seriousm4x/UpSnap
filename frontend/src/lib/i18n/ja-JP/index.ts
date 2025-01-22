import type { Translation } from '../i18n-types';
import ja from '../ja';

const ja_JP = {
	...(ja as unknown as Translation)
} satisfies Translation;

export default ja_JP;
