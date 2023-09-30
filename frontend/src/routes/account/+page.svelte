<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import { faSave } from '@fortawesome/free-solid-svg-icons';
	import { loadLocaleAsync } from '$lib/i18n/i18n-util.async';
	import { setLocale } from '$lib/i18n/i18n-svelte';
	import { baseLocale, locales } from '$lib/i18n/i18n-util';
	import { detectLocale, navigatorDetector } from 'typesafe-i18n/detectors';
	import LL from '$lib/i18n/i18n-svelte';
	import type { Locales } from '$lib/i18n/i18n-types';

	let newAvatar: number;

	// locales
	const languageEmojis = {
		'en-US': '🇺🇸',
		'de-DE': '🇩🇪'
	};
	let localStorageLang: Locales | 'auto' = 'auto';
	let selectedLanguage: Locales | 'auto' = localStorageLang;
	$: selectedLanguage = localStorageLang as Locales;

	// password change
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

	onMount(() => {
		localStorageLang = (localStorage.getItem('lang') as Locales) || 'auto';
	});

	async function saveUser() {
		let newLang = selectedLanguage;
		if (newLang === 'auto') {
			localStorage.removeItem('lang');
			newLang = detectLocale(baseLocale, locales, navigatorDetector);
		} else {
			localStorage.setItem('lang', newLang);
		}
		await loadLocaleAsync(newLang);
		setLocale(newLang);

		if ($pocketbase.authStore.isAdmin) {
			if (!$pocketbase.authStore.model?.id) return;
			$pocketbase.admins
				.update($pocketbase.authStore.model.id, { avatar: newAvatar })
				.then(() => {
					toast.success($LL.account.toast_admin_saved());
				})
				.catch((err) => {
					toast.error(err.message);
				});
		} else {
			if (!$pocketbase.authStore.model?.id) return;
			$pocketbase
				.collection('users')
				.update($pocketbase.authStore.model.id, { avatar: newAvatar })
				.then(() => {
					toast.success($LL.account.toast_user_saved());
				})
				.catch((err) => {
					toast.error(err.message);
				});
		}
	}

	function changePassword() {
		fetch(
			`${backendUrl}api/${$pocketbase.authStore.isAdmin ? 'admins' : `collections/users/records`}/${
				$pocketbase.authStore.model?.id
			}`,
			{
				method: 'PATCH',
				headers: {
					Authorization: $pocketbase.authStore.token,
					'Content-Type': 'application/json'
				},
				body: $pocketbase.authStore.isAdmin ? adminBody : userBody
			}
		)
			.then(async (data) => {
				if (data.ok) {
					toast.success($LL.account.toast_password_changed());
					$pocketbase.authStore.clear();
					goto('/login');
				} else {
					const j = await data.json();
					if (j?.data?.password?.message) {
						toast.error(j?.data?.password?.message);
					} else if (j?.data?.passwordConfirm?.message) {
						toast.error($LL.account.toast_passwords_missmatch());
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

	function localeToFullName(lang: Locales) {
		const languageNames = new Intl.DisplayNames([lang], { type: 'language' });
		return languageNames.of(lang) ?? '';
	}
</script>

<h1 class="text-3xl font-bold mb-8">{$LL.account.page_title()}</h1>
<div class="card w-full bg-base-300 shadow-xl">
	<div class="card-body">
		<div class="flex flex-row gap-4 items-center">
			<div class="avatar">
				<div class="w-16 rounded-full">
					<!-- svelte static build will fail, because the image gets served from pocketbase
								and is not a local static file -->
					{#if $pocketbase.authStore.model?.id}
						<img
							src="{backendUrl}_/images/avatars/avatar{newAvatar ??
								$pocketbase.authStore.model?.avatar}.svg"
							alt="Avatar {newAvatar ?? $pocketbase.authStore.model?.avatar}"
						/>
					{/if}
				</div>
			</div>
			<div>
				<h2 class="card-title">
					{$pocketbase.authStore.isAdmin
						? $pocketbase.authStore.model?.email
						: $pocketbase.authStore.model?.username}
				</h2>
				<h3>
					{$pocketbase.authStore.isAdmin
						? $LL.account.account_type_admin()
						: $LL.account.account_type_user()}
				</h3>
			</div>
		</div>
		<form on:submit|preventDefault={saveUser}>
			<h2 class="card-title mt-4 mb-2">{$LL.account.avatar_title()}</h2>
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
							<!-- svelte static build will fail, because the image gets served from pocketbase
								and is not a local static file -->
							{#if $pocketbase.authStore.model?.id}
								<img
									src="{backendUrl}_/images/avatars/avatar{i}.svg"
									alt="{$LL.account.avatar_title()} {i}"
								/>
							{/if}
						</div>
					</div>
				{/each}
			</div>
			<h2 class="card-title mt-4 mb-2">{$LL.account.language_title()}</h2>
			<select class="select select-bordered w-full max-w-xs" bind:value={selectedLanguage}>
				<option value="auto" selected={localStorageLang === null}
					>🌐 {$LL.account.language_option_auto()}</option
				>
				{#each locales.sort( (a, b) => a.localeCompare( b, undefined, { sensitivity: 'base' } ) ) as lang}
					<option value={lang} selected={localStorageLang === lang}>
						{languageEmojis[lang]}
						{localeToFullName(lang)} [{lang}]
					</option>
				{/each}
			</select>
			<div class="mt-2">
				<button type="submit" class="btn btn-success mt-2"
					><Fa icon={faSave} />{$LL.buttons.save()}</button
				>
			</div>
		</form>
	</div>
</div>
<div class="card w-full bg-base-300 shadow-xl mt-6">
	<div class="card-body">
		<h2 class="card-title">{$LL.account.change_password_title()}</h2>
		<p>{$LL.account.change_password_body()}</p>
		<form on:submit|preventDefault={changePassword}>
			<div class="form-control w-full max-w-xs">
				{#if !$pocketbase.authStore.isAdmin}
					<label class="label" for="password-old">
						<span class="label-text">{$LL.account.change_password_label()}</span>
					</label>
					<input
						id="password-old"
						type="password"
						placeholder={$LL.account.change_password_label()}
						class="input input-bordered w-full max-w-xs"
						minlength="5"
						maxlength="72"
						bind:value={newPassword.old}
						required
					/>
				{/if}
				<label class="label" for="password-new">
					<span class="label-text">{$LL.account.change_password_new()}</span>
				</label>
				<input
					id="password-new"
					type="password"
					placeholder={$LL.account.change_password_new()}
					class="input input-bordered w-full max-w-xs"
					minlength={$pocketbase.authStore.isAdmin ? 10 : 5}
					maxlength="72"
					bind:value={newPassword.password}
					required
				/>
			</div>
			<div class="form-control w-full max-w-xs">
				<label class="label" for="password-confirm">
					<span class="label-text">{$LL.account.change_password_confirm()}</span>
				</label>
				<input
					id="password-confirm"
					type="password"
					placeholder={$LL.account.change_password_confirm()}
					class="input input-bordered w-full max-w-xs"
					minlength={$pocketbase.authStore.isAdmin ? 10 : 5}
					maxlength="72"
					bind:value={newPassword.confirm}
					required
				/>
			</div>
			<div class="mt-2">
				<button type="submit" class="btn btn-success mt-2"
					><Fa icon={faSave} />{$LL.buttons.change()}</button
				>
			</div>
		</form>
	</div>
</div>
