<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { pocketbase, isAdmin, backendUrl } from '$lib/stores/pocketbase';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import { faSave } from '@fortawesome/free-solid-svg-icons';
	import type { User } from '$lib/types/user';
	import type { Device } from '$lib/types/device';
	import type { Permission } from '$lib/types/permission';

	let users = [] as User[];
	let permissions = [] as Permission[];
	let devices = [] as Device[];

	$: users.map((user) => {
		if (!permissions.find((perm) => perm.user === user.id)) {
			permissions = [...permissions, { user: user.id } as Permission];
		}
	});

	onMount(() => {
		if (!$isAdmin) {
			toast(`You don't have permission to visit ${$page.url.pathname}`, {
				icon: '⛔'
			});
			goto('/');
			return;
		}
	});

	async function getUsers() {
		$pocketbase
			.collection('users')
			.getFullList(1000, { sort: 'username' })
			.then((data) => {
				users = data as User[];
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	async function getPermissions() {
		$pocketbase
			.collection('permissions')
			.getFullList(1000, { expand: 'user' })
			.then((data) => {
				permissions = data as Permission[];
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	async function getDevices() {
		$pocketbase
			.collection('devices')
			.getFullList(1000, { sort: 'name' })
			.then((resp) => {
				devices = resp as Device[];
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	function save() {
		let errors = 0;
		permissions.forEach(async (permission) => {
			if (permission.id === undefined) {
				await $pocketbase
					.collection('permissions')
					.create(permission)
					.catch((err) => {
						toast.error(err.message);
						errors += 1;
					});
			} else {
				await $pocketbase
					.collection('permissions')
					.update(permission.id, permission)
					.catch((err) => {
						toast.error(err.message);
						errors += 1;
					});
			}
		});
		if (errors === 0) {
			toast.success('Permissions saved.');
		} else {
			toast.error(`Failed to save ${errors} ${errors > 1 ? 'users' : 'user'}.`);
		}
	}
</script>

{#await Promise.all([getUsers(), getPermissions(), getDevices()])}
	<PageLoading />
{:then}
	<h1 class="text-3xl font-bold mb-8">Permissions</h1>
	<form on:submit|preventDefault={save}>
		{#each users as user, index}
			<div class="card w-full bg-base-300 shadow-xl mt-6">
				<div class="card-body gap-4">
					<h2 class="card-title">
						<label tabindex="-1" class="avatar" for="avatar{index}">
							<div class="w-10 rounded-full" id="avatar{index}">
								<img
									src="{backendUrl}_/images/avatars/avatar{Math.floor(
										Math.random() * 9
									).toString()}.svg"
									alt="Profile pic"
								/>
							</div>
						</label>
						{user.username}
					</h2>
					{#each permissions.filter((perm) => perm.user === user.id) as permission}
						<div class="form-control w-fit">
							<label class="label cursor-pointer gap-2">
								<input
									type="checkbox"
									bind:checked={permission.create}
									class="checkbox checked:checkbox-primary"
								/>
								<span class="label-text font-bold"
									>Allow {user.username} to create new devices and edit device groups?</span
								>
							</label>
						</div>
					{/each}
					<div class="grid grid-cols-5 gap-4 place-items-center">
						<div class="font-bold" />
						<div class="font-bold">Read</div>
						<div class="font-bold">Update</div>
						<div class="font-bold">Delete</div>
						<div class="font-bold">Power</div>
						{#each devices as device}
							{#each permissions.filter((perm) => perm.user === user.id) as permission}
								<div class="flex flex-row gap-2 place-self-start">
									<span class="font-bold">{device.name}</span>
									<span class="badge hidden sm:block">{device.ip}</span>
								</div>
								<input
									type="checkbox"
									class="checkbox checked:checkbox-primary"
									bind:group={permission.read}
									value={device.id}
								/>
								<input
									type="checkbox"
									class="checkbox checked:checkbox-primary"
									bind:group={permission.update}
									value={device.id}
								/>
								<input
									type="checkbox"
									class="checkbox checked:checkbox-primary"
									bind:group={permission.delete}
									value={device.id}
								/>
								<input
									type="checkbox"
									class="checkbox checked:checkbox-primary"
									bind:group={permission.power}
									value={device.id}
								/>
							{/each}
						{/each}
					</div>
				</div>
			</div>
		{:else}
			<div class="container text-center">
				<p>No users</p>
			</div>
		{/each}
		{#if users.length > 0}
			<div class="card-actions mt-6 justify-end gap-4">
				<button class="btn btn-success" type="submit"><Fa icon={faSave} />Save</button>
			</div>
		{/if}
	</form>
{/await}
