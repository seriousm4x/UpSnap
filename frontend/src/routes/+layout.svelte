<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { theme } from '@stores/theme';
    import Navbar from '@components/Navbar.svelte';
    import Transition from '@components/Transition.svelte';
    import { pocketbase, settings, devices } from '@stores/pocketbase';

    let preferesDark;

    onMount(async () => {
        let settingsRes = {};
        settingsRes = await $pocketbase.collection('settings').getList(1, 1);
        settings.set(settingsRes.items[0]);

        let tempDevices = {};
        const devicesRes = await $pocketbase.collection('devices').getFullList(200, {
            sort: 'name',
            expand: 'ports'
        });
        devicesRes.forEach((device) => {
            tempDevices[device.id] = device;
        });
        devices.set(tempDevices);
    });

    onMount(() => {
        // import bootstrap js
        import('bootstrap/js/dist/dropdown');

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
