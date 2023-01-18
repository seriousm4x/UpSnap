<script>
	import Navbar from '@components/Navbar.svelte';
	import Transition from '@components/Transition.svelte';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	let preferesDark;

	onMount(() => {
		preferesDark = window.matchMedia('(prefers-color-scheme: dark)');
		preferesDark.addEventListener('change', (e) => {
			setTheme(e.matches ? 'dark' : 'light');
		});
	});

	function setTheme(newTheme) {
		document.documentElement.setAttribute('data-bs-theme', newTheme);
		localStorage.setItem('theme', newTheme);
	}
</script>

<svelte:head>
	<script>
		if (document) {
			theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
			document.documentElement.setAttribute('data-bs-theme', theme);
		}
	</script>
</svelte:head>

<Navbar />

<Transition url={$page.url}>
	<slot />
</Transition>

<style lang="scss" global>
	@import '../scss/main.scss';
</style>
