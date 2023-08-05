<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { pocketbase, backendUrl, devices } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import Navbar from '$lib/components/Navbar.svelte';
	import Transition from '$lib/components/Transition.svelte';
	import type { SettingsPublic } from '$lib/types/settings';

	onMount(async () => {
		// set settingsPub store on load
		if (!$settingsPub) {
			const res = await $pocketbase.collection('settings_public').getList(1, 1);
			settingsPub.set(res.items[0] as SettingsPublic);
		}

		// redirect to welcome page if setup is not completed
		if (!$settingsPub.setup_completed && $page.url.pathname !== '/welcome') {
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

		if ($pocketbase.authStore.model?.collectionName === 'users') {
			$pocketbase.collection('users').authRefresh();
		} else {
			$pocketbase.admins.authRefresh();
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

{#if $pocketbase.authStore.isValid}
	<Navbar />
{/if}

<Transition url={$page.url}>
	<div class="container mx-auto p-2 mb-4">
		<slot />
	</div>
</Transition>
