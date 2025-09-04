<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { toggleVisibility } from '$lib/helpers/forms';
	import { m } from '$lib/paraglide/messages';
	import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import { faEye, faLockOpen, faQuestionCircle } from '@fortawesome/free-solid-svg-icons';
	import type { AuthProviderInfo } from 'pocketbase';
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
			goto(resolve('/'));
		}
	});

	function tryAdminThenUser() {
		$pocketbase
			.collection('_superusers')
			.authWithPassword(form.email, form.password)
			.then(() => {
				goto(resolve('/'));
			})
			.catch(() => {
				$pocketbase
					.collection('users')
					.authWithPassword(form.email, form.password)
					.then(() => {
						goto(resolve('/'));
					})
					.catch((err) => {
						toast.error(err.message);
					});
			});
	}

	function loginWithProvider(provider: AuthProviderInfo) {
		$pocketbase
			.collection('users')
			.authWithOAuth2({ provider: provider.name })
			.then(() => {
				goto(resolve('/'));
			})
			.catch((err) => {
				console.log(err);
				toast.error(err.message);
			});
	}
</script>

<div class="mt-10 flex items-center justify-center">
	<div class="card bg-base-200 w-full max-w-lg shadow-sm">
		<div class="card-body">
			{#if $settingsPub?.website_title}
				<h1 class="text-center text-xl font-bold">{$settingsPub?.website_title}</h1>
			{/if}
			<div class="flex flex-row gap-4">
				<figure class="w-16">
					<img
						src={$settingsPub?.id && $settingsPub?.favicon
							? `${backendUrl}api/files/settings_public/${$settingsPub?.id}/${$settingsPub?.favicon}`
							: '/gopher.svg'}
						alt={$settingsPub?.website_title ? $settingsPub?.website_title : 'Gopher'}
					/>
				</figure>
				<h2 class="card-title">{m.login_welcome()}</h2>
			</div>
			<form
				class="w-full"
				onsubmit={(e) => {
					e.preventDefault();
					tryAdminThenUser();
				}}
			>
				<label class="label" for="email">
					<span>{m.login_email_label()}</span>
				</label>
				<input id="email" type="text" class="input w-full" required bind:value={form.email} />
				<label class="label" for="password">
					<span>{m.login_password_label()}</span>
				</label>
				<label class="relative block">
					<div
						class="absolute top-1/2 right-4 -translate-y-1/2 cursor-pointer"
						role="none"
						onclick={() => toggleVisibility(inputPassword)}
						onkeydown={() => toggleVisibility(inputPassword)}
					>
						<Fa icon={faEye} />
					</div>
					<input
						id="password"
						type="password"
						class="input w-full"
						maxlength="72"
						required
						bind:value={form.password}
						bind:this={inputPassword}
					/>
				</label>
				<div class="card-actions mt-4">
					<div class="dropdown">
						<label tabindex="-1" class="btn" for="other-providers">{m.login_btn_more()}</label>
						<ul
							id="other-providers"
							tabindex="-1"
							class="menu dropdown-content rounded-box bg-base-200 z-1 w-56 p-2 shadow-sm"
						>
							<li class="menu-title flex flex-row items-center gap-2">
								{m.login_menu_title_auth_providers()}
								<a href="https://github.com/seriousm4x/UpSnap/wiki/Auth-providers" target="_blank"
									><Fa icon={faQuestionCircle} /></a
								>
							</li>
							{#await $pocketbase.collection('users').listAuthMethods() then authMethods}
								{#if authMethods.oauth2.enabled && authMethods.oauth2.providers.length > 0}
									{#each authMethods.oauth2.providers as provider (provider)}
										<li>
											<button type="button" onclick={() => loginWithProvider(provider)}
												>{provider.displayName}</button
											>
										</li>
									{/each}
								{/if}
							{/await}
						</ul>
					</div>
					<button class="btn btn-success ms-auto" type="submit"
						>{m.login_btn_login()} <Fa icon={faLockOpen} /></button
					>
				</div>
			</form>
		</div>
	</div>
</div>
