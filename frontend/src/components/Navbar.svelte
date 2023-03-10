<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { theme } from '@stores/theme';
    import { pocketbase, authorizedStore, settings_public } from '@stores/pocketbase';
    import Fa from 'svelte-fa/src/fa.svelte';
    import {
        faSun,
        faMoon,
        faCircleHalfStroke,
        faBrush,
        faRightFromBracket
    } from '@fortawesome/free-solid-svg-icons';

    let userInfo = {
        usernameOrEmail: '',
        role: ''
    };

    onMount(async () => {
        if ($pocketbase.authStore.baseModel?.collectionName === 'users') {
            userInfo.role = 'user';
        } else {
            userInfo.role = 'admin';
        }
        if ($pocketbase.authStore.baseModel?.username) {
            userInfo.usernameOrEmail = $pocketbase.authStore.baseModel.username;
        } else {
            userInfo.usernameOrEmail = $pocketbase.authStore.baseModel.email;
        }
    });

    async function logout() {
        $pocketbase.authStore.clear();
        authorizedStore.set(false);
    }
</script>

<nav class="navbar navbar-expand-sm pt-0">
    <div class="container-fluid">
        <a class="navbar-brand" href="/">
            <img
                src={$settings_public.favicon !== ''
                    ? `${$pocketbase.baseUrl}/api/files/settings_public/${$settings_public.id}/${$settings_public.favicon}`
                    : '/gopher.svg'}
                alt={$settings_public.website_title ? $settings_public.website_title : 'UpSnap'}
                width="45"
                height="45"
                class:me-2={$settings_public.website_title !== ''}
            />{$settings_public.website_title ? $settings_public.website_title : ''}
        </a>
        <button
            class="navbar-toggler border-0"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent"
            aria-expanded="false"
            aria-label="Toggle navigation"
        >
            <span class="navbar-toggler-icon" />
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav">
                {#if userInfo.role !== 'user'}
                    <li class="nav-item">
                        <a
                            class="nav-link"
                            class:active={$page.url.pathname === '/' ? true : false}
                            href="/">Home</a
                        >
                    </li>
                    <li class="nav-item me-3">
                        <a
                            class="nav-link"
                            class:active={$page.url.pathname === '/settings/' ? true : false}
                            href="/settings/">Settings</a
                        >
                    </li>
                {/if}
                <li class="nav-item dropdown">
                    <div
                        class="nav-link dropdown-toggle"
                        role="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        <Fa icon={faBrush} class="me-2" />
                        Theme
                    </div>
                    <ul class="dropdown-menu border-0 p-1 shadow-sm mb-2">
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
            <div class="ms-auto d-flex">
                <button
                    class="text-dark-emphasis nav-link active border-0"
                    data-toggle="tooltip"
                    title="Logged in as {userInfo.role} &quot;{userInfo.usernameOrEmail}&quot;"
                    on:click={() => logout()}
                >
                    <Fa icon={faRightFromBracket} class="me-2" />Logout ({userInfo.usernameOrEmail})
                </button>
            </div>
        </div>
    </div>
</nav>
