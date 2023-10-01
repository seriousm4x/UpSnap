import type { BaseTranslation } from '../i18n-types';

const en_US = {
	home: {
		order_groups: 'Groups',
		order_name: 'Name',
		order_ip: 'IP',
		order_tooltip: 'Order',
		message_no_devices: 'No devices here.',
		message_add_first_device: 'Add your first device',
		message_grant_permissions:
			'Please ask the admin to grant you permissions to existing devices or to create new ones.'
	},
	account: {
		page_title: 'Account',
		account_type_admin: 'Admin',
		account_type_user: 'User',
		avatar_title: 'Avatar',
		language_title: 'Language',
		language_option_auto: 'Automatic',
		change_password_title: 'Change password',
		change_password_body: 'After the password has been changed, you will need to log in again.',
		change_password_label: 'Old password',
		change_password_new: 'New password',
		change_password_confirm: 'Confirm password'
	},
	device: {
		page_title: 'New device',
		tabs: ['Manual', 'Network Scan']
	},
	login: {
		welcome: 'Welcome',
		email_label: 'Email or Username:',
		password_label: 'Password:',
		btn_more: 'More',
		menu_title_auth_providers: 'Other Auth Providers',
		btn_login: 'Login'
	},
	settings: {
		settings_title: 'Settings',
		ping_interval_title: 'Ping interval',
		ping_interval_desc1:
			'Sets the interval in which the devices are pinged. Leave blank to use default value of <span class="badge">@every 3s</span>.',
		ping_interval_desc2:
			'Learn more about the correct syntax for cron on <a class="link" href="https://en.wikipedia.org/wiki/Cron" target="_blank">Wikipedia</a> or refer to the <a class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3" target="_blank">package documentation</a>.',
		lazy_ping_title: 'Lazy ping',
		lazy_ping_desc:
			"When lazy ping is turned on, UpSnap will only ping devices if there is an active user visiting the website. If it's turned off, UpSnap will always ping devices.",
		lazy_ping_enable: 'Enable',
		website_title_title: 'Website title',
		website_title_desc: 'Sets the title of the website and in the browser tab.',
		icon_title: 'Icon',
		icon_desc: 'Set a custom favicon. Supported file types are:',
		upsnap_version: 'UpSnap version'
	},
	users: {
		page_title: 'Users',
		allow_create_devices: 'Allow {username} to create new devices and edit device groups',
		device_permissions: 'Device permissions',
		create_new_device: 'Create new device',
		read: 'Read',
		update: 'Update',
		delete: 'Delete',
		power: 'Power',
		toggle: 'Toggle',
		confirm_delete_title: 'Confirm delete',
		confirm_delete_desc: 'Are you sure you want to delete {username}?',
		create_new_user: 'Create new user',
		username: 'Username',
		password: 'Password',
		password_confirm: 'Password confirm',
		required_field: 'required field'
	},
	buttons: {
		save: 'Save',
		delete: 'Delete',
		change: 'Change',
		reset: 'Reset',
		cancle: 'Cancle',
		add: 'Add'
	},
	welcome: {
		step1_page_title: 'Welcome to UpSnap 🥳',
		step1_setup_desc: 'Please complete the following steps to finish the setup.',
		step1_setup_btn_next: 'Next',
		step2_page_title: 'Create an admin account',
		step2_label_email: 'Email:',
		step2_label_password: 'Password:',
		step2_label_min_chars: 'min. 10 characters',
		step2_label_password_confirm: 'Password confirm:',
		step2_btn_create: 'Create',
		step3_page_title: 'You are all set! 🎉',
		step3_page_desc: 'Go ahead and add some devices to your dashboard.',
		step3_btn_done: 'Lets go!',
		not_expected_title: "I didn't expect you here! 🧐",
		not_expected_desc: 'You are already done with the setup! Nothing to do.',
		not_expected_back: 'Take me back',
		progress_step1: 'Welcome',
		progress_step2: 'Create account',
		progress_step3: 'Done'
	},
	toasts: {
		admin_saved: 'Admin saved',
		user_saved: 'User saved',
		user_created: 'User {username} created',
		user_deleted: 'User {username} deleted',
		settings_saved: 'Saved settings',
		password_changed: 'Password changed. Please login again.',
		passwords_missmatch: "Passwords don't match",
		permissions_created: 'Permissions for {username} created',
		permissions_deleted: 'Permissions for {username} deleted',
		permissions_updated: 'Permissions for {username} updated',
		no_permission: "You don't have permission to visit {url}"
	}
} satisfies BaseTranslation;

export default en_US;
