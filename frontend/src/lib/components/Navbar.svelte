<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import LL from '$lib/i18n/i18n-svelte';
	import { backendUrl, permission, pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import {
		faCheck,
		faChevronDown,
		faCog,
		faDoorOpen,
		faHome,
		faPlus,
		faSwatchbook,
		faUserGear,
		faUsersGear
	} from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import { themeChange } from 'theme-change';

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
		'winter',
		'dim',
		'nord',
		'sunset',
		'caramellatte',
		'abyss',
		'silk'
	];
	let activeTheme: string | null = '';
	$: avatar = $pocketbase.authStore.record?.avatar;

	onMount(() => {
		themeChange(false);
		activeTheme = document.documentElement.getAttribute('data-theme');

		$pocketbase.authStore.onChange(() => {
			avatar = $pocketbase.authStore.record?.avatar;
		});
	});

	async function logout() {
		await $pocketbase.collection('devices').unsubscribe('*');
		await $pocketbase.collection('ports').unsubscribe('*');
		await $pocketbase.collection('permissions').unsubscribe('*');
		$pocketbase.authStore.clear();
		goto('/login', { invalidateAll: true });
	}
</script>

<div class="navbar bg-base-100">
	<div class="justify-start">
		<div class="dropdown">
			<label
				tabindex="-1"
				class="btn btn-ghost {$settingsPub?.website_title ? 'lg:hidden' : 'md:hidden'}"
				for="mobile-menu"
			>
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
				class="menu dropdown-content rounded-box bg-base-200 z-[1] mt-3 w-max gap-1 p-2 shadow-sm"
			>
				{#if $settingsPub?.website_title}
					<div class="menu-title">
						{$settingsPub?.website_title}
					</div>
				{/if}
				<li>
					<a href="/" class="px-4 py-2" class:active={$page.url.pathname === '/'}
						><Fa icon={faHome} />{$LL.home.page_title()}</a
					>
				</li>
				{#if $pocketbase.authStore.isSuperuser}
					<li>
						<a
							href="/users"
							class="px-4 py-2"
							class:active={$page.url.pathname.startsWith('/users')}
							><Fa icon={faUsersGear} />{$LL.users.page_title()}</a
						>
					</li>
					<li>
						<a
							href="/settings/"
							class="px-4 py-2"
							class:active={$page.url.pathname.startsWith('/settings')}
							><Fa icon={faCog} />{$LL.settings.page_title()}</a
						>
					</li>
				{/if}
				<li>
					<details>
						<summary>
							<Fa icon={faSwatchbook} />
							Themes</summary
						>
						<ul>
							<div class="h-fit max-h-72 overflow-scroll">
								{#each availableThemes as theme}
									<li class="w-full">
										<button
											class="gap-3 px-2"
											data-set-theme={theme}
											on:click={() => (activeTheme = theme)}
											on:keydown={() => (activeTheme = theme)}
										>
											<div
												data-theme={theme}
												class="bg-base-100 grid shrink-0 grid-cols-2 gap-0.5 rounded-md p-1 shadow-sm"
											>
												<div class="bg-base-content size-1 rounded-full"></div>
												<div class="bg-primary size-1 rounded-full"></div>
												<div class="bg-secondary size-1 rounded-full"></div>
												<div class="bg-accent size-1 rounded-full"></div>
											</div>
											<div class="truncate">{theme}</div>
											<Fa icon={faCheck} class={activeTheme === theme ? 'visible' : 'invisible'} />
										</button>
									</li>
								{/each}
							</div>
						</ul>
					</details>
				</li>
			</ul>
		</div>
		<a class="btn btn-ghost h-full border-0 px-2 text-xl" href="/">
			<img
				src={$settingsPub?.id && $settingsPub?.favicon
					? `${backendUrl}api/files/settings_public/${$settingsPub?.id}/${$settingsPub?.favicon}`
					: '/gopher.svg'}
				alt={$settingsPub?.website_title ? $settingsPub?.website_title : 'UpSnap'}
				width="45"
				height="45"
			/>
		</a>
	</div>
	<div class="hidden {$settingsPub?.website_title ? 'lg:flex' : 'md:flex'}">
		{#if $settingsPub?.website_title}
			<span class="px-2">{$settingsPub?.website_title}</span>
		{/if}
		<ul class="menu menu-horizontal h-full gap-1 px-1">
			<li class="h-full">
				<a href="/" class="p-2" class:menu-active={$page.url.pathname === '/'}
					><Fa icon={faHome} />{$LL.home.page_title()}</a
				>
			</li>
			{#if $pocketbase.authStore.isSuperuser}
				<li class="h-full">
					<a href="/users" class="p-2" class:menu-active={$page.url.pathname.startsWith('/users')}
						><Fa icon={faUsersGear} />{$LL.users.page_title()}</a
					>
				</li>
				<li class="h-full">
					<a
						href="/settings/"
						class="p-2"
						class:menu-active={$page.url.pathname.startsWith('/settings')}
						><Fa icon={faCog} />{$LL.settings.page_title()}</a
					>
				</li>
			{/if}
			<div class="dropdown dropdown-end">
				<button class="btn btn-ghost hover:bg-base-content/10 h-full border-0 p-2">
					<Fa icon={faSwatchbook} />
					<span class="font-normal">{$LL.navbar.theme()}</span>
					<Fa icon={faChevronDown} />
				</button>
				<div class="dropdown-content bg-base-200 rounded-box z-1 mt-3 w-52 shadow-sm">
					<ul class="menu menu-horizontal h-fit max-h-96 overflow-y-auto">
						{#each availableThemes as theme}
							<li class="w-full">
								<button
									class="gap-3 px-2"
									data-set-theme={theme}
									on:click={() => (activeTheme = theme)}
									on:keydown={() => (activeTheme = theme)}
								>
									<div
										data-theme={theme}
										class="bg-base-100 grid shrink-0 grid-cols-2 gap-0.5 rounded-md p-1 shadow-sm"
									>
										<div class="bg-base-content size-1 rounded-full"></div>
										<div class="bg-primary size-1 rounded-full"></div>
										<div class="bg-secondary size-1 rounded-full"></div>
										<div class="bg-accent size-1 rounded-full"></div>
									</div>
									<div class="truncate">{theme}</div>
									<Fa icon={faCheck} class={activeTheme === theme ? 'visible' : 'invisible'} />
								</button>
							</li>
						{/each}
					</ul>
				</div>
			</div>
		</ul>
	</div>
	<div class="ms-auto justify-end">
		{#if $pocketbase.authStore?.record !== null}
			{#if $pocketbase.authStore.isSuperuser || $permission.create}
				<a class="btn btn-success me-4" href="/device/new">
					<Fa icon={faPlus} />
					{$LL.navbar.new()}
				</a>
			{/if}
			<div class="dropdown dropdown-end">
				<label tabindex="-1" class="avatar btn btn-circle btn-ghost" for="avatar">
					<div class="w-10 rounded-full" id="avatar">
						<img src="/avatars/avatar{avatar}.svg" alt="Avatar {avatar}" />
					</div>
				</label>
				<ul
					tabindex="-1"
					class="menu dropdown-content rounded-box bg-base-200 z-[1] mt-3 w-52 p-2 shadow"
				>
					<li class="menu-title">
						{$pocketbase.authStore.isSuperuser
							? $pocketbase.authStore.record?.email
							: $pocketbase.authStore.record?.username}
					</li>
					<li>
						<a href="/account"><Fa icon={faUserGear} />{$LL.navbar.edit_account()}</a>
					</li>
					<li>
						<div on:click={async () => logout()} on:keydown={async () => logout()} role="none">
							<Fa icon={faDoorOpen} />{$LL.navbar.logout()}
						</div>
					</li>
				</ul>
			</div>
		{/if}
	</div>
</div>
