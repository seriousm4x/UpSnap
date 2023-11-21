import type { RecordModel } from 'pocketbase';

export type SettingsPublic = RecordModel & {
	collectionId: string;
	favicon: string;
	setup_completed: boolean;
	website_title: string;
};

export type SettingsPrivate = RecordModel & {
	interval: number;
	lazy_ping: boolean;
	scan_range: string;
};
