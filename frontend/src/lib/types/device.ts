import type { Record } from 'pocketbase';

export type Device = Record & {
	expand: {
		ports: Record[];
	};
};
