export type SettingsPublic = {
	collectionId: string;
	favicon: string;
	setup_completed: boolean;
	website_title: string;
};

export type SettingsPrivate = {
	interval: number;
	scan_range: string;
};
