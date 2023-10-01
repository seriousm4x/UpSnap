import en_US from '../en-US';
import type { Translation } from '../i18n-types';

const de_DE = {
	...(en_US as unknown as Translation),
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
		change_password_confirm: 'Password bestätigen'
	},
	device: {
		page_title: 'Neues Gerät',
		tabs: ['Manuell', 'Netzwerkscan']
	},
	login: {
		welcome: 'Willkommen',
		email_label: 'Email oder Nutzername:',
		password_label: 'Passwort:',
		btn_more: 'Mehr',
		menu_title_auth_providers: 'Andere Login Provider',
		btn_login: 'Login'
	},
	settings: {
		settings_title: 'Einstellungen',
		ping_interval_title: 'Ping Intervall',
		ping_interval_desc1:
			'Legt den Intervall fest, in dem Geräte gepingt werden. Leer lassen um den Standardwert von <span class="badge">@every 3s</span> zu verwenden.',
		ping_interval_desc2:
			'Lerne mehr über die korrekte Cron Syntax auf <a class="link" href="https://de.wikipedia.org/wiki/Cron" target="_blank">Wikipedia</a> oder sieh dir die <a class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3" target="_blank">Paketdokumentation</a> an.',
		lazy_ping_title: 'Lazy Ping',
		lazy_ping_desc:
			'Wenn Lazy Ping aktiviert ist, pingt UpSnap die Geräte nur an, wenn ein aktiver Benutzer die Website besucht. Wenn es ausgeschaltet ist, pingt UpSnap die Geräte immer an.',
		lazy_ping_enable: 'Aktivieren',
		website_title_title: 'Seitentitel',
		website_title_desc: 'Setzt den Seiten- und den Tabtitel.',
		icon_title: 'Icon',
		icon_desc: 'Eigenes Favicon nutzen. Unterstützte Dateitypen sind:',
		upsnap_version: 'UpSnap Version'
	},
	users: {
		page_title: 'Benutzer',
		allow_create_devices:
			'{username} erlauben neue Geräte zu erstellen und Gerätegruppen zu bearbeiten',
		device_permissions: 'Geräteberechtigungen',
		create_new_device: 'Neues Gerät erstellen',
		read: 'Lesen',
		update: 'Ändern',
		delete: 'Löschen',
		power: 'Power',
		toggle: 'Toggle',
		confirm_delete_title: 'Löschen bestätigen',
		confirm_delete_desc: 'Soll {username} wirklich gelöscht werden?',
		create_new_user: 'Neuen Benutzer erstellen',
		username: 'Benutzername',
		password: 'Passwort',
		password_confirm: 'Passwort bestätigen',
		required_field: 'Pflichtfeld'
	},
	buttons: {
		save: 'Speichern',
		delete: 'Löschen',
		change: 'Ändern',
		reset: 'Zurücksetzen',
		cancle: 'Abbrechen',
		add: 'Hinzufügen'
	},
	welcome: {
		step1_page_title: 'Willkommen bei UpSnap 🥳',
		step1_setup_desc:
			'Bitte schließe die nachfolgenden Schritte ab, um die Einrichtung abzuschließen.',
		step1_setup_btn_next: 'Weiter',
		step2_page_title: 'Erstelle einen Admin Account',
		step2_label_email: 'Email:',
		step2_label_password: 'Passwort:',
		step2_label_min_chars: 'min. 10 Zeichen',
		step2_label_password_confirm: 'Passwort bestätigen:',
		step2_btn_create: 'Erstellen',
		step3_page_title: 'Du bist fertig! 🎉',
		step3_page_desc: 'Füge Geräte zu deinem Dashboard hinzu.',
		step3_btn_done: 'Lets go!',
		not_expected_title: 'Ich hab dich hier nicht erwartet! 🧐',
		not_expected_desc: 'Du hast die Einrichtung bereits abgeschlossen! Hier gibts nichts zu tun.',
		not_expected_back: 'Bring mich zurück',
		progress_step1: 'Willkommen',
		progress_step2: 'Account erstellen',
		progress_step3: 'Fertig'
	},
	toasts: {
		admin_saved: 'Admin gespeichert',
		user_saved: 'Benutzer gespeichert',
		user_created: 'Benutzer {username} erstellt',
		user_deleted: 'Benutzer {username} gelöscht',
		settings_saved: 'Einstellungen gespeichert',
		password_changed: 'Passwort geändert. Bitte neu einloggen.',
		passwords_missmatch: 'Passwörter stimmen nicht überein',
		permissions_created: 'Berechtigungen für {username} erstellt',
		permissions_deleted: 'Berechtigungen für {username} gelöscht',
		permissions_updated: 'Berechtigungen für {username} geändert',
		no_permission: 'Du hast keine Berechtigung für {url}'
	}
} satisfies Translation;

export default de_DE;
