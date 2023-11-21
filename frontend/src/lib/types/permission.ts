import type { RecordModel } from 'pocketbase';

export type Permission = RecordModel & {
	user: string;
	create: boolean;
	read: string[];
	update: string[];
	delete: string[];
	power: string[];
};
