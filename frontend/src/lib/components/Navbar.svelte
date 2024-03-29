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
		'sunset'
	];
	let activeTheme: string | null = '';
	let avatar = $pocketbase.authStore.model?.avatar;

	onMount(() => {
		themeChange(false);
		activeTheme = document.documentElement.getAttribute('data-theme');

		$pocketbase.authStore.onChange(() => {
			avatar = $pocketbase.authStore.model?.avatar;
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
				class="menu dropdown-content mt-3 z-[1] p-2 gap-1 shadow bg-base-300 rounded-box w-max"
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
				{#if $pocketbase.authStore.isAdmin}
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
									<li>
										<button
											class="outline-base-content overflow-hidden rounded-lg text-left"
											data-set-theme={theme}
											on:click={() => (activeTheme = theme)}
											on:keydown={() => (activeTheme = theme)}
										>
											<div
												data-theme={theme}
												class="bg-base-100 text-base-content rounded w-full cursor-pointer font-sans"
											>
												<div class="grid grid-cols-5 grid-rows-3">
													<div
														class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"
													>
														<Fa
															icon={faCheck}
															class={activeTheme === theme ? 'visible' : 'invisible'}
														/>
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
									</li>
								{/each}
							</div>
						</ul>
					</details>
				</li>
			</ul>
		</div>
		<a class="btn btn-ghost normal-case text-xl px-2" href="/">
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
		<ul class="menu menu-horizontal px-1 gap-1">
			<li>
				<a href="/" class="p-2" class:active={$page.url.pathname === '/'}
					><Fa icon={faHome} />{$LL.home.page_title()}</a
				>
			</li>
			{#if $pocketbase.authStore.isAdmin}
				<li>
					<a href="/users" class="p-2" class:active={$page.url.pathname.startsWith('/users')}
						><Fa icon={faUsersGear} />{$LL.users.page_title()}</a
					>
				</li>
				<li>
					<a href="/settings/" class="p-2" class:active={$page.url.pathname.startsWith('/settings')}
						><Fa icon={faCog} />{$LL.settings.page_title()}</a
					>
				</li>
			{/if}
		</ul>
		<div class="dropdown dropdown-end">
			<div tabindex="-1" class="btn normal-case btn-ghost">
				<Fa icon={faSwatchbook} />
				<span class="font-normal">{$LL.navbar.theme()}</span>
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
	</div>
	<div class="justify-end ms-auto">
		{#if $pocketbase.authStore?.model !== null}
			{#if $pocketbase.authStore.isAdmin || $permission.create}
				<a class="btn btn-success me-4" href="/device/new">
					<Fa icon={faPlus} />
					{$LL.navbar.new()}
				</a>
			{/if}
			<div class="dropdown dropdown-end">
				<label tabindex="-1" class="btn btn-ghost btn-circle avatar" for="avatar">
					<div class="w-10 rounded-full" id="avatar">
						<img src="{backendUrl}_/images/avatars/avatar{avatar}.svg" alt="Avatar {avatar}" />
					</div>
				</label>
				<ul
					tabindex="-1"
					class="menu dropdown-content mt-3 z-[1] p-2 shadow bg-base-300 rounded-box"
				>
					<li class="menu-title">
						{$pocketbase.authStore.isAdmin
							? $pocketbase.authStore.model.email
							: $pocketbase.authStore.model.username}
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
