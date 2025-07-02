<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import { m } from '$lib/paraglide/messages';
	import { pocketbase } from '$lib/stores/pocketbase';
	import type { Device } from '$lib/types/device';
	import type { Permission } from '$lib/types/permission';
	import type { User } from '$lib/types/user';
	import { faPlus, faRetweet, faSave, faTrash } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let getUsersPromise = getUsers();
	let getPermissionsPromise = getPermissions();
	let getDevicesPromise = getDevices();
	let deleteModal = [] as HTMLDialogElement[];

	let users = [] as User[];
	let permissions = [] as Permission[];
	let devices = [] as Device[];
	let newUser = {
		username: '',
		password: '',
		passwordConfirm: ''
	};

	$: users.map((user) => {
		if (!permissions.find((perm) => perm.user === user.id)) {
			permissions = [...permissions, { user: user.id } as Permission];
		}
	});

	onMount(() => {
		if (!$pocketbase.authStore.isSuperuser) {
			toast(m.toasts_no_permission({ url: $page.url.pathname }), {
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

	async function createUser() {
		await $pocketbase
			.collection('users')
			.create(newUser)
			.then((data) => {
				toast.success(m.toasts_user_created({ username: data.username }));
				newUser.username = '';
				newUser.password = '';
				newUser.passwordConfirm = '';
				reload();
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}

	async function deleteUser(user: User) {
		let permission = permissions.find((perm) => perm.user === user.id);
		if (permission?.id !== undefined) {
			await $pocketbase
				.collection('permissions')
				.delete(permission.id)
				.then(() => {
					toast.success(m.toasts_permissions_deleted({ username: user.username }));
				})
				.catch((err) => {
					toast.error(err.message);
				});
		}

		await $pocketbase
			.collection('users')
			.delete(user.id)
			.then(async () => {
				toast.success(m.toasts_user_deleted({ username: user.username }));
				reload();
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

	async function save(user: User) {
		let permission = permissions.find((perm) => perm.user === user.id);
		if (permission?.id === undefined) {
			await $pocketbase
				.collection('permissions')
				.create(permission)
				.then((data) => {
					permission = data as Permission;
					// const i = permissions.findIndex((perm) => perm.user === user.id);
					// if (i > -1) {
					// 	permissions[i] = data as Permission;
					// }
					toast.success(m.toasts_permissions_created({ username: user.username }));
					reload();
				})
				.catch((err) => {
					toast.error(err.message);
					return;
				});
		} else {
			await $pocketbase
				.collection('permissions')
				.update(permission.id, permission)
				.then(() => {
					toast.success(m.toasts_permissions_updated({ username: user.username }));
					reload();
				})
				.catch((err) => {
					toast.error(err.message);
				});
		}
	}

	function reload() {
		getUsersPromise = getUsers();
		getPermissionsPromise = getPermissions();
		getDevicesPromise = getDevices();
	}

	function toggleAllPermissions(user: User, key: string) {
		const perm = permissions.find((perm) => perm.user === user.id);
		if (!perm) return;
		if (perm[key] === undefined) perm[key] = [];
		if (devices.length === perm[key].length) {
			// remove all perms
			perm[key] = [];
		} else {
			// grant all perms
			devices.forEach((dev) => {
				if (!perm[key].includes(dev.id)) {
					perm[key] = [...perm[key], dev.id];
				}
			});
		}
		devices = devices;
	}
</script>

{#await Promise.all([getUsersPromise, getPermissionsPromise, getDevicesPromise])}
	<PageLoading />
{:then}
	<h1 class="mb-8 text-3xl font-bold">{m.users_page_title()}</h1>
	{#each users as user, index (user.id)}
		<form on:submit|preventDefault={() => save(user)}>
			<div class="card bg-base-200 mt-6 w-full shadow-sm">
				<div class="card-body gap-4">
					<h2 class="card-title">
						<label tabindex="-1" class="avatar" for="avatar{index}">
							<div class="w-10 rounded-full" id="avatar{index}">
								<img src="/avatars/avatar{user.avatar}.svg" alt="Avatar {user.avatar}" />
							</div>
						</label>
						{user.username}
					</h2>
					{#each permissions.filter((perm) => perm.user === user.id) as permission (permission.id)}
						<div class="w-fit">
							<label class="label cursor-pointer gap-2 text-wrap">
								<input type="checkbox" bind:checked={permission.create} class="checkbox" />
								{m.users_allow_create_devices({ username: user.username })}
							</label>
						</div>
					{/each}
					<div class="collapse-arrow bg-base-200 collapse">
						<input type="checkbox" />
						<div class="collapse-title text-xl font-medium">{m.users_device_permissions()}</div>
						<div class="collapse-content">
							{#if devices.length === 0}
								<p>
									{m.home_no_devices()}
									<a href="/device/new" class="link">{m.users_create_new_device()}</a>
								</p>
							{:else}
								<div class="grid grid-cols-4 justify-items-center gap-4 md:grid-cols-5">
									<div class="font-bold md:col-start-2">{m.users_read()}</div>
									<div class="font-bold">{m.users_update()}</div>
									<div class="font-bold">{m.users_delete()}</div>
									<div class="font-bold">{m.users_power()}</div>
									{#each devices as device (device.id)}
										<hr class="border-neutral col-span-full w-full border-b-1 opacity-30" />
										{#each permissions.filter((perm) => perm.user === user.id) as permission (permission.id)}
											<div
												class="col-span-full flex flex-row flex-wrap gap-2 place-self-start break-all md:col-span-1"
											>
												<span class="font-bold">{device.name}</span>
												<span class="badge hidden md:block">{device.ip}</span>
											</div>
											<input
												type="checkbox"
												class="checkbox md:col-start-2"
												bind:group={permission.read}
												value={device.id}
											/>
											<input
												type="checkbox"
												class="checkbox"
												bind:group={permission.update}
												value={device.id}
											/>
											<input
												type="checkbox"
												class="checkbox"
												bind:group={permission.delete}
												value={device.id}
											/>
											<input
												type="checkbox"
												class="checkbox"
												bind:group={permission.power}
												value={device.id}
											/>
										{/each}
									{/each}
									<button
										class="btn btn-sm md:col-start-2"
										type="button"
										on:click={() => {
											toggleAllPermissions(user, 'read');
										}}
										><span class="hidden md:block">{m.users_toggle()}</span><Fa
											icon={faRetweet}
											size="1.5x"
										/></button
									>
									<button
										class=" btn btn-sm"
										type="button"
										on:click={() => {
											toggleAllPermissions(user, 'update');
										}}
										><span class="hidden md:block">{m.users_toggle()}</span><Fa
											icon={faRetweet}
											size="1.5x"
										/></button
									>
									<button
										class=" btn btn-sm"
										type="button"
										on:click={() => {
											toggleAllPermissions(user, 'delete');
										}}
										><span class="hidden md:block">{m.users_toggle()}</span><Fa
											icon={faRetweet}
											size="1.5x"
										/></button
									>
									<button
										class=" btn btn-sm"
										type="button"
										on:click={() => {
											toggleAllPermissions(user, 'power');
										}}
										><span class="hidden md:block">{m.users_toggle()}</span><Fa
											icon={faRetweet}
											size="1.5x"
										/></button
									>
								</div>
								<div class="mt-4 flex flex-row flex-wrap justify-end gap-4"></div>
							{/if}
						</div>
					</div>
					<div class="join justify-end">
						<button
							class="btn btn-error join-item"
							type="button"
							on:click={() => deleteModal[index].showModal()}
							><Fa icon={faTrash} />{m.buttons_delete()}</button
						>
						<button class="btn btn-success join-item" type="submit"
							><Fa icon={faSave} />{m.buttons_save()}</button
						>
					</div>
				</div>
			</div>
		</form>
		<dialog class="modal" bind:this={deleteModal[index]}>
			<div class="modal-box">
				<h3 class="text-lg font-bold">{m.users_confirm_delete_title()}</h3>
				<p class="py-4">{m.users_confirm_delete_desc({ username: user.username })}</p>
				<div class="modal-action">
					<form method="dialog">
						<button class="btn">{m.buttons_cancel()}</button>
						<button class="btn btn-error" on:click={() => deleteUser(user)}
							>{m.buttons_delete()}</button
						>
					</form>
				</div>
			</div>
		</dialog>
	{/each}
	<div class="card bg-base-200 mt-6 w-full shadow-sm">
		<div class="card-body">
			<h2 class="card-title">{m.users_create_new_user()}</h2>
			<form on:submit|preventDefault={createUser}>
				<div class="flex flex-row flex-wrap gap-4">
					<div class="w-full max-w-xs">
						<label class="label" for="username">
							<div>
								<span>{m.users_username()}</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="username"
							type="text"
							placeholder={m.users_username()}
							class="input w-full max-w-xs"
							required
							bind:value={newUser.username}
						/>
					</div>
					<div class="w-full max-w-xs">
						<label class="label" for="password">
							<div>
								<span>{m.users_password()}</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="password"
							type="password"
							placeholder={m.users_password()}
							class="input w-full max-w-xs"
							minlength="5"
							maxlength="72"
							required
							bind:value={newUser.password}
						/>
					</div>
					<div class="w-full max-w-xs">
						<label class="label" for="passwordConfirm">
							<div>
								<span>{m.users_password_confirm()}</span>
								<span class="text-error">*</span>
							</div>
						</label>
						<input
							id="passwordConfirm"
							type="password"
							placeholder={m.users_password_confirm()}
							class="input w-full max-w-xs"
							minlength="5"
							maxlength="72"
							required
							bind:value={newUser.passwordConfirm}
						/>
					</div>
				</div>
				<span class="badge text-error mt-4 self-center">* {m.users_required_field()}</span>
				<div class="card-actions justify-end">
					<button type="submit" class="btn btn-success mt-2"
						><Fa icon={faPlus} />{m.buttons_add()}</button
					>
				</div>
			</form>
		</div>
	</div>
{/await}
