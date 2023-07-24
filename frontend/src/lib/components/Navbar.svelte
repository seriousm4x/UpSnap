<script lang="ts">
	import { goto } from '$app/navigation';
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import Fa from 'svelte-fa';
	import { faCog, faHome, faCircleQuestion, faKey } from '@fortawesome/free-solid-svg-icons';

	$: avatar = $pocketbase.authStore?.model?.avatar
		? $pocketbase.authStore?.model.avatar
		: Math.floor(Math.random() * 9).toString();

	function logout() {
		$pocketbase.authStore.clear();
		goto('/login');
	}
</script>

<div class="navbar bg-base-100">
	<div class="navbar-start">
		<a class="btn btn-ghost normal-case text-xl" href="/">
			<img
				src={$settingsPub?.collectionId && $settingsPub?.favicon
					? `/api/files/settings_public/${$settingsPub?.collectionId}/${$settingsPub?.favicon}`
					: '/gopher.svg'}
				alt={$settingsPub?.website_title ? $settingsPub?.website_title : 'UpSnap'}
				width="45"
				height="45"
			/>{$settingsPub?.website_title ? $settingsPub?.website_title : ''}
		</a>
	</div>
	<div class="navbar-center hidden md:flex">
		<ul class="menu menu-horizontal px-1">
			<li><a href="/"><Fa icon={faHome} />Home</a></li>
			{#if $pocketbase.authStore.model?.collectionName !== 'users'}
				<li><a href="/"><Fa icon={faKey} />Permissions</a></li>
				<li><a href="/"><Fa icon={faCog} />Settings</a></li>
			{/if}
			<li>
				<details>
					<summary><Fa icon={faCircleQuestion} />FAQ</summary>
					<ul class="p-2 bg-base-300">
						<li><a href="/">Submenu 1</a></li>
						<li><a href="/">Submenu 2</a></li>
					</ul>
				</details>
			</li>
		</ul>
	</div>
	<div class="navbar-end">
		{#if $pocketbase.authStore?.model !== null}
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
