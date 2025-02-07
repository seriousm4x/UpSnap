import en from '../en';
import type { Translation } from '../i18n-types';

const pl = {
	...(en as unknown as Translation),
	home: {
		page_title: 'Strona Główna',
		order_groups: 'Grupowanie',
		order_name: 'Nazwa',
		order_ip: 'IP',
		order_tooltip: 'Kolejność',
		no_devices: 'Brak urządzeń.',
		add_first_device: 'Utwórz pierwsze urządzenie',
		grant_permissions:
			'Poproś administratora o przyznanie uprawnień do istniejących urządzeń lub do tworzenia nowych.',
		search_placeholder: 'Wyszukiwanie urządzeń'
	},
	account: {
		page_title: 'Konto',
		account_type_admin: 'Administrator',
		account_type_user: 'Użytkownik',
		avatar_title: 'Awatar',
		language_title: 'Język',
		language_option_auto: 'Automatyczny',
		change_password_title: 'Zmiana hasła',
		change_password_body: 'Po zmianie hasła, będzie konieczne ponowne zalogowanie.',
		change_password_label: 'Stare hasło',
		change_password_new: 'Nowe hasło',
		change_password_confirm: 'Potwierdź hasło'
	},
	device: {
		page_title: 'Nowe urządzenie',
		tabs: ['Ręcznie', 'Skanuj sieć'],
		card_btn_more: 'Więcej',
		card_btn_more_edit: 'Edytuj',
		card_btn_more_sleep: 'Wyłącz',
		card_btn_more_reboot: 'Zrestartuj',
		card_tooltip_wake_cron: 'Cron do włączania',
		card_tooltip_shutdown_cron: 'Cron do wyłączania',
		card_tooltip_wake_password: 'Hasło do włączania',
		card_tooltip_last_status_change: 'Ostatnia zmiana statusu',
		card_password: 'Hasło',
		card_nic_tooltip_pending: 'Oczekiwanie',
		card_nic_tooltip_shutdown: 'Wyłącz',
		card_nic_tooltip_shutdown_no_cmd: 'Brak ustawionej polecenia do wyłączania',
		card_nic_tooltip_shutdown_no_permission: 'Brak uprawnień do wyłączania tego urządzenia',
		card_nic_tooltip_power: 'Włącz',
		card_nic_tooltip_power_no_permission: 'Brak uprawnień do włączania tego urządzenia',
		modal_confirm_wake_title: 'Włączyć {device}?',
		modal_confirm_wake_desc: 'Potwierdź włączenie {device}.',
		modal_confirm_shutdown_title: 'Wyłączyć {device}?',
		modal_confirm_shutdown_desc: 'Potwierdź wyłączenie {device}.',
		general: 'Ogólne',
		general_name: 'Nazwa',
		general_ip: 'IP',
		general_mac: 'MAC',
		general_netmask: 'Maska sieci',
		general_description: 'Opis',
		general_description_placeholder: 'Przykładowy opis urządzenia',
		general_required_field: 'pole wymagane',
		ports: 'Porty',
		ports_desc: 'UpSnap może sprawdzić czy podane porty są otwarte.',
		ports_add_new: 'Dodaj nowy port',
		ports_name: 'Nazwa',
		ports_number: 'Numer',
		link: 'Link',
		link_desc: 'Twoje urządzenie może być klikalnym łączem. Świetne do linkowania np. pulpitów.',
		link_open: 'Otwórz link automatycznie',
		link_open_no: 'Nie',
		link_open_same_tab: 'W tej samej karcie',
		link_open_new_tab: 'W nowej karcie',
		ping: 'Ping',
		ping_desc:
			'Możesz wykorzystać niestandardowe polecenie powłoki, aby sprawdzić, czy urządzenie jest włączone. Polecenie powinno zwrócić kod zakończenia <span class="badge">0</span> aby wskazać, że urządzenie jest włączone. Dowolny inny kod zakończenia poskutkuje oznaczeniem urządzenia jako wyłączone.',
		ping_cmd: 'Niestandardowe polecenie ping',
		wake: 'Włączanie',
		wake_desc: 'Możesz włączyć to urządzenie za pomocą zaplanowanego zadania cron.',
		wake_cmd: 'Niestandardowe polecenie do włączania',
		wake_cron: 'Cron do włączania',
		wake_cron_enable: 'Włącz cron do włączania',
		sol: 'Sleep-On-LAN',
		sol_desc1:
			'Możesz wyłączać urządzenia za pomocą narzędzia <a class="link" href="https://github.com/SR-G/sleep-on-lan" target="_blank">Sleep-On-LAN</a>. Sleep-On-LAN (SOL) jest zewnętrznym narzędziem/demonem, działającym na urządzeniach które chcesz wyłączać. Tworzy on endpoint REST. Aby dowiedzieć się jak skonfigurować SOL, zapoznaj się z <a href="https://github.com/SR-G/sleep-on-lan#usage" class="link" target="_blank">instrukcją użytkownika</a>.',
		sol_desc2:
			'SOL przyjmuje żądania za pomocą HTTP zamiast UDP, aby zwiększyć niezawodność i umożliwić autoryzację.',
		sol_desc3:
			'W związku z tym, uwzględnij <span class="badge">HTTP:&lt;TWÓJPORT&gt;</span> w sekcji <span class="badge">Listeners</span> twojej <a href="https://github.com/SR-G/sleep-on-lan#configuration" class="link" target="_blank">konfiguracji SOL</a>.',
		sol_enable: 'Włącz Sleep-On-LAN',
		sol_port: 'Port SOL',
		sol_authorization: 'Autoryzacja',
		sol_user: 'Użytkownik SOL',
		sol_password: 'Hasło SOL',
		shutdown: 'Wyłącz',
		shutdown_desc:
			'To <strong>polecenie powłoki</strong> zostanie wywołane w Twoim kontenerze (jeżeli używasz Dockera) lub na Twoim hoście (jeżeli używasz pliku binarnego). Aby zweryfikować jego działanie, możesz je najpierw uruchomić wewnątrz kontenera lub na hoście. Powszechne polecenia to <span class="badge">net rpc</span> dla Windowsa, <span class="badge">sshpass</span> dla Linuxa lub <span class="badge">curl</span> do wykonywania żądań internetowych.',
		shutdown_examples: 'Przykłady:',
		shutdown_examples_windows: 'Wyłącz zdalną maszynę Windows:',
		shutdown_examples_linux: 'Wyłącz zdalną maszynę Linux:',
		shutdown_cmd: 'Polecenie do wyłączania',
		shutdown_cron_desc:
			'Podobnie jak możesz harmonogramować zadanie cron do włączania urządzenia, możesz robić to również do jego wyłączania.',
		shutdown_cron: 'Cron do wyłączania',
		shutdown_cron_enable: 'Włącz cron do wyłączania',
		password: 'Hasło',
		password_desc:
			'Niektóre karty sieciowe mają możliwość ustawienia hasła dla magicznych pakietów, znanych również jako <span class="badge">SecureON</span>. Hasło może mieć 0, 4 lub 6 znaków długości.',
		groups: 'Grupy',
		groups_desc: 'Możesz dodać urządzenie do grup, aby sortować je według grup na pulpicie',
		groups_placeholder: "np. 'Piwnica' lub 'Biuro'",
		network_scan_range_saved: 'Zapisano zakres skanowania',
		network_scan_desc:
			'Automatycznie skanuj sieć w poszukiwaniu urządzeń. Aby to zadziałało, musisz uruchomić UpSnap jako root/admin i mieć nmap zainstalowany oraz dostępny w $PATH (preinstalowany w obrazach dockerowych). Skanowanie może chwilę potrwać.',
		network_scan_ip_range: 'Zakres IP',
		network_scan_no_range: 'Brak zakresu skanowania',
		network_scan_unsaved_changes: 'Niezapisane zmiany',
		network_scan_running: 'Skanowanie w toku',
		network_scan: 'Skanuj',
		network_scan_ip: 'IP:',
		network_scan_mac: 'MAC:',
		network_scan_mac_vendor: 'Dostawca MAC:',
		network_scan_netmask: 'Maska sieci:',
		network_scan_add_all: 'Dodaj wszystkie urządzenia',
		network_scan_replace_netmask: 'Zastąpić maskę sieci dla wszystkich urządzeń?',
		network_scan_new_netmask: 'Nowa maska sieci',
		network_scan_include_unknown: 'Uwzględnij urządzenia, których nazwa to "Unknown"',
		require_confirmation: 'Wymagaj potwierdzenia'
	},
	login: {
		welcome: 'Witaj',
		email_label: 'Email lub nazwa użytkownika:',
		password_label: 'Hasło:',
		btn_more: 'Więcej',
		menu_title_auth_providers: 'Inni dostawcy autoryzacji',
		btn_login: 'Zaloguj'
	},
	settings: {
		page_title: 'Ustawienia',
		ping_interval_title: 'Interwał pingowania',
		ping_interval_desc1:
			'Ustawia interwał, w którym urządzenia są pingowane. Pozostaw puste aby użyć domyślnej wartości <span class="badge">@every 3s</span>.',
		ping_interval_desc2:
			'Dowiedz się więcej na temat składni cron na <a class="link" href="https://pl.wikipedia.org/wiki/Cron" target="_blank">Wikipedii</a> lub w <a class="link" href="https://github.com/harrisiirak/cron-parser" target="_blank">dokumentacji pakietu cron</a>.',
		lazy_ping_title: 'Leniwe pingowanie',
		lazy_ping_desc:
			'Gdy leniwe pingowanie jest włączone, UpSnap będzie pingować urządzenia tylko wtedy, gdy użytkownicy odwiedzają stronę. Gdy ta opcja jest wyłączona, UpSnap będzie pingować urządzenia nieustannie.',
		lazy_ping_enable: 'Włącz',
		website_title_title: 'Tytuł strony',
		website_title_desc: 'Ustawia tytuł strony internetowej i na karcie przeglądarki.',
		icon_title: 'Ikona',
		icon_desc: 'Ustaw niestandardowy favicon. Obsługiwane są typy plików:',
		upsnap_version: 'Wersja UpSnap',
		invalid_cron: 'Nieprawidłowa składnia cron'
	},
	users: {
		page_title: 'Użytkownicy',
		allow_create_devices: 'Pozwól {username} tworzyć nowe urządzenia i edytować grupy urządzeń',
		device_permissions: 'Uprawnienia do urządzeń',
		create_new_device: 'Utwórz nowe urządzenie',
		read: 'Odczyt',
		update: 'Aktualizacja',
		delete: 'Usuwanie',
		power: 'Włączanie',
		toggle: 'Przełącz',
		confirm_delete_title: 'Potwierdź usunięcie',
		confirm_delete_desc: 'Czy na pewno chcesz usunąć {username}?',
		create_new_user: 'Utwórz użytkownika',
		username: 'Nazwa użytkownika',
		password: 'Hasło',
		password_confirm: 'Potwierdź hasło',
		required_field: 'pole wymagane'
	},
	buttons: {
		save: 'Zapisz',
		delete: 'Usuń',
		change: 'Zmień',
		reset: 'Resetuj',
		cancel: 'Anuluj',
		add: 'Nowy',
		confirm: 'Potwierdź'
	},
	welcome: {
		step1_page_title: 'Witaj w UpSnap 🥳',
		step1_setup_desc: 'Wykonaj następujące kroki, aby dokończyć konfigurację.',
		step1_setup_btn_next: 'Dalej',
		step2_page_title: 'Utwórz konto administratora',
		step2_label_email: 'Email:',
		step2_label_password: 'Hasło:',
		step2_label_min_chars: 'min. 10 znaków',
		step2_label_password_confirm: 'Potwierdź hasło:',
		step2_btn_create: 'Utwórz',
		step3_page_title: 'Wszystko gotowe! 🎉',
		step3_page_desc: 'A teraz dodaj kilka urządzeń do Twojego pulpitu.',
		step3_btn_done: 'Lecimy!',
		not_expected_title: 'Stało się coś niespodziewanego! 🧐',
		not_expected_desc: 'Konfiguracja zakończona! Nie musisz nic więcej robić.',
		not_expected_back: 'Powrót',
		progress_step1: 'Witaj',
		progress_step2: 'Utwórz konto',
		progress_step3: 'Gotowe'
	},
	toasts: {
		admin_saved: 'Administrator zapisany',
		user_saved: 'Użytkownik zapisany',
		user_created: 'Użytkownik {username} utworzony',
		user_deleted: 'Użytkownik {username} usunięty',
		settings_saved: 'Zapisano ustawienia',
		password_changed: 'Hasło zostało zmienione. Zaloguj się ponownie.',
		passwords_missmatch: 'Hasła się różnią',
		permissions_created: 'Uprawnienia dla {username} zostały utworzone',
		permissions_deleted: 'Uprawnienia dla {username} zostały usunięte',
		permissions_updated: 'Uprawnienia dla {username} zostały zaktualizowane',
		permissions_updated_personal: 'Twoje uprawnienia zostały zaktualizowane',
		no_permission: 'Nie masz uprawnień aby odwiedzić {url}',
		device_created: 'Utworzono {device}',
		devices_created_multiple: 'Utworzono {count|int} urządzeń',
		device_updated: 'Zaktualizowano {device}',
		device_deleted: 'Usunięto {device}',
		group_created: 'Utworzono grupę {group}',
		group_deleted: 'Usunięto grupę {group}'
	},
	navbar: {
		theme: 'Motyw',
		new: 'Utwórz',
		edit_account: 'Edytuj konto',
		logout: 'Wyloguj'
	}
} satisfies Translation;

export default pl;
