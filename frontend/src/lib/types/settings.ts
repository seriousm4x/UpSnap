import type { Record } from 'pocketbase';

export type SettingsPublic = Record & {
	collectionId: string;
	favicon: string;
	setup_completed: boolean;
	website_title: string;
};

export type SettingsPrivate = Record & {
	interval: number;
	lazy_ping: boolean;
	scan_range: string;
};
