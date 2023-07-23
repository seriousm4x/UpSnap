<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { pocketbase, backendUrl, devices } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import Navbar from '$lib/components/Navbar.svelte';

	onMount(async () => {
		// set settingsPub store on load
		if (!$settingsPub) {
			const res = await $pocketbase.collection('settings_public').getList(1, 1);
			settingsPub.set({
				collectionId: res.items[0].collectionId,
				favicon: res.items[0].favicon,
				setup_completed: res.items[0].setup_completed,
				website_title: res.items[0].website_title
			});
		}

		// redirect to welcome page if setup is not completed
		if (!$settingsPub?.setup_completed && $page.url.pathname !== '/welcome') {
			goto('/welcome');
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
		}

		if ($pocketbase.authStore.model?.collectionName === 'users') {
			$pocketbase.collection('users').authRefresh();
		} else {
			$pocketbase.admins.authRefresh();
		}

		// fill devices store
		$pocketbase
			.collection('devices')
			.getFullList(1000, { sort: 'name', expand: 'ports,groups' })
			.then((data) => {
				devices.set(data);
			});
	});
</script>

<svelte:head>
	<link
		rel="shortcut icon"
		href={$settingsPub?.collectionId && $settingsPub?.favicon
			? `${backendUrl}/api/files/settings_public/${$settingsPub?.collectionId}/${$settingsPub?.favicon}`
			: '/gopher.svg'}
	/>
</svelte:head>

{#if $pocketbase.authStore.isValid}
	<Navbar />
{/if}

<div class="container mx-auto">
	<slot />
</div>
