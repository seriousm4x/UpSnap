import en from '../en';
import type { Translation } from '../i18n-types';

const nl = {
	...(en as unknown as Translation),
	home: {
		page_title: 'Home',
		order_groups: 'Groepen',
		order_name: 'Naam',
		order_ip: 'IP',
		order_tooltip: 'Sorteren',
		no_devices: 'Geen apparaten gevonden.',
		add_first_device: 'Voeg je eerste apparaat toe',
		grant_permissions:
			'Vraag de admin om je toestemming te geven voor bestaande apparaten of om nieuwe te creëren.'
	},
	account: {
		page_title: 'Account',
		account_type_admin: 'Beheerder',
		account_type_user: 'Gebruiker',
		avatar_title: 'Avatar',
		language_title: 'Taal',
		language_option_auto: 'Automatisch',
		change_password_title: 'Wachtwoord wijzigen',
		change_password_body: 'Na het wijzigen van het wachtwoord moet je opnieuw inloggen.',
		change_password_label: 'Oud wachtwoord',
		change_password_new: 'Nieuw wachtwoord',
		change_password_confirm: 'Bevestig wachtwoord'
	},
	device: {
		page_title: 'Nieuw apparaat',
		tabs: ['Handmatig', 'Netwerkscan'],
		card_btn_more: 'Meer',
		card_btn_more_edit: 'Bewerken',
		card_btn_more_sleep: 'Slaapstand',
		card_btn_more_reboot: 'Herstarten',
		card_tooltip_wake_cron: 'Wek cron',
		card_tooltip_shutdown_cron: 'Afsluiten cron',
		card_tooltip_wake_password: 'Wek wachtwoord',
		card_tooltip_last_status_change: 'Laatste statuswijziging',
		card_password: 'Wachtwoord',
		card_nic_tooltip_pending: 'In behandeling',
		card_nic_tooltip_shutdown: 'Afsluiten',
		card_nic_tooltip_shutdown_no_cmd: 'Geen uitschakel-opdracht ingesteld',
		card_nic_tooltip_shutdown_no_permission: 'Geen toestemming om dit apparaat uit te schakelen',
		card_nic_tooltip_power: 'Aanzetten',
		card_nic_tooltip_power_no_permission: 'Geen toestemming om dit apparaat aan te zetten',
		modal_confirm_wake_title: 'Wek {device}?',
		modal_confirm_wake_desc: 'Bevestig om {device} te wekken.',
		modal_confirm_shutdown_title: 'Afsluiten {device}?',
		modal_confirm_shutdown_desc: 'Bevestig om {device} af te sluiten.',
		general: 'Algemeen',
		general_name: 'Naam',
		general_name_placeholder: 'Apparaatnaam',
		general_ip: 'IP',
		general_mac: 'Mac',
		general_netmask: 'Subnetmask',
		general_required_field: 'vereist veld',
		ports: 'Poorten',
		ports_desc: 'UpSnap kan ook controleren of gegeven poorten open zijn.',
		ports_add_new: 'Voeg nieuwe poort toe',
		ports_name: 'Naam',
		ports_number: 'Nummer',
		link: 'Link',
		link_desc:
			'Maakt de apparaatsnaam een klikbare link, ideaal om bijvoorbeeld een dashboard te linken.',
		ping: 'Ping',
		ping_desc:
			'Je kunt een aangepaste shell-opdracht gebruiken om te zien of het apparaat aan staat. De opdracht moet een exitcode van <span class="badge">0</span> retourneren om aan te geven dat het apparaat aan is, elke andere exitcode markeert het apparaat als uitgeschakeld.',
		ping_cmd: 'Aangepaste ping-opdracht',
		wake: 'Wek',
		wake_desc: 'Je kunt dit apparaat inschakelen met een geplande cron-taak.',
		wake_cmd: 'Aangepaste wek-opdracht',
		wake_cron: 'Wek cron',
		wake_cron_enable: 'Wek cron inschakelen',
		sol: 'Sleep-On-LAN',
		sol_desc1:
			'Je kunt computers in de slaapstand zetten met de <a class="link" href="https://github.com/SR-G/sleep-on-lan" target="_blank">Sleep-On-LAN</a> tool. Sleep-On-LAN (SOL) is een externe tool/daemon die op de PC\'s werkt die je in de slaapstand wilt zetten en biedt een REST-endpoint. Voor instructies over het instellen van Sleep-On-LAN, zie de <a href="https://github.com/SR-G/sleep-on-lan#usage" class="link" target="_blank">Gebruik</a> sectie.',
		sol_desc2:
			'SOL is geconfigureerd om verzoeken over HTTP in plaats van UDP te verzenden om autorisatie mogelijk te maken en verzoeken betrouwbaarder te maken.',
		sol_desc3:
			'Zorg er daarom voor dat je <span class="badge">HTTP:&lt;JOUWPOORT&gt;</span> opneemt in de <span class="badge">Listeners</span> sectie van de <a href="https://github.com/SR-G/sleep-on-lan#configuration" class="link" target="_blank">SOL-configuratie</a>.',
		sol_enable: 'Sleep-On-LAN inschakelen',
		sol_port: 'SOL Poort',
		sol_authorization: 'Autorisatie',
		sol_user: 'SOL Gebruiker',
		sol_password: 'SOL Wachtwoord',
		shutdown: 'Afsluiten',
		shutdown_desc:
			'Deze <strong>shell-opdracht</strong> wordt uitgevoerd binnen de container (als je docker gebruikt) of op jouw host (als je de binaire versie gebruikt). Om te verifiëren dat het werkt, kun je de opdracht eerst binnen de container of op jouw host-shell uitvoeren. Veelvoorkomende opdrachten zijn <span class="badge">net rpc</span> voor Windows, <span class="badge">sshpass</span> voor Linux of <span class="badge">curl</span> in het algemeen om webverzoeken te doen.',
		shutdown_examples: 'Voorbeelden:',
		shutdown_examples_windows: 'Afsluiten van een externe Windows-machine:',
		shutdown_examples_linux: 'Afsluiten van een externe Linux-machine:',
		shutdown_cmd: 'Afsluitopdracht',
		shutdown_cron_desc:
			'Net als het instellen van een cron om het apparaat te wekken, kun je ook een cron-taak plannen om dit apparaat af te sluiten.',
		shutdown_cron: 'Afsluitcron',
		shutdown_cron_enable: 'Afsluitcron inschakelen',
		password: 'Wachtwoord',
		password_desc:
			'Sommige netwerkkaarten hebben de optie om een wachtwoord voor magic packets in te stellen, ook wel <span class="badge">SecureON</span> genoemd. Het wachtwoord kan alleen 0, 4 of 6 tekens lang zijn.',
		groups: 'Groepen',
		groups_desc:
			'Je kunt apparaten aan een groep toevoegen om ze op het dashboard per groep te sorteren.',
		groups_placeholder: "bijv. 'Kelder' of 'Kantoor'",
		network_scan_range_saved: 'Scanbereik opgeslagen',
		network_scan_desc:
			'Automatisch jouw netwerk scannen op apparaten. Voor dit te laten werken, moet je UpSnap als root/admin uitvoeren en moet nmap geïnstalleerd en beschikbaar zijn in jouw $PATH (Voor docker-gebruikers is dit al het geval en hoef je niets meer te doen). Scannen kan enkele seconden duren.',
		network_scan_ip_range: 'IP-bereik',
		network_scan_no_range: 'Geen scanbereik',
		network_scan_unsaved_changes: 'Niet-opgeslagen wijzigingen',
		network_scan_running: 'Scan loopt',
		network_scan: 'Scannen',
		network_scan_ip: 'IP:',
		network_scan_mac: 'Mac:',
		network_scan_mac_vendor: 'Mac leverancier:',
		network_scan_netmask: 'Subnetmasker:',
		network_scan_add_all: 'Alle apparaten toevoegen',
		network_scan_replace_netmask: 'Subnetmasker voor alle apparaten vervangen?',
		network_scan_new_netmask: 'Nieuw subnetmasker',
		network_scan_include_unknown: 'Voeg apparaten toe waarvan de naam "Onbekend" is',
		require_confirmation: 'Bevestiging vereist'
	},
	login: {
		welcome: 'Welkom',
		email_label: 'E-mail of gebruikersnaam:',
		password_label: 'Wachtwoord:',
		btn_more: 'Meer',
		menu_title_auth_providers: 'Andere auth-providers',
		btn_login: 'Inloggen'
	},
	settings: {
		page_title: 'Instellingen',
		ping_interval_title: 'Ping-interval',
		ping_interval_desc1:
			'Stelt de interval in waarin de apparaten worden gepingt. Laat leeg om de standaardwaarde van <span class="badge">@every 3s</span> te gebruiken.',
		ping_interval_desc2:
			'Leer meer over de juiste syntax voor cron op <a class="link" href="https://en.wikipedia.org/wiki/Cron" target="_blank">Wikipedia</a> of raadpleeg de <a class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3" target="_blank">pakketdocumentatie</a>.',
		lazy_ping_title: 'Lazy ping',
		lazy_ping_desc:
			'Wanneer lazy ping is ingeschakeld, pingt UpSnap alleen apparaten als er een actieve gebruiker de website bezoekt. Als het is uitgeschakeld, pingt UpSnap altijd apparaten.',
		lazy_ping_enable: 'Inschakelen',
		website_title_title: 'Website titel',
		website_title_desc: 'Stelt de titel van de website en in het browsertabblad in.',
		icon_title: 'Icoon',
		icon_desc: 'Stel een aangepast favicon in. Ondersteunde bestandstypen zijn:',
		upsnap_version: 'UpSnap versie'
	},
	users: {
		page_title: 'Gebruikers',
		allow_create_devices:
			'Sta {username} toe om nieuwe apparaten te maken en apparaatgroepen te bewerken',
		device_permissions: 'Apparaattoestemmingen',
		create_new_device: 'Nieuw apparaat maken',
		read: 'Lezen',
		update: 'Bijwerken',
		delete: 'Verwijderen',
		power: 'Aanzetten',
		toggle: 'Wisselen',
		confirm_delete_title: 'Bevestig verwijderen',
		confirm_delete_desc: 'Weet je zeker dat je {username} wilt verwijderen?',
		create_new_user: 'Nieuwe gebruiker aanmaken',
		username: 'Gebruikersnaam',
		password: 'Wachtwoord',
		password_confirm: 'Bevestig wachtwoord',
		required_field: 'vereist veld'
	},
	buttons: {
		save: 'Opslaan',
		delete: 'Verwijderen',
		change: 'Wijzigen',
		reset: 'Resetten',
		cancel: 'Annuleren',
		add: 'Toevoegen',
		confirm: 'Bevestigen'
	},
	welcome: {
		step1_page_title: 'Welkom bij UpSnap 🥳',
		step1_setup_desc: 'Voltooi de volgende stappen om de installatie te voltooien.',
		step1_setup_btn_next: 'Volgende',
		step2_page_title: 'Maak een beheerdersaccount aan',
		step2_label_email: 'E-mail:',
		step2_label_password: 'Wachtwoord:',
		step2_label_min_chars: 'min. 10 tekens',
		step2_label_password_confirm: 'Bevestig wachtwoord:',
		step2_btn_create: 'Aanmaken',
		step3_page_title: 'Je bent helemaal klaar! 🎉',
		step3_page_desc: 'Ga je gang en voeg enkele apparaten toe aan je dashboard.',
		step3_btn_done: 'Laten we beginnen!',
		not_expected_title: 'Ik zag je hier niet aankomen! 🧐',
		not_expected_desc: 'Je bent al klaar met de installatie! Er is niets meer te doen.',
		not_expected_back: 'Breng me terug',
		progress_step1: 'Welkom',
		progress_step2: 'Account aanmaken',
		progress_step3: 'Klaar'
	},
	toasts: {
		admin_saved: 'Beheerder opgeslagen',
		user_saved: 'Gebruiker opgeslagen',
		user_created: 'Gebruiker {username} aangemaakt',
		user_deleted: 'Gebruiker {username} verwijderd',
		settings_saved: 'Instellingen opgeslagen',
		password_changed: 'Wachtwoord gewijzigd. Log alsjeblieft opnieuw in.',
		passwords_missmatch: 'Wachtwoorden komen niet overeen',
		permissions_created: 'Rechten voor {username} aangemaakt',
		permissions_deleted: 'Rechten voor {username} verwijderd',
		permissions_updated: 'Rechten voor {username} bijgewerkt',
		permissions_updated_personal: 'Jouw rechten zijn bijgewerkt',
		no_permission: 'Je hebt geen toestemming om {url} te bezoeken',
		device_created: 'Aangemaakt {device}',
		devices_created_multiple: 'Aangemaakt {count|int} apparaten',
		device_updated: 'Bijgewerkt {device}',
		device_deleted: 'Verwijderd {device}',
		group_created: 'Groep {group} aangemaakt',
		group_deleted: 'Groep {group} verwijderd'
	},
	navbar: {
		theme: 'Thema',
		new: 'Nieuw',
		edit_account: 'Account bewerken',
		logout: 'Uitloggen'
	}
} satisfies Translation;

export default nl;