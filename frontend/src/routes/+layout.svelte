<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
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

	let authIsValid = false;

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
			authIsValid = $pocketbase.authStore.isValid;

			// load user permissions
			if ($pocketbase.authStore.record?.collectionName === 'users') {
				$pocketbase
					.collection('permissions')
					.getFirstListItem(`user.id = '${$pocketbase.authStore.record.id}'`)
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
		if ($settingsPub.setup_completed === false && page.url.pathname !== '/welcome') {
			$pocketbase.authStore.clear();
			goto('/welcome');
			return;
		}

		// refresh auth token
		if ($pocketbase.authStore.isSuperuser) {
			await $pocketbase
				.collection('_superusers')
				.authRefresh()
				.catch(() => {
					$pocketbase.authStore.clear();
					goto('/login');
				});
		} else {
			await $pocketbase
				.collection('users')
				.authRefresh()
				.catch(() => {
					$pocketbase.authStore.clear();
					goto('/login');
				});
		}
	});
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

{#if authIsValid && !page.url.pathname.startsWith('/welcome')}
	<Navbar />
{/if}

<Toaster position="bottom-center" {toastOptions} />

<Transition url={page.url}>
	<div class="container mx-auto mb-4 p-2">
		<slot />
	</div>
</Transition>
