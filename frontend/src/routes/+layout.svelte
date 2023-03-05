<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import Navbar from '@components/Navbar.svelte';
    import Login from '@components/Login.svelte';
    import Transition from '@components/Transition.svelte';
    import { theme } from '@stores/theme';
    import { pocketbase, authorizedStore } from '@stores/pocketbase';

    let preferesDark;
    let isAuth = false;

    onMount(() => {
        // import bootstrap js
        import('bootstrap/js/dist/dropdown');
        import('bootstrap/js/dist/collapse');

        // set dark mode
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

        authorizedStore.subscribe((state) => {
            isAuth = state;
        });

        const pbCookie = localStorage.getItem('pocketbase_auth');
        if (pbCookie) {
            $pocketbase.authStore.loadFromCookie('pb_auth=' + pbCookie);
            authorizedStore.set($pocketbase.authStore.isValid);
        }
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

{#if isAuth}
    <Navbar />
    <Transition url={$page.url}>
        <slot />
    </Transition>
{:else}
    <Login />
{/if}

<style lang="scss" global>
    @import '../scss/main.scss';
</style>
