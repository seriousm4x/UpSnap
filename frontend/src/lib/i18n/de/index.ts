import en from '../en';
import type { Translation } from '../i18n-types';

const de = {
	...(en as unknown as Translation),
	home: {
		page_title: 'Home',
		order_groups: 'Gruppen',
		order_name: 'Name',
		order_ip: 'IP',
		order_tooltip: 'Sortieren',
		no_devices: 'Keine Geräte.',
		add_first_device: 'Füge dein erstes Gerät hinzu',
		grant_permissions:
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
		tabs: ['Manuell', 'Netzwerkscan'],
		card_btn_more: 'Mehr',
		card_btn_more_edit: 'Bearbeiten',
		card_btn_more_sleep: 'Ruhezustand',
		card_btn_more_reboot: 'Neustart',
		card_tooltip_wake_cron: 'Einschalten Cron',
		card_tooltip_shutdown_cron: 'Ausschalten Cron',
		card_tooltip_wake_password: 'Einschalten Passwort',
		card_tooltip_last_status_change: 'Letzte Statusänderung',
		card_password: 'Passwort',
		card_nic_tooltip_pending: 'Warten',
		card_nic_tooltip_shutdown: 'Ausschalten',
		card_nic_tooltip_shutdown_no_cmd: 'Kein Ausschaltbefehl gesetzt',
		card_nic_tooltip_shutdown_no_permission: 'Keine Berechtigung zum ausschalten',
		card_nic_tooltip_power: 'Einschalten',
		card_nic_tooltip_power_no_permission: 'Keine Berechtigung zum einschalten',
		general: 'Allgemein',
		general_name: 'Name',
		general_name_placeholder: 'Gerätename',
		general_ip: 'IP',
		general_mac: 'Mac',
		general_netmask: 'Netmask',
		general_required_field: 'Pflichtfeld',
		ports: 'Ports',
		ports_desc: 'UpSnap kann auch prüfen, ob bestimmte Ports offen sind.',
		ports_add_new: 'Port hinzufügen',
		ports_name: 'Name',
		ports_number: 'Nummer',
		link: 'Link',
		link_desc:
			'Macht Ihren Gerätenamen zu einem anklickbaren Link, ideal zum Beispiel für die Verknüpfung eines Dashboards.',
		wake: 'Einschalten',
		wake_desc: 'Du kannst das Gerät mit einem Cron-Job einschalten.',
		wake_cron: 'Cron',
		wake_cron_enable: 'Aktivieren',
		sol: 'Sleep-On-LAN',
		sol_desc1:
			'Du kannst Computer mithilfe des Tools <a class="link" href="https://github.com/SR-G/sleep-on-lan" target="_blank">Sleep-On-LAN</a> in den Ruhezustand versetzen. Sleep-On-LAN (SOL) ist ein externes Tool/Daemon, das auf den PCs arbeitet, die du in den Ruhezustand versetzen möchtest, und stellt einen REST-Endpunkt bereit. Für Anweisungen zur Einrichtung von Sleep-On-LAN verweise bitte auf den Abschnitt <a href="https://github.com/SR-G/sleep-on-lan#usage" class="link" target="_blank">Usage</a>.',
		sol_desc2:
			'SOL ist so konfiguriert, dass es Anfragen über HTTP statt über UDP sendet, um eine Autorisierung zu ermöglichen und die Zuverlässigkeit der Anfragen zu erhöhen.',
		sol_desc3:
			'Daher stellen Sie bitte sicher, dass Sie <span class="badge">HTTP:&lt;DEINPORT&gt;</span> im Abschnitt <span class="badge">Listener</span> der <a href="https://github.com/SR-G/sleep-on-lan#configuration" class="link" target="_blank">SOL-Konfiguration</a> hinzufügen.',
		sol_enable: 'Sleep-On-LAN aktivieren',
		sol_port: 'SOL Port',
		sol_authorization: 'Authorisierung',
		sol_user: 'SOL Benutzer',
		sol_password: 'SOL Passwort',
		shutdown: 'Ausschalten',
		shutdown_desc:
			'Dieser <strong>Shell-Befehl</strong> wird in deinem Container ausgeführt (wenn du Docker verwendest) oder auf deinem Host (wenn du die Binärdatei verwendest). Um zu überprüfen, ob das funktioniert, kannst du den Befehl zuerst im Container oder auf deiner Host-Shell ausführen. Übliche Befehle sind <span class="badge">net rpc</span> für Windows, <span class="badge">sshpass</span> für Linux oder <span class="badge">curl</span> im Allgemeinen, um Webanfragen durchzuführen.',
		shutdown_examples: 'Beispiele:',
		shutdown_examples_windows: 'Entfernten Windows PC ausschalten:',
		shutdown_examples_linux: 'Entfernten Linux PC ausschalten:',
		shutdown_cmd: 'Befehl zum ausschalten',
		shutdown_cron_desc:
			'Genau wie das Einrichten eines Cron-Jobs, um das Gerät aufzuwecken, kannst du auch einen Cron-Job planen, um dieses Gerät herunterzufahren.',
		shutdown_cron: 'Cron',
		shutdown_cron_enable: 'Aktivieren',
		password: 'Password',
		password_desc:
			'Einige Netzwerkkarten bieten die Möglichkeit, ein Passwort für Magic Packets zu setzen, auch <span class="badge">SecureON</span> genannt. Das Passwort kann nur 0, 4 oder 6 Zeichen lang sein.',
		groups: 'Gruppen',
		groups_desc:
			'Du kannst Geräte einer Gruppe hinzufügen, um sie auf dem Dashboard nach Gruppen sortiert anzuzeigen.',
		groups_placeholder: "z.B. 'Keller' or 'Büro'",
		network_scan_range_saved: 'Scan-Bereich gespeichert',
		network_scan_desc:
			'Automatisches Scannen deines Netzwerks nach Geräten. Damit dies funktioniert, musst du UpSnap als Administrator/root ausführen und nmap installiert und im $PATH verfügbar haben (Für Docker-Benutzer ist das bereits der Fall, und du musst nichts weiter tun). Das Scannen kann einige Sekunden dauern.',
		network_scan_ip_range: 'IP-Bereich',
		network_scan_no_range: 'Kein Scan-Bereich',
		network_scan_unsaved_changes: 'Ungespeicherte Änderungen',
		network_scan_running: 'Scan läuft',
		network_scan: 'Scan',
		network_scan_ip: 'IP:',
		network_scan_mac: 'Mac:',
		network_scan_mac_vendor: 'Mac Hersteller:',
		network_scan_netmask: 'Netmask:',
		network_scan_add_all: 'Alle Geräte hinzufügen',
		network_scan_replace_netmask: 'Netmask für alle Geräte ändern?',
		network_scan_new_netmask: 'Neue Netmask',
		network_scan_include_unknown: 'Schließe Geräte ein, bei denen der Name "Unknown" ist.'
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
		page_title: 'Einstellungen',
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
		cancel: 'Abbrechen',
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
		permissions_updated_personal: 'Deine Berechtigungen wurden geupdated',
		no_permission: 'Du hast keine Berechtigung für {url}',
		device_created: '{device} erstellt',
		devices_created_multiple: '{count|int} Geräte erstellt',
		device_updated: '{device} geändert',
		device_deleted: '{device} gelöscht',
		group_created: 'Gruppe {group} erstellt',
		group_deleted: 'Gruppe {group} gelöscht'
	},
	navbar: {
		theme: 'Theme',
		new: 'Neu',
		edit_account: 'Account bearbeiten',
		logout: 'Abmelden'
	}
} satisfies Translation;

export default de;
