import zh from '../zh';
import type { Translation } from '../i18n-types';

const zh_CN = {
	...(zh as unknown as Translation)
} satisfies Translation;

export default zh_CN;
