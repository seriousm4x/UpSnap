<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import Fa from 'svelte-fa';
	import {
		faCog,
		faHome,
		faPlus,
		faKey,
		faSwatchbook,
		faChevronDown,
		faCheck
	} from '@fortawesome/free-solid-svg-icons';
	import { themeChange } from 'theme-change';
	import { onMount } from 'svelte';

	let availableThemes = [
		'light',
		'dark',
		'cupcake',
		'bumblebee',
		'emerald',
		'corporate',
		'synthwave',
		'retro',
		'cyberpunk',
		'valentine',
		'halloween',
		'garden',
		'forest',
		'aqua',
		'lofi',
		'pastel',
		'fantasy',
		'wireframe',
		'black',
		'luxury',
		'dracula',
		'cmyk',
		'autumn',
		'business',
		'acid',
		'lemonade',
		'night',
		'coffee',
		'winter'
	];
	let activeTheme: string | null = '';

	onMount(() => {
		themeChange(false);
		activeTheme = document.documentElement.getAttribute('data-theme');
	});

	$: avatar = $pocketbase.authStore?.model?.avatar
		? $pocketbase.authStore?.model.avatar
		: Math.floor(Math.random() * 9).toString();

	function logout() {
		$pocketbase.authStore.clear();
		goto('/login');
	}
</script>

<div class="navbar bg-base-100">
	<div class="justify-start">
		<div class="dropdown">
			<label tabindex="-1" class="btn btn-ghost md:hidden" for="mobile-menu">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 6h16M4 12h8m-8 6h16"
					/></svg
				>
			</label>
			<ul
				id="mobile-menu"
				tabindex="-1"
				class="menu dropdown-content mt-3 z-[1] p-2 gap-1 shadow bg-base-300 rounded-box w-52"
			>
				<li>
					<a href="/" class="rounded-lg px-4 py-2" class:active={$page.url.pathname === '/'}
						><Fa icon={faHome} />Home</a
					>
				</li>
				{#if !$pocketbase.authStore.model?.collectionName}
					<li>
						<a
							href="/permissions"
							class="rounded-lg px-4 py-2"
							class:active={$page.url.pathname.startsWith('/permissions')}
							><Fa icon={faKey} />Permissions</a
						>
					</li>
					<li>
						<a
							href="/settings/"
							class="rounded-lg px-4 py-2"
							class:active={$page.url.pathname.startsWith('/settings')}
							><Fa icon={faCog} />Settings</a
						>
					</li>
				{/if}
			</ul>
		</div>
		<a class="btn btn-ghost normal-case text-xl" href="/">
			<img
				src={$settingsPub?.id && $settingsPub?.favicon
					? `${backendUrl}api/files/settings_public/${$settingsPub?.id}/${$settingsPub?.favicon}`
					: '/gopher.svg'}
				alt={$settingsPub?.website_title ? $settingsPub?.website_title : 'UpSnap'}
				width="45"
				height="45"
			/>{$settingsPub?.website_title ? $settingsPub?.website_title : ''}
		</a>
	</div>
	<div class="hidden md:flex">
		<ul class="menu menu-horizontal px-1 gap-1">
			<li>
				<a href="/" class="p-2" class:active={$page.url.pathname === '/'}
					><Fa icon={faHome} />Home</a
				>
			</li>
			{#if !$pocketbase.authStore.model?.collectionName}
				<li>
					<a
						href="/permissions"
						class="p-2"
						class:active={$page.url.pathname.startsWith('/permissions')}
						><Fa icon={faKey} />Permissions</a
					>
				</li>
				<li>
					<a href="/settings/" class="p-2" class:active={$page.url.pathname.startsWith('/settings')}
						><Fa icon={faCog} />Settings</a
					>
				</li>
			{/if}
		</ul>
	</div>
	<div class="dropdown md:dropdown-end">
		<div tabindex="-1" class="btn normal-case btn-ghost">
			<Fa icon={faSwatchbook} />
			<span class="hidden font-normal md:inline">Theme</span>
			<Fa icon={faChevronDown} />
		</div>
		<div
			class="dropdown-content bg-base-200 text-base-content rounded-box h-fit max-h-96 w-56 overflow-y-auto shadow mt-3 z-[1]"
		>
			<div class="grid grid-cols-1 gap-3 p-3" tabindex="-1">
				{#each availableThemes as theme}
					<button
						class="outline-base-content overflow-hidden rounded-lg text-left"
						data-set-theme={theme}
						on:click={() => (activeTheme = theme)}
						on:keydown={() => (activeTheme = theme)}
					>
						<div
							data-theme={theme}
							class="bg-base-100 text-base-content w-full cursor-pointer font-sans"
						>
							<div class="grid grid-cols-5 grid-rows-3">
								<div class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3">
									<Fa icon={faCheck} class={activeTheme === theme ? 'visible' : 'invisible'} />
									<div class="flex-grow text-sm">{theme}</div>
									<div class="flex h-full flex-shrink-0 flex-wrap gap-1">
										<div class="bg-primary w-2 rounded" />
										<div class="bg-secondary w-2 rounded" />
										<div class="bg-accent w-2 rounded" />
										<div class="bg-neutral w-2 rounded" />
									</div>
								</div>
							</div>
						</div>
					</button>
				{/each}
			</div>
		</div>
	</div>
	<div class="justify-end ms-auto">
		{#if $pocketbase.authStore?.model !== null}
			<a class="btn btn-success me-4" href="/device/new">
				<Fa icon={faPlus} />
				New
			</a>
			<div class="dropdown dropdown-end">
				<label tabindex="-1" class="btn btn-ghost btn-circle avatar" for="avatar">
					<div class="w-10 rounded-full" id="avatar">
						<img src="{backendUrl}_/images/avatars/avatar{avatar}.svg" alt="Profile pic" />
					</div>
				</label>
				<ul
					tabindex="-1"
					class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-300 rounded-box w-40"
				>
					<li>
						<a class="justify-between" href="/account">Edit profile</a>
					</li>
					<li>
						<div on:click={() => logout()} on:keydown={() => logout()} role="none">Logout</div>
					</li>
				</ul>
			</div>
		{/if}
	</div>
</div>
