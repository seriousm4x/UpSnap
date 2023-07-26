import type { Record } from 'pocketbase';

export type Device = Record & {
	name: string;
	ip: string;
	mac: string;
	netmask: string;
	status: string;
	ports: string[];
	link: URL;
	wake_cron: string;
	wake_cron_enabled: boolean;
	shutdown_cron: string;
	shutdown_cron_enabled: boolean;
	shutdown_cmd: string;
	password: string;
	groups: Group[];
	expand: {
		ports: Port[];
	};
};

export type Port = Record & {
	name: string;
	number: number;
};

export type Group = Record & {
	name: string;
};
