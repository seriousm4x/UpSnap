import type { RecordModel } from 'pocketbase';

export type User = RecordModel & {
	username: string;
	email: string;
	avatar: number;
};
