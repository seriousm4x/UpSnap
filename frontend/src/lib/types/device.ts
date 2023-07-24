import type { Record } from 'pocketbase';

export type Device = Record & {
	expand: {
		ports: Record[];
	};
};

export type Port = Record & {
	name: string;
	number: number;
};
