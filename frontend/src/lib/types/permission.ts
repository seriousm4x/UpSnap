import type { Record } from 'pocketbase';

export type Permission = Record & {
	user: string;
	create: boolean;
	read: string[];
	update: string[];
	delete: string[];
	power: string[];
};
