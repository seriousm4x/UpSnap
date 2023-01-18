<script>
	import Navbar from '@components/Navbar.svelte';
	import Transition from '@components/Transition.svelte';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { theme } from '@stores/themestore';

	let preferesDark;

    onMount(() => {
        preferesDark = window.matchMedia('(prefers-color-scheme: dark)');
        preferesDark.addEventListener('change', (e) => {
            if ($theme === 'auto') {
                document.documentElement.setAttribute(
                    'data-bs-theme',
                    e.matches ? 'dark' : 'light'
                );
            }
        });

        theme.subscribe((t) => {
            localStorage.setItem('theme', t);
            if (t === 'auto') {
                t = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
            }
            document.documentElement.setAttribute('data-bs-theme', t);
        });
    });
</script>

<svelte:head>
    <script>
        if (document) {
            preferesDark = window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light';
            let t = localStorage.getItem('theme');
            if (!t) {
                localStorage.setItem('theme', 'auto');
                t = preferesDark;
            } else if (t === 'auto') {
                t = preferesDark;
            }
            document.documentElement.setAttribute('data-bs-theme', t);
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
