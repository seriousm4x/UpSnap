<script lang="ts">
	import { goto } from '$app/navigation';
	import { pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';

	$: console.log($pocketbase);
	$: avatar = $pocketbase.authStore?.model?.avatar ? $pocketbase.authStore?.model.avatar : '0';

	function logout() {
		$pocketbase.authStore.clear();
		goto('/login');
	}
</script>

<div class="navbar bg-base-100">
	<div class="flex-1">
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
	<div class="flex-none">
		{#if $pocketbase.authStore?.model !== null}
			<div class="dropdown dropdown-end">
				<label tabindex="-1" class="btn btn-ghost btn-circle avatar" for="avatar">
					<div class="w-10 rounded-full" id="avatar">
						<img
							src="{$pocketbase.baseUrl}/_/images/avatars/avatar{avatar}.svg"
							alt="Profile pic"
						/>
					</div>
				</label>
				<ul
					tabindex="-1"
					class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-40"
				>
					<li>
						<a class="justify-between" href="/profile/{$pocketbase.authStore.model?.id}"
							>Edit profile</a
						>
					</li>
					<li>
						<div on:click={() => logout()} on:keydown={() => logout()} role="none">Logout</div>
					</li>
				</ul>
			</div>
		{/if}
	</div>
</div>
