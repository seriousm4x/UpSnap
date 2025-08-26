import type { RecordModel } from 'pocketbase';

export type Device = RecordModel & {
	name: string;
	ip: string;
	mac: string;
	netmask: string;
	description: string;
	status: 'pending' | 'online' | 'offline' | '';
	ports: string[];
	link: string;
	link_open: '' | 'same_tab' | 'new_tab';
	ping_cmd: string;
	wake_cron: string;
	wake_cron_enabled: boolean;
	wake_cmd: string;
	wake_confirm: boolean;
	wake_timeout: number;
	shutdown_cron: string;
	shutdown_cron_enabled: boolean;
	shutdown_cmd: string;
	shutdown_confirm: boolean;
	shutdown_timeout: number;
	password: string;
	groups: string[];
	expand: {
		ports: Port[];
		groups: Group[];
	};
	created_by: string;
	sol_enabled: boolean;
	sol_auth: boolean;
	sol_user: string;
	sol_password: string;
	sol_port: number;
};

export type Port = RecordModel & {
	name: string;
	number: number;
	link: string;
};

export type Group = RecordModel & {
	name: string;
};
