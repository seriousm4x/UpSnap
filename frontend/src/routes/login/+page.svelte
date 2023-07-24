<script lang="ts">
	import { goto } from '$app/navigation';
	import { pocketbase, devices } from '$lib/stores/pocketbase';
	import Fa from 'svelte-fa';
	import { faLockOpen, faEye } from '@fortawesome/free-solid-svg-icons';
	import { toggleVisibility } from '$lib/helpers/forms';
	import type { Device } from '$lib/types/device';

	let inputPassword: HTMLInputElement;
	let form = {
		email: '',
		password: ''
	};
	let errorMsg = '';

	function tryAdminThenUser() {
		errorMsg = '';
		$pocketbase.admins
			.authWithPassword(form.email, form.password)
			.then(() => {
				getDevices();
				goto('/');
			})
			.catch(() => {
				$pocketbase
					.collection('users')
					.authWithPassword(form.email, form.password)
					.then(() => {
						getDevices();
						goto('/');
					})
					.catch((err) => {
						errorMsg = err;
					});
			});
	}

	function getDevices() {
		$pocketbase
			.collection('devices')
			.getFullList(1000, { sort: 'name', expand: 'ports,groups' })
			.then((data) => {
				devices.set(data as Device[]);
			});
	}
</script>

<div class="flex h-screen">
	<div class="m-auto">
		<div class="flex flex-col gap-16 w-screen max-w-lg my-4">
			<div class="card bg-base-300 shadow-xl">
				<div class="card-body">
					<div class="flex flex-row gap-4">
						<figure class="w-16"><img src="/gopher.svg" alt="Gopher" /></figure>
						<h2 class="card-title">Welcome</h2>
					</div>
					<form class="form-control w-full" on:submit|preventDefault={tryAdminThenUser}>
						<label class="label" for="email">
							<span class="label-text">Email or Username:</span>
						</label>
						<input
							id="email"
							type="text"
							class="input input-bordered w-full"
							bind:value={form.email}
						/>
						<label class="label" for="password">
							<span class="label-text">Password:</span>
						</label>
						<label class="relative block">
							<div
								class="absolute top-1/2 -translate-y-1/2 right-4 cursor-pointer"
								role="none"
								on:click={() => toggleVisibility(inputPassword)}
								on:keydown={() => toggleVisibility(inputPassword)}
							>
								<Fa icon={faEye} />
							</div>
							<input
								id="password"
								type="password"
								class="input input-bordered w-full"
								maxlength="72"
								bind:value={form.password}
								bind:this={inputPassword}
							/>
						</label>
						{#if errorMsg !== ''}
							<div class="alert alert-error mt-4">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="stroke-current shrink-0 h-6 w-6"
									fill="none"
									viewBox="0 0 24 24"
									><path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
									/></svg
								>
								<span>{errorMsg}</span>
							</div>
						{/if}
						<div class="card-actions justify-end mt-4">
							<button class="btn btn-primary" type="submit">Login <Fa icon={faLockOpen} /></button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
</div>
