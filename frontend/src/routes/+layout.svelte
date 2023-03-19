<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import Navbar from '@components/Navbar.svelte';
    import Login from '@components/Login.svelte';
    import Transition from '@components/Transition.svelte';
    import { theme } from '@stores/theme';
    import {
        pocketbase,
        authorizedStore,
        devices,
        settings_private,
        settings_public
    } from '@stores/pocketbase';

    let favicon;
    let preferesDark;
    let isAuth = false;

    onMount(async () => {
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

        // load public settings and subscribe to changes and change favicon
        const settingsPublicRes = await $pocketbase.collection('settings_public').getList(1, 1);
        settings_public.set(settingsPublicRes.items[0]);
        updateSettingsPublic();
        $pocketbase.collection('settings_public').subscribe('*', function (e) {
            settings_public.set(e.record);
            updateSettingsPublic();
        });

        // load auth from localstorage
        const pbCookie = localStorage.getItem('pocketbase_auth');
        if (pbCookie) {
            $pocketbase.authStore.loadFromCookie('pb_auth=' + pbCookie);

            // try to refresh auth token if valid
            if ($pocketbase.authStore.isValid) {
                if ($pocketbase.authStore.model?.collectionName === 'users') {
                    await $pocketbase.collection('users').authRefresh();
                } else {
                    await $pocketbase.admins.authRefresh();
                }
            }
            authorizedStore.set($pocketbase.authStore.isValid);
        }
    });

    function updateSettingsPublic() {
        favicon.href =
            $settings_public.favicon === ''
                ? '/gopher.svg'
                : `${$pocketbase.baseUrl}/api/files/settings_public/${$settings_public.id}/${$settings_public.favicon}`;
        document.title =
            $settings_public.website_title === '' ? 'UpSnap' : $settings_public.website_title;
    }

    async function getSettingsPrivateAndDevices() {
        if ($pocketbase.authStore.model?.collectionName !== 'users') {
            const settingsPrivateRes = await $pocketbase
                .collection('settings_private')
                .getList(1, 1);
            settings_private.set(settingsPrivateRes.items[0]);
        }

        let tempDevices = {};
        const devicesRes = await $pocketbase.collection('devices').getFullList(200, {
            sort: 'name',
            expand: 'ports'
        });
        devicesRes.forEach((device) => {
            tempDevices[device.id] = device;
        });
        devices.set(tempDevices);
    }

    $: if (isAuth) {
        getSettingsPrivateAndDevices();
    }
</script>

<svelte:head>
    <link rel="shortcut icon" href="/gopher.svg" bind:this={favicon} />
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
