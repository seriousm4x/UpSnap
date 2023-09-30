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
		change_password_confirm: 'Confirm password',
		toast_admin_saved: 'Admin saved',
		toast_user_saved: 'User saved',
		toast_password_changed: 'Password changed. Please login again.',
		toast_passwords_missmatch: "Passwords don't match"
	},
	buttons: {
		save: 'Save',
		change: 'Change'
	}
} satisfies BaseTranslation;

export default en_US;
