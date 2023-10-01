<script lang="ts">
	import { goto } from '$app/navigation';
	import { toggleVisibility } from '$lib/helpers/forms';
	import LL from '$lib/i18n/i18n-svelte';
	import { pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import { faEye, faLockOpen, faQuestionCircle } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let inputPassword: HTMLInputElement;
	let form = {
		email: '',
		password: ''
	};

	onMount(() => {
		if ($pocketbase.authStore.isValid) {
			goto('/');
		}
	});

	function tryAdminThenUser() {
		$pocketbase.admins
			.authWithPassword(form.email, form.password)
			.then(() => {
				goto('/');
			})
			.catch(() => {
				$pocketbase
					.collection('users')
					.authWithPassword(form.email, form.password)
					.then(() => {
						goto('/');
					})
					.catch((err) => {
						toast.error(err.message);
					});
			});
	}

	function loginOIDC() {
		$pocketbase
			.collection('users')
			.authWithOAuth2({ provider: 'oidc' })
			.then(() => {
				goto('/');
			})
			.catch((err) => {
				console.log(err);

				toast.error(err.message);
			});
	}
</script>

<div class="mt-10 flex items-center justify-center">
	<div class="card bg-base-300 shadow-xl w-full max-w-lg">
		<div class="card-body">
			{#if $settingsPub?.website_title}
				<h1 class="text-xl text-center font-bold">{$settingsPub?.website_title}</h1>
			{/if}
			<div class="flex flex-row gap-4">
				<figure class="w-16"><img src="/gopher.svg" alt="Gopher" /></figure>
				<h2 class="card-title">{$LL.login.welcome()}</h2>
			</div>
			<form class="form-control w-full" on:submit|preventDefault={tryAdminThenUser}>
				<label class="label" for="email">
					<span class="label-text">{$LL.login.email_label()}</span>
				</label>
				<input
					id="email"
					type="text"
					class="input input-bordered w-full"
					required
					bind:value={form.email}
				/>
				<label class="label" for="password">
					<span class="label-text">{$LL.login.password_label()}</span>
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
						required
						bind:value={form.password}
						bind:this={inputPassword}
					/>
				</label>
				<div class="card-actions mt-4">
					<div class="dropdown">
						<label tabindex="-1" class="btn btn-neutral" for="other-providers"
							>{$LL.login.btn_more()}</label
						>
						<ul
							id="other-providers"
							tabindex="-1"
							class="dropdown-content z-[1] menu p-2 shadow bg-base-200 rounded-box w-56"
						>
							<li class="menu-title flex flex-row gap-2 items-center">
								{$LL.login.menu_title_auth_providers()}
								<a href="https://github.com/seriousm4x/UpSnap/wiki/Auth-providers" target="_blank"
									><Fa icon={faQuestionCircle} /></a
								>
							</li>
							<li>
								<button type="button" on:click={loginOIDC}>OIDC 1</button>
							</li>
						</ul>
					</div>
					<button class="btn btn-primary ms-auto" type="submit"
						>{$LL.login.btn_login()} <Fa icon={faLockOpen} /></button
					>
				</div>
			</form>
		</div>
	</div>
</div>
