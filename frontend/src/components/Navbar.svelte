<script>
    import { page } from '$app/stores';
    import { theme } from '@stores/theme';
    import { pocketbase, settings } from '@stores/pocketbase';
    import { faSun, faMoon, faCircleHalfStroke, faBrush } from '@fortawesome/free-solid-svg-icons';
    import Fa from 'svelte-fa';
    import { onMount } from 'svelte';

    onMount(async () => {
        $pocketbase.collection('settings').subscribe('*', function (e) {
            settings.set(e.record);
            document.title = $settings.website_title;
        });
    });
</script>

<svelte:head>
    {#if $settings.website_title !== ''}
        <title>{$settings.website_title}</title>
    {:else}
        <title>UpSnap</title>
    {/if}
</svelte:head>

<nav class="navbar navbar-expand">
    <div class="container-fluid">
        <a class="navbar-brand" href="/">
            <img
                src="/favicon.png"
                alt="UpSnap"
                width="30"
                height="30"
                class:me-2={$settings.website_title !== ''}
            />{$settings.website_title ? $settings.website_title : ''}
        </a>
        <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav">
                <li class="nav-item">
                    <a
                        class="nav-link"
                        class:active={$page.url.pathname === '/' ? true : false}
                        href="/">Home</a
                    >
                </li>
                <li class="nav-item">
                    <a
                        class="nav-link"
                        class:active={$page.url.pathname === '/settings' ? true : false}
                        href="/settings">Settings</a
                    >
                </li>
                <li class="nav-item dropdown ms-3">
                    <div
                        class="nav-link dropdown-toggle"
                        role="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        <Fa icon={faBrush} class="me-2" />
                        Theme
                    </div>
                    <ul class="dropdown-menu border-0 p-1">
                        <li>
                            <div
                                class="dropdown-item"
                                class:active={$theme === 'light'}
                                on:click={() => ($theme = 'light')}
                                on:keydown={() => ($theme = 'light')}
                            >
                                <Fa icon={faSun} class="me-2" />
                                Light
                            </div>
                        </li>
                        <li>
                            <div
                                class="dropdown-item"
                                class:active={$theme === 'dark'}
                                on:click={() => ($theme = 'dark')}
                                on:keydown={() => ($theme = 'dark')}
                            >
                                <Fa icon={faMoon} class="me-2" />
                                Dark
                            </div>
                        </li>
                        <li>
                            <div
                                class="dropdown-item"
                                class:active={$theme === 'auto'}
                                on:click={() => ($theme = 'auto')}
                                on:keydown={() => ($theme = 'auto')}
                            >
                                <Fa icon={faCircleHalfStroke} class="me-2" />
                                Auto
                            </div>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
    </div>
</nav>
