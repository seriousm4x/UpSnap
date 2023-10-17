import en from '../en';
import type { Translation } from '../i18n-types';

const fr = {
	...(en as unknown as Translation),
	home: {
		page_title: 'Accueil',
		order_groups: 'Groupes',
		order_name: 'Nom',
		order_ip: 'IP',
		order_tooltip: 'Tri',
		no_devices: 'Aucun appareil ici.',
		add_first_device: 'Ajoutez votre premier appareil',
		grant_permissions:
			'Veuillez demander à votre administrateur les permissions aux appareils existants ou créez-en de nouveaux.'
	},
	account: {
		page_title: 'Compte',
		account_type_admin: 'Administrateur',
		account_type_user: 'Utilisateur',
		avatar_title: 'Avatar',
		language_title: 'Langue',
		language_option_auto: 'Automatique',
		change_password_title: 'Modification du mot de passe',
		change_password_body: 'Après la modification du mot de passe, vous devrez vous reconnecter.',
		change_password_label: 'Ancien mot de passe',
		change_password_new: 'Nouveau mot de passe',
		change_password_confirm: 'Confirmation du mot de passe'
	},
	device: {
		page_title: 'Nouvel appareil',
		tabs: ['Manuel', 'Analyse du réseau'],
		card_btn_more: 'Plus',
		card_btn_more_edit: 'Modifier',
		card_btn_more_sleep: 'Veille',
		card_btn_more_reboot: 'Redémarrer',
		card_tooltip_wake_cron: 'Réveil avec cron',
		card_tooltip_shutdown_cron: 'Arrêt avec cron',
		card_tooltip_wake_password: 'Mot de passe de réveil',
		card_tooltip_last_status_change: 'Dernier changement de statut',
		card_password: 'Mot de passe',
		card_nic_tooltip_pending: 'En attente',
		card_nic_tooltip_shutdown: 'Éteindre',
		card_nic_tooltip_shutdown_no_cmd: "Aucune commande d'arrêt n'est définie",
		card_nic_tooltip_shutdown_no_permission:
			"Vous n'avez pas la permission pour éteindre cet appareil",
		card_nic_tooltip_power: 'Allumer',
		card_nic_tooltip_power_no_permission: "Vous n'avez pas la permission pour allumer cet appareil",
		general: 'Général',
		general_name: 'Nom',
		general_name_placeholder: "Nom de l'appareil",
		general_ip: 'IP',
		general_mac: 'Mac',
		general_netmask: 'Masque réseau',
		general_required_field: 'champ requis',
		ports: 'Ports',
		ports_desc: 'UpSnap peut aussi vérifier si certains ports sont ouverts.',
		ports_add_new: 'Ajouter un nouveau port',
		ports_name: 'Nom',
		ports_number: 'Numéro',
		link: 'Lien',
		link_desc:
			'Rend le nom de votre appareil cliquable, parfait pour y lier un tableau de bord par exemple.',
		wake: 'Réveil',
		wake_desc: 'Vous pouvez allumer cet appareil en utilisant un job cron planifié.',
		wake_cron: 'Réveil avec cron',
		wake_cron_enable: 'Activer le réveil avec cron',
		sol: 'Sleep-On-LAN',
		sol_desc1:
			'Vous pouvez mettre les ordinateurs en veille à l\'aide de l\'outil <a class="link" href="https://github.com/SR-G/sleep-on-lan" target="_blank">Sleep-On-LAN</a>. Sleep-On-LAN (SOL) est un outil/daemon externe qui fonctionne sur les PC que vous souhaitez mettre en veille, fournissant un point de terminaison REST. Pour obtenir des instructions sur la configuration de Sleep-On-LAN, veuillez vous référer à la section <a href="https://github.com/SR-G/sleep-on-lan#usage" class="link" target="_blank">Usage</a>.',
		sol_desc2:
			"SOL est configuré pour envoyer des requêtes avec le protocole HTTP au lieu d'UDP pour activer l'autorisation et faire des requêtes plus fiables.",
		sol_desc3:
			'Par conséquent, veuillez vous assurer que vous incluez <span class="badge">HTTP:&lt;VOTREPORT&gt;</span> dans la section <span class="badge">Listeners</span> de la <a href="https://github.com/SR-G/sleep-on-lan#configuration" class="link" target="_blank">configuration SOL</a>.',
		sol_enable: 'Activer Sleep-On-LAN',
		sol_port: 'Port SOL',
		sol_authorization: 'Autorisation',
		sol_user: 'Utilisateur SOL',
		sol_password: 'Mot de passe SOL',
		shutdown: 'Arrêt',
		shutdown_desc:
			'Cette <strong>commande shell</strong> va s\'exécuter au sein de votre conteneur (si vous utilisez Docker) ou sur votre hôte (si vous utilisez l\'exécutable). Pour vérifier que cela fonctionne, vous pouvez d\'abord exécuter la commande à l\'intérieur du conteneur ou sur le shell de l\'hôte. Les commandes courantes sont <span class="badge">net rpc</span> pour Windows, <span class="badge">sshpass</span> pour Linux ou <span class="badge">curl</span> en général pour faire des requêtes web.',
		shutdown_examples: 'Exemples:',
		shutdown_examples_windows: 'Arrêter une machine Windows distante:',
		shutdown_examples_linux: 'Arrêter une machine Linux distante:',
		shutdown_cmd: "Commande d'arrêt",
		shutdown_cron_desc:
			"Exactement comme configurer cron pour réveiller l'appareil, vous pouvez aussi planifier un job cron pour arrêter cet appareil.",
		shutdown_cron: 'Arrêt avec cron',
		shutdown_cron_enable: "Activer l'arrêt avec cron",
		password: 'Mot de passe',
		password_desc:
			'Certaines cartes réseau ont l\'option pour définir un mot de passe pour les paquets magiques, aussi appelée <span class="badge">SecureON</span>. Le mot de passe ne peut être que de 0, 4 ou 6 caractères.',
		groups: 'Groupes',
		groups_desc:
			'Vous pouvez ajouter des appareils à un groupe pour pouvoir les trier sur le tableau de bord.',
		groups_placeholder: 'par ex. «Bureau»', // I had to strink because the original English text in French gives a too long text
		network_scan_range_saved: "Plage d'analyse sauvegardée",
		network_scan_desc:
			"Analyser automatiquement votre réseau pour trouver des appareils. Pour que cela fonctionne, vous devez exécuter UpSnap en tant que root/administrateur et avoir nmap installé et disponible dans votre $PATH (pour les utilisateurs Docker, cela est déjà le cas et vous n'avez besoin de rien faire). L'analyse peut prendre quelques secondes.",
		network_scan_ip_range: 'Plage IP',
		network_scan_no_range: "Aucune plage d'analyse",
		network_scan_unsaved_changes: 'Modifications non sauvegardées',
		network_scan_running: 'Analyse en cours',
		network_scan: 'Analyser',
		network_scan_ip: 'IP:',
		network_scan_mac: 'Mac:',
		network_scan_mac_vendor: "Fabricant ayant l'adresse Mac:",
		network_scan_netmask: 'Masque réseau:',
		network_scan_add_all: 'Ajouter tous les appareils',
		network_scan_replace_netmask: 'Remplacer le masque réseau pour tous les appareils ?',
		network_scan_new_netmask: 'Nouveau masque réseau',
		network_scan_include_unknown: 'Inclure les appareils dont nom est « Inconnu »'
	},
	login: {
		welcome: 'Bienvenue',
		email_label: "Courriel ou nom d'utilisateur:",
		password_label: 'Mot de passe:',
		btn_more: 'Plus',
		menu_title_auth_providers: "Autres fournisseurs d'authentification",
		btn_login: 'Connexion'
	},
	settings: {
		page_title: 'Paramètres',
		ping_interval_title: 'Intervalle de ping',
		ping_interval_desc1:
			'Définit l\'intervalle dans lequel les appareils reçoivent une requête ping. Laissez vide pour utiliser la valeur par défaut de <span class="badge">@every 3s</span>.',
		ping_interval_desc2:
			'Pour écrire une syntaxe correcte pour l\'utilitaire cron, consultez <a class="link" href="https://fr.wikipedia.org/wiki/Cron" target="_blank">Wikipédia</a> ou référez-vous à la <a class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3" target="_blank">documentation du paquet</a>.',
		lazy_ping_title: 'Ping fainéant',
		lazy_ping_desc:
			"Lorsque le mode ping fainéant est activé, UpSnap ne va effectuer des requêtes de ping sur les appareils que lorsqu'un utilisateur actif visite le site. Si vous désactivez cette fonction, UpSnap exécutera toujours ces requêtes.",
		lazy_ping_enable: 'Activer',
		website_title_title: 'Titre du site',
		website_title_desc:
			"Définit le titre du site ainsi que le texte affiché dans l'onglet du navigateur.",
		icon_title: 'Icône',
		icon_desc: 'Définit un favicon personnalisé. Les types de fichier compatibles sont:',
		upsnap_version: 'Version de UpSnap'
	},
	users: {
		page_title: 'Utilisateurs',
		allow_create_devices:
			"Autoriser {username} à créer de nouveaux appareils et à modifier les groupes d'appareils",
		device_permissions: 'Permissions des appareils',
		create_new_device: 'Créer un nouvel appareil',
		read: 'Lire',
		update: 'Mettre à jour',
		delete: 'Supprimer',
		power: 'Allumer',
		toggle: 'Inverser',
		confirm_delete_title: 'Confirmation de la suppression',
		confirm_delete_desc: 'Voulez-vous vraiment supprimer {username}?',
		create_new_user: "Création d'un nouvel utilisateur",
		username: "Nom d'utilisateur",
		password: 'Mot de passe',
		password_confirm: 'Confirmation du mot de passe',
		required_field: 'champ requis'
	},
	buttons: {
		save: 'Sauvegarder',
		delete: 'Supprimer',
		change: 'Modifier',
		reset: 'Réinitialiser',
		cancel: 'Annuler',
		add: 'Ajouter'
	},
	welcome: {
		step1_page_title: 'Bienvenue sur UpSnap 🥳',
		step1_setup_desc: "Veuillez suivre les étapes suivantes pour terminer l'installation.",
		step1_setup_btn_next: 'Suivant',
		step2_page_title: 'Création du compte administrateur',
		step2_label_email: 'Courriel:',
		step2_label_password: 'Mot de passe:',
		step2_label_min_chars: 'min. 10 caractères',
		step2_label_password_confirm: 'Confirmation du mot de passe:',
		step2_btn_create: 'Créer',
		step3_page_title: 'Vous êtes tout bon ! 🎉',
		step3_page_desc: 'Allez-y et ajoutez quelques appareils à votre tableau de bord.',
		step3_btn_done: "C'est parti !",
		not_expected_title: 'Je ne vous attendais pas ici ! 🧐',
		not_expected_desc: "Vous avez déjà terminé l'installation ! Il n'y a plus rien à faire.",
		not_expected_back: 'Ramène-moi',
		progress_step1: 'Bienvenue',
		progress_step2: 'Créer un compte',
		progress_step3: 'Fin'
	},
	toasts: {
		admin_saved: 'Administrateur sauvegardé',
		user_saved: 'Utilisateur sauvegardé',
		user_created: 'Utilisateur {username} créé',
		user_deleted: 'Utilisateur {username} supprimé',
		settings_saved: 'Paramètres sauvegardés',
		password_changed: 'Mot de passe modifié. Veuillez vous reconnecter.',
		passwords_missmatch: 'Les mots de passe ne correspondent pas',
		permissions_created: 'Permissions pour {username} créées',
		permissions_deleted: 'Permissions pour {username} supprimées',
		permissions_updated: 'Permissions pour {username} mises à jour',
		permissions_updated_personal: 'Vos permissions ont été mises à jour',
		no_permission: "Vous n'avez pas la permission de visiter {url}",
		device_created: '{device} créé',
		devices_created_multiple: '{count|int} appareils créés',
		device_updated: '{device} mis à jour',
		device_deleted: '{device} supprimé',
		group_created: 'Groupe {group} créé',
		group_deleted: 'Groupe {group} supprimé'
	},
	navbar: {
		theme: 'Thème',
		new: 'Nouveau',
		edit_account: 'Modifier le compte',
		logout: 'Déconnexion'
	}
} satisfies Translation;

export default fr;
