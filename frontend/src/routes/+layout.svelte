<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Navbar from '$lib/components/Navbar.svelte';
	import Transition from '$lib/components/Transition.svelte';
	import LL, { setLocale } from '$lib/i18n/i18n-svelte';
	import { baseLocale, locales } from '$lib/i18n/i18n-util';
	import { loadLocaleAsync } from '$lib/i18n/i18n-util.async';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import type { Permission } from '$lib/types/permission';
	import type { SettingsPublic } from '$lib/types/settings';
	import { onMount } from 'svelte';
	import toast, { Toaster, type ToastOptions } from 'svelte-french-toast';
	import { detectLocale, localStorageDetector, navigatorDetector } from 'typesafe-i18n/detectors';
	import '../app.css';

	const toastOptions: ToastOptions = {
		duration: 5000
	};

	onMount(async () => {
		// load locale
		const detectedLocale = detectLocale(
			baseLocale,
			locales,
			localStorageDetector,
			navigatorDetector
		);
		await loadLocaleAsync(detectedLocale);
		setLocale(detectedLocale);

		$pocketbase.authStore.onChange(() => {
			// load user permissions
			if ($pocketbase.authStore.model?.collectionName === 'users') {
				$pocketbase
					.collection('permissions')
					.getFirstListItem(`user.id = '${$pocketbase.authStore.model.id}'`)
					.then((data) => {
						permission.set(data as Permission);
					})
					.catch(() => {
						return;
					});

				$pocketbase.collection('permissions').subscribe('*', (event) => {
					permission.set(event.record as Permission);
					toast.success($LL.toasts.permissions_updated_personal());
				});
			}
		});

		// set settingsPub store on load
		if (!$settingsPub) {
			const res = await $pocketbase.collection('settings_public').getFirstListItem('');
			settingsPub.set(res as SettingsPublic);
		}

		// redirect to welcome page if setup is not completed
		if ($settingsPub.setup_completed === false && $page.url.pathname !== '/welcome') {
			goto('/welcome');
			return;
		}

		// load auth from localstorage
		const pbCookie = localStorage.getItem('pocketbase_auth');
		if (!pbCookie) {
			goto('/login');
			return;
		}

		$pocketbase.authStore.loadFromCookie('pb_auth=' + pbCookie);
		if (!$pocketbase.authStore.isValid) {
			goto('/login');
			return;
		}

		// only refresh token if valid less than 1 day
		const jwt = parseJwt($pocketbase.authStore.token);
		if (jwt.exp > Date.now() / 1000 + 60 * 60 * 24) {
			return;
		}

		if ($pocketbase.authStore.isAdmin) {
			await $pocketbase.admins.authRefresh().catch(() => {
				goto('/login');
			});
		} else {
			await $pocketbase
				.collection('users')
				.authRefresh()
				.catch(() => {
					goto('/login');
				});
		}
	});

	function parseJwt(token: string) {
		var base64Url = token.split('.')[1];
		var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
		var jsonPayload = decodeURIComponent(
			window
				.atob(base64)
				.split('')
				.map(function (c) {
					return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
				})
				.join('')
		);

		return JSON.parse(jsonPayload);
	}
</script>

<svelte:head>
	<link
		rel="shortcut icon"
		href={$settingsPub?.id && $settingsPub?.favicon
			? `${backendUrl}api/files/settings_public/${$settingsPub?.id}/${$settingsPub?.favicon}`
			: '/gopher.svg'}
	/>
	{#if $settingsPub === undefined}
		<title>UpSnap</title>
	{:else}
		<title>{$settingsPub.website_title === '' ? 'UpSnap' : $settingsPub.website_title}</title>
	{/if}
</svelte:head>

{#if $pocketbase.authStore.isValid}
	<Navbar />
{/if}

<Toaster position="bottom-center" {toastOptions} />

<Transition url={$page.url}>
	<div class="container mx-auto p-2 mb-4">
		<slot />
	</div>
</Transition>
