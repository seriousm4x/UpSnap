import en_US from '../en-US';
import type { Translation } from '../i18n-types';

const de_DE = {
	...(en_US as Translation),
	home: {
		order_groups: 'Gruppen',
		order_name: 'Name',
		order_ip: 'IP',
		order_tooltip: 'Sortieren',
		message_no_devices: 'Keine Geräte.',
		message_add_first_device: 'Füge dein erstes Gerät hinzu',
		message_grant_permissions:
			'Bitte frag den Administrator, dir Berechtigungen für bestehende Geräte zu erteilen oder neue Geräte zu erstellen.'
	},
	account: {
		page_title: 'Account',
		account_type_admin: 'Admin',
		account_type_user: 'Benutzer',
		avatar_title: 'Avatar',
		language_title: 'Sprache',
		language_option_auto: 'Automatisch',
		change_password_title: 'Passwort ändern',
		change_password_body: 'Nachdem das Passwort geändert wurde, musst du dich neu einloggen.',
		change_password_label: 'Altes Passwort',
		change_password_new: 'Neues Password',
		change_password_confirm: 'Password bestätigen',
		toast_admin_saved: 'Admin gespeichert',
		toast_user_saved: 'Benutzer gespeichert',
		toast_password_changed: 'Passwort geändert. Bitte neu einloggen.',
		toast_passwords_missmatch: 'Passwörter stimmen nicht überein'
	},
	buttons: {
		save: 'Speichern',
		change: 'Ändern'
	}
} satisfies Translation;

export default de_DE;
