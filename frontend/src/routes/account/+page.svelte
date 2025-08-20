<script lang="ts">
	import { goto } from '$app/navigation';
	import { m } from '$lib/paraglide/messages';
	import type { Locale } from '$lib/paraglide/runtime';
	import { getLocale, locales, setLocale } from '$lib/paraglide/runtime';
	import { localeStore } from '$lib/stores/locale';
	import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
	import { faSave } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let newAvatar: number;

	let selectedLang: Locale;

	// locales
	const languageEmojis = {
		'bg-BG': '🇧🇬',
		'de-DE': '🇩🇪',
		'en-US': '🇺🇸',
		'es-ES': '🇪🇸',
		'fr-FR': '🇫🇷',
		'id-ID': '🇮🇩',
		'it-IT': '🇮🇹',
		'ja-JP': '🇯🇵',
		'ko-KR': '🇰🇷',
		'nl-NL': '🇳🇱',
		'nb-NO': '🇳🇴',
		'pl-PL': '🇵🇱',
		'pt-PT': '🇵🇹',
		'zh-TW': '🇹🇼',
		'zh-CN': '🇨🇳'
	};

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
		selectedLang = getLocale();
	});

	async function saveUser() {
		if (newAvatar !== undefined) {
			if ($pocketbase.authStore.isSuperuser) {
				if (!$pocketbase.authStore.record?.id) return;
				await $pocketbase
					.collection('_superusers')
					.update($pocketbase.authStore.record.id, { avatar: newAvatar })
					.then(() => {
						toast.success(m.toasts_admin_saved());
					})
					.catch((err) => {
						toast.error(err.message);
					});
			} else {
				if (!$pocketbase.authStore.record?.id) return;
				await $pocketbase
					.collection('users')
					.update($pocketbase.authStore.record.id, { avatar: newAvatar })
					.then(() => {
						toast.success(m.toasts_user_saved());
					})
					.catch((err) => {
						toast.error(err.message);
					});
			}
		}

		if ($localeStore !== selectedLang) {
			setLocale(selectedLang);
		}
	}

	function changePassword() {
		fetch(
			`${backendUrl}api/collections/${
				$pocketbase.authStore.isSuperuser ? '_superusers' : `users`
			}/records/${$pocketbase.authStore.record?.id}`,
			{
				method: 'PATCH',
				headers: {
					Authorization: $pocketbase.authStore.token,
					'Content-Type': 'application/json'
				},
				body: $pocketbase.authStore.isSuperuser ? adminBody : userBody
			}
		)
			.then(async (data) => {
				if (data.ok) {
					toast.success(m.toasts_password_changed());
					$pocketbase.authStore.clear();
					goto('/login');
				} else {
					const j = await data.json();
					if (j?.data?.password?.message) {
						toast.error(j?.data?.password?.message);
					} else if (j?.data?.passwordConfirm?.message) {
						toast.error(m.toasts_passwords_missmatch());
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

	function localeToFullName(lang: Locale) {
		const languageNames = new Intl.DisplayNames([lang], { type: 'language' });
		return languageNames.of(lang) ?? '';
	}
</script>

<h1 class="mb-8 text-3xl font-bold">{m.account_page_title()}</h1>
<div class="card bg-base-200 w-full shadow-sm">
	<div class="card-body">
		<div class="flex flex-row items-center gap-4">
			<div class="avatar">
				<div class="w-16 rounded-full">
					{#if $pocketbase.authStore.record?.id}
						<img
							src="/avatars/avatar{newAvatar ?? $pocketbase.authStore.record?.avatar}.svg"
							alt="Avatar {newAvatar ?? $pocketbase.authStore.record?.avatar}"
						/>
					{/if}
				</div>
			</div>
			<div>
				<h2 class="card-title">
					{$pocketbase.authStore.isSuperuser
						? $pocketbase.authStore.record?.email
						: $pocketbase.authStore.record?.username}
				</h2>
				<h3>
					{$pocketbase.authStore.isSuperuser
						? m.account_account_type_admin()
						: m.account_account_type_user()}
				</h3>
			</div>
		</div>
		<form on:submit|preventDefault={saveUser}>
			<h2 class="card-title mt-4 mb-2">{m.account_avatar_title()}</h2>
			<div class="flex flex-row flex-wrap gap-4">
				{#each [...Array(10).keys()] as i (i)}
					<div class="avatar">
						<div
							class="hover:ring-neutral w-11 rounded-full hover:ring-2 hover:ring-offset-2"
							class:ring-2={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.record?.avatar}
							class:ring-success={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.record?.avatar}
							class:ring-offset-base-100={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.record?.avatar}
							class:ring-offset-2={newAvatar !== undefined
								? i === newAvatar
								: i === $pocketbase.authStore.record?.avatar}
							on:click={() => (newAvatar = i)}
							role="none"
						>
							{#if $pocketbase.authStore.record?.id}
								<img src="/avatars/avatar{i}.svg" alt="{m.account_avatar_title()} {i}" />
							{/if}
						</div>
					</div>
				{/each}
			</div>
			<h2 class="card-title mt-4 mb-2">{m.account_language_title()}</h2>
			<select class="select w-full max-w-xs" bind:value={selectedLang}>
				{#each [...locales].sort( (a, b) => a.localeCompare( b, $localeStore, { sensitivity: 'base' } ) ) as lang (lang)}
					<option value={lang} selected={$localeStore === lang}>
						{languageEmojis[lang]}
						{localeToFullName(lang)} [{lang}]
					</option>
				{/each}
			</select>
			<div class="mt-2">
				<button type="submit" class="btn btn-success mt-2"
					><Fa icon={faSave} />{m.buttons_save()}</button
				>
			</div>
		</form>
	</div>
</div>
<div class="card bg-base-200 mt-6 w-full shadow-sm">
	<div class="card-body">
		<h2 class="card-title">{m.account_change_password_title()}</h2>
		<p>{m.account_change_password_body()}</p>
		<form on:submit|preventDefault={changePassword}>
			<div class="w-full max-w-xs">
				{#if !$pocketbase.authStore.isSuperuser}
					<fieldset class="fieldset">
						<label class="floating-label mt-2">
							<span>{m.account_change_password_label()} <span class="text-error">*</span></span>
							<input
								type="password"
								placeholder={m.account_change_password_label()}
								class="input"
								minlength="5"
								maxlength="72"
								bind:value={newPassword.old}
								required
							/>
						</label>
					</fieldset>
				{/if}
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.account_change_password_new()} <span class="text-error">*</span></span>
						<input
							type="password"
							placeholder={m.account_change_password_new()}
							class="input"
							minlength={$pocketbase.authStore.isSuperuser ? 10 : 5}
							maxlength="72"
							bind:value={newPassword.password}
							required
						/>
					</label>
				</fieldset>
			</div>
			<div class="w-full max-w-xs">
				<fieldset class="fieldset">
					<label class="floating-label mt-2">
						<span>{m.account_change_password_confirm()} <span class="text-error">*</span></span>
						<input
							type="password"
							placeholder={m.account_change_password_confirm()}
							class="input"
							minlength={$pocketbase.authStore.isSuperuser ? 10 : 5}
							maxlength="72"
							bind:value={newPassword.confirm}
							required
						/>
					</label>
				</fieldset>
			</div>
			<div class="mt-2">
				<button type="submit" class="btn btn-success mt-2"
					><Fa icon={faSave} />{m.buttons_change()}</button
				>
			</div>
		</form>
	</div>
</div>
