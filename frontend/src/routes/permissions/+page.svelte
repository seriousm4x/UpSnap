<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { isAdmin } from '$lib/stores/pocketbase';
	import toast from 'svelte-french-toast';

	onMount(() => {
		if (!$isAdmin) {
			toast(`You don't have permission to visit ${$page.url.pathname}`, {
				icon: '⛔'
			});
			goto('/');
			return;
		}
	});
</script>

permissions
