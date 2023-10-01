<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { PUBLIC_VERSION } from '$env/static/public';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import LL from '$lib/i18n/i18n-svelte';
	import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
	import { settingsPriv, settingsPub } from '$lib/stores/settings';
	import type { SettingsPrivate, SettingsPublic } from '$lib/types/settings';
	import { faSave } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let settingsPubClone: SettingsPublic | undefined;
	let settingsPrivClone: SettingsPrivate | undefined;
	let faviconPreview: HTMLImageElement;
	let faviconInputElement: HTMLInputElement;

	onMount(() => {
		if (!$pocketbase.authStore.isAdmin) {
			toast($LL.toasts.no_permission({ url: $page.url.pathname }), {
				icon: '⛔'
			});
			goto('/');
			return;
		}

		settingsPub.subscribe((settings) => {
			settingsPubClone = settings;
		});
		settingsPriv.subscribe((settings) => {
			settingsPrivClone = settings;
		});

		if (!$settingsPriv) {
			$pocketbase
				.collection('settings_private')
				.getFirstListItem('')
				.then((res) => {
					settingsPriv.set(res as SettingsPrivate);
				});
		}
	});

	function resetFavicon() {
		if (settingsPubClone === undefined) return;
		settingsPubClone.favicon = '';
		faviconInputElement.value = '';
		faviconPreview.src = '/gopher.svg';
	}

	async function save() {
		if (
			settingsPubClone === undefined ||
			settingsPrivClone === undefined ||
			faviconInputElement === undefined
		)
			return;

		if (faviconInputElement.files !== null && faviconInputElement.files?.length > 0) {
			let form = new FormData();
			form.append('favicon', faviconInputElement.files[0]);
			const res = await $pocketbase.collection('settings_public').update(settingsPubClone.id, form);
			settingsPub.set(res as SettingsPublic);
		}
		await $pocketbase
			.collection('settings_public')
			.update(settingsPubClone.id, settingsPubClone)
			.then((res) => {
				settingsPub.set(res as SettingsPublic);
			})
			.catch((err) => {
				toast.error(err.message);
			});
		await $pocketbase
			.collection('settings_private')
			.update(settingsPrivClone.id, settingsPrivClone)
			.then((res) => {
				toast.success($LL.toasts.settings_saved());
				settingsPriv.set(res as SettingsPrivate);
				goto('/');
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}
</script>

{#if settingsPubClone === undefined || settingsPrivClone === undefined}
	<PageLoading />
{:else}
	<h1 class="text-3xl font-bold mb-8">{$LL.settings.settings_title()}</h1>
	<form on:submit|preventDefault={save}>
		<div class="card w-full bg-base-300 shadow-xl mt-6">
			<div class="card-body">
				<h2 class="card-title">{$LL.settings.ping_interval_title()}</h2>
				<p class="mt-2">
					<!-- eslint-disable svelte/no-at-html-tags -->
					{@html $LL.settings.ping_interval_desc1()}
				</p>
				<p>
					<!-- Usage of  @html is not a nice way to do this, but it's safe in this case. If someones finds
						a good solution to avoid @html but still be able to use multiple links in a single sentence,
						please create a pr/issue. -->
					<!-- eslint-disable svelte/no-at-html-tags -->
					{@html $LL.settings.ping_interval_desc2()}
				</p>
				<div class="form-control w-full mt-2">
					<input
						type="text"
						placeholder="e.g. '@every 5s' or '@every 1m'"
						class="input w-full max-w-xs"
						bind:value={settingsPrivClone.interval}
					/>
				</div>
				<h2 class="card-title mt-2">{$LL.settings.lazy_ping_title()}</h2>
				<p class="mt-2">
					{$LL.settings.lazy_ping_desc()}
				</p>
				<div class="form-control w-fit">
					<label class="label cursor-pointer gap-2">
						<input
							type="checkbox"
							class="checkbox checkbox-primary"
							bind:checked={settingsPrivClone.lazy_ping}
						/>
						<span class="label-text">{$LL.settings.lazy_ping_enable()}</span>
					</label>
				</div>
			</div>
		</div>
		<div class="card w-full bg-base-300 shadow-xl mt-6">
			<div class="card-body">
				<h2 class="card-title">{$LL.settings.website_title_title()}</h2>
				<p class="my-2">{$LL.settings.website_title_desc()}</p>
				<div class="form-control w-full">
					<input
						type="text"
						placeholder="e.g. 'UpSnap'"
						class="input w-full max-w-xs"
						bind:value={settingsPubClone.website_title}
					/>
				</div>
			</div>
		</div>
		<div class="card w-full bg-base-300 shadow-xl mt-6">
			<div class="card-body">
				<h2 class="card-title">{$LL.settings.icon_title()}</h2>
				<p class="my-2">
					{$LL.settings.icon_desc()}
					<span class="badge">.ico</span> <span class="badge">.png</span>
					<span class="badge">.svg</span> <span class="badge">.gif</span>
					<span class="badge">.jpg/.jpeg</span>
				</p>
				<div>
					<img
						class="h-36"
						src={$settingsPub.favicon !== ''
							? `${backendUrl}api/files/settings_public/${$settingsPub.id}/${$settingsPub.favicon}`
							: '/gopher.svg'}
						alt="Favicon preview"
						bind:this={faviconPreview}
					/>
				</div>
				<div class="form-control w-full max-w-md flex flex-row gap-4">
					<input
						type="file"
						class="file-input w-full max-w-xs"
						accept=".ico,.png,.svg,.gif,.jpg,.jpeg"
						bind:this={faviconInputElement}
					/>
					<button
						class="btn btn-outline btn-error"
						on:click={() => resetFavicon()}
						on:keydown={() => resetFavicon()}>{$LL.buttons.reset()}</button
					>
				</div>
			</div>
		</div>
		<div class="card-actions justify-end mt-6">
			<button class="btn btn-success" type="submit"><Fa icon={faSave} />{$LL.buttons.save()}</button
			>
		</div>
	</form>
	<div class="container mx-auto text-center mt-6">
		{#if PUBLIC_VERSION === ''}
			{$LL.settings.upsnap_version()}: (untracked)
		{:else}
			{$LL.settings.upsnap_version()}:
			<a href="https://github.com/seriousm4x/UpSnap/releases/tag/{PUBLIC_VERSION}" class="link"
				>{PUBLIC_VERSION}</a
			>
		{/if}
	</div>
{/if}
