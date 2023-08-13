import type { Record } from 'pocketbase';

export type User = Record & {
	username: string;
	email: string;
	avatar: number;
};
