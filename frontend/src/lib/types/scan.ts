export type ScannedDevice = {
	name: string;
	ip: string;
	mac: string;
	mac_vendor: string;
	netmask: string;
};

export type ScanResponse = {
	netmask: string;
	devices: ScannedDevice[];
};
