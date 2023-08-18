<script lang="ts">
	import { goto } from '$app/navigation';
	import { settingsPub, settingsPriv } from '$lib/stores/settings';
	import { pocketbase, backendUrl, isAdmin } from '$lib/stores/pocketbase';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import toast from 'svelte-french-toast';
	import Fa from 'svelte-fa';
	import { faSave } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { PUBLIC_VERSION } from '$env/static/public';
	import type { SettingsPublic, SettingsPrivate } from '$lib/types/settings';

	let settingsPubClone: SettingsPublic | undefined;
	let settingsPrivClone: SettingsPrivate | undefined;
	let faviconPreview: HTMLImageElement;
	let faviconInputElement: HTMLInputElement;

	onMount(() => {
		if (!$isAdmin) {
			toast(`You don't have permission to visit ${$page.url.pathname}`, {
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
				toast.success('Saved settings');
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
	<h1 class="text-3xl font-bold mb-8">Settings</h1>
	<form on:submit|preventDefault={save}>
		<div class="card w-full bg-base-300 shadow-xl mt-6">
			<div class="card-body">
				<h2 class="card-title">Ping interval</h2>
				<p class="mt-2">
					Sets the interval in which the devices are pinged. Leave blank to use default value of <span
						class="badge">@every 3s</span
					>.
				</p>
				<p>
					Learn more about the correct syntax for cron on
					<a target="_blank" class="link" href="https://en.wikipedia.org/wiki/Cron">Wikipedia</a>
					or refer to the
					<a target="_blank" class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3"
						>package documentation</a
					>.
				</p>
				<div class="form-control w-full mt-2">
					<input
						type="text"
						placeholder="e.g. '@every 5s' or '@every 1m'"
						class="input w-full max-w-xs"
						bind:value={settingsPrivClone.interval}
					/>
				</div>
				<h2 class="card-title mt-2">Lazy ping</h2>
				<p class="mt-2">
					When lazy ping is turned on, UpSnap will only ping devices if there is an active user
					visiting the website. If it's turned off, UpSnap will always ping devices.
				</p>
				<div class="form-control w-fit">
					<label class="label cursor-pointer gap-2">
						<input
							type="checkbox"
							class="checkbox checkbox-primary"
							bind:checked={settingsPrivClone.lazy_ping}
						/>
						<span class="label-text">Enable</span>
					</label>
				</div>
			</div>
		</div>
		<div class="card w-full bg-base-300 shadow-xl mt-6">
			<div class="card-body">
				<h2 class="card-title">Website title</h2>
				<p class="my-2">Sets the title of the website and in the browser tab.</p>
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
				<h2 class="card-title">Icon</h2>
				<p class="my-2">
					Set a custom favicon. Supported file types are
					<span class="badge">.ico</span>, <span class="badge">.png</span>,
					<span class="badge">.svg</span>, <span class="badge">.gif</span> and
					<span class="badge">.jpg/.jpeg</span>.
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
						on:keydown={() => resetFavicon()}>Reset</button
					>
				</div>
			</div>
		</div>
		<div class="card-actions justify-end mt-6">
			<button class="btn btn-success" type="submit"><Fa icon={faSave} />Save</button>
		</div>
	</form>
	<div class="container mx-auto text-center mt-6">
		{#if PUBLIC_VERSION === ''}
			UpSnap version: (untracked)
		{:else}
			UpSnap version: <a
				href="https://github.com/seriousm4x/UpSnap/releases/tag/{PUBLIC_VERSION}"
				class="link">{PUBLIC_VERSION}</a
			>
		{/if}
	</div>
{/if}
