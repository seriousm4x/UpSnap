import pt from '../pt';
import type { Translation } from '../i18n-types';

const pt_PT = {
	...(pt as unknown as Translation)
} satisfies Translation;

export default pt_PT;
