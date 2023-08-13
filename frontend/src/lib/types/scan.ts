export type ScannedDevice = {
	name: string;
	ip: string;
	mac: string;
	mac_vendor: string;
	netmask: string;
	status: 'pending' | 'online' | 'offline' | '';
};

export type ScanResponse = {
	netmask: string;
	devices: ScannedDevice[];
};
