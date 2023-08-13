<script lang="ts">
	import { goto } from '$app/navigation';
	import { isAdmin, pocketbase } from '$lib/stores/pocketbase';
	import { backendUrl } from '$lib/stores/pocketbase';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import { faSave } from '@fortawesome/free-solid-svg-icons';

	let newAvatar: number;
	let newPassword = {
		old: '',
		password: '',
		confirm: ''
	};
	$: adminBody = JSON.stringify({
		password: newPassword.password,
		passwordConfirm: newPassword.confirm
	});
	$: userBody = JSON.stringify({
		oldPassword: newPassword.old,
		password: newPassword.password,
		passwordConfirm: newPassword.confirm
	});

	function changeAvatar() {
		if ($isAdmin) {
			$pocketbase.admins
				.authRefresh()
				.then((data) => {
					data.admin.avatar = newAvatar;
					$pocketbase.authStore.save(data.token, data.admin);
					toast.success('Avatar saved.');
				})
				.catch((err) => {
					toast.error(err.message);
				});
		} else {
			$pocketbase
				.collection('users')
				.authRefresh()
				.then((data) => {
					data.record.avatar = newAvatar;
					$pocketbase.authStore.save(data.token, data.record);
					toast.success('Avatar saved.');
				})
				.catch((err) => {
					toast.error(err.message);
				});
		}
	}

	function changePassword() {
		fetch(
			`${backendUrl}api/${$isAdmin ? 'admins' : `collections/users/records`}/${
				$pocketbase.authStore.model?.id
			}`,
			{
				method: 'PATCH',
				headers: {
					Authorization: $pocketbase.authStore.token,
					'Content-Type': 'application/json'
				},
				body: $isAdmin ? adminBody : userBody
			}
		)
			.then(async (data) => {
				if (data.ok) {
					toast.success('Password changed. Please login again.');
					$pocketbase.authStore.clear();
					goto('/login');
				} else {
					const j = await data.json();
					if (j?.data?.password?.message) {
						toast.error(j?.data?.password?.message);
					} else if (j?.data?.passwordConfirm?.message) {
						toast.error("Passwords don't match.");
					} else if (j.data?.oldPassword?.message) {
						toast.error(j.data.oldPassword.message);
					} else {
						toast.error(j.message || data.statusText);
					}
				}
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}
</script>

<h1 class="text-3xl font-bold mb-8">Account</h1>
<div class="card w-full bg-base-300 shadow-xl">
	<div class="card-body">
		<div class="flex flex-row gap-4 items-center">
			<div class="avatar">
				<div class="w-16 rounded-full">
					<img
						src="{backendUrl}_/images/avatars/avatar{newAvatar ??
							$pocketbase.authStore.model?.avatar}.svg"
						alt="Avatar {newAvatar ?? $pocketbase.authStore.model?.avatar}"
					/>
				</div>
			</div>
			<div>
				<h2 class="card-title">
					{$isAdmin ? $pocketbase.authStore.model?.email : $pocketbase.authStore.model?.username}
				</h2>
				<h3>{$isAdmin ? 'Admin' : 'User'}</h3>
			</div>
		</div>
		<h2 class="card-title mt-4">Avatar</h2>
		<form on:submit|preventDefault={changeAvatar}>
			<div class="flex flex-row flex-wrap gap-4">
				{#each [...Array(10).keys()] as i}
					<div class="avatar btn btn-ghost btn-circle">
						<div
							class="w-11 rounded-full"
							class:ring={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.model?.avatar}
							class:ring-primary={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.model?.avatar}
							class:ring-offset-base-100={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.model?.avatar}
							class:ring-offset-2={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.model?.avatar}
							on:click={() => (newAvatar = i)}
							role="none"
						>
							<img src="{backendUrl}_/images/avatars/avatar{i}.svg" alt="Avatar {i}" />
						</div>
					</div>
				{/each}
			</div>
			<button type="submit" class="btn btn-success mt-2"><Fa icon={faSave} />Save</button>
		</form>
	</div>
</div>
<div class="card w-full bg-base-300 shadow-xl mt-6">
	<div class="card-body">
		<h2 class="card-title">Change password</h2>
		<p>After the password has been changed, you will need to log in again.</p>
		<form on:submit|preventDefault={changePassword}>
			<div class="form-control w-full max-w-xs">
				{#if !$isAdmin}
					<label class="label" for="password-old">
						<span class="label-text">Old password</span>
					</label>
					<input
						id="password-old"
						type="password"
						placeholder="Old password"
						class="input input-bordered w-full max-w-xs"
						minlength="5"
						maxlength="72"
						bind:value={newPassword.old}
						required
					/>
				{/if}
				<label class="label" for="password-new">
					<span class="label-text">New password</span>
				</label>
				<input
					id="password-new"
					type="password"
					placeholder="New password"
					class="input input-bordered w-full max-w-xs"
					minlength={$isAdmin ? 10 : 5}
					maxlength="72"
					bind:value={newPassword.password}
					required
				/>
			</div>
			<div class="form-control w-full max-w-xs">
				<label class="label" for="password-confirm">
					<span class="label-text">Confirm password</span>
				</label>
				<input
					id="password-confirm"
					type="password"
					placeholder="New password"
					class="input input-bordered w-full max-w-xs"
					minlength={$isAdmin ? 10 : 5}
					maxlength="72"
					bind:value={newPassword.confirm}
					required
				/>
			</div>
			<button type="submit" class="btn btn-success mt-2"><Fa icon={faSave} />Save</button>
		</form>
	</div>
</div>
