<script>
    import { page } from '$app/stores';
    import { theme } from '@stores/theme';
    import { pocketbase } from '@stores/pocketbase';
    import { faSun, faMoon, faCircleHalfStroke, faBrush } from '@fortawesome/free-solid-svg-icons';
    import Fa from 'svelte-fa';
    import { onMount } from 'svelte';

    let pb;
    let website_title = '';

    onMount(async () => {
        pocketbase.subscribe((conn) => {
            pb = conn;
        });
        const result = await pb.collection('settings').getList(1, 1);
        website_title = result.items[0].website_title;

        pb.collection('settings').subscribe('*', function (e) {
            website_title = e.record.website_title;
            document.title = website_title;
        });
    });
</script>

<svelte:head>
    {#if website_title !== ''}
        <title>{website_title}</title>
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
                class:me-2={website_title !== ''}
            />{website_title ? website_title : ''}
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
