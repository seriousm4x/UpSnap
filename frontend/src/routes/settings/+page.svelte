<script lang="ts">
	import { settingsPub, settingsPriv } from '$lib/stores/settings';
	import { pocketbase, backendUrl } from '$lib/stores/pocketbase';
	import type { SettingsPublic, SettingsPrivate } from '$lib/types/settings';

	let settingsPubClone: SettingsPublic | undefined;
	let settingsPrivClone: SettingsPrivate | undefined;
	let faviconPreview: HTMLImageElement;
	let faviconInputElement: HTMLInputElement;
	let saveBtn: HTMLButtonElement;
	let version = import.meta.env.UPSNAP_VERSION;

	settingsPub.subscribe((settings) => {
		settingsPubClone = settings;
	});
	settingsPriv.subscribe((settings) => {
		settingsPrivClone = settings;
	});

	if (!$settingsPriv && $pocketbase.authStore.isValid) {
		$pocketbase
			.collection('settings_private')
			.getList(1, 1)
			.then((res) => {
				settingsPriv.set(res.items[0] as SettingsPrivate);
			});
	}

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
			faviconInputElement === undefined ||
			saveBtn === undefined
		)
			return;

		saveBtn.classList.remove('btn-primary');
		saveBtn.classList.add('btn-warning');

		if (faviconInputElement.files !== null && faviconInputElement.files?.length > 0) {
			let form = new FormData();
			form.append('favicon', faviconInputElement.files[0]);
			await $pocketbase.collection('settings_public').update(settingsPubClone.id, form);
		}
		const res = await $pocketbase
			.collection('settings_public')
			.update(settingsPubClone.id, settingsPubClone);
		settingsPub.set(res as SettingsPublic);

		saveBtn.classList.remove('btn-warning');
		saveBtn.classList.add('btn-primary');
	}
</script>

{#if settingsPubClone === undefined || settingsPrivClone === undefined}
	<div class="container mx-auto text-center"><span class="loading loading-ring loading-md" /></div>
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
				<p class="mb-2">
					Learn more about the correct syntax for cron on
					<a target="_blank" class="link" href="https://en.wikipedia.org/wiki/Cron">Wikipedia</a>
					or refer to the
					<a target="_blank" class="link" href="https://pkg.go.dev/github.com/robfig/cron/v3"
						>package documentation</a
					>.
				</p>
				<div class="form-control w-full">
					<input
						type="text"
						placeholder="e.g. '@every 5s' or '@every 1m'"
						class="input w-full max-w-xs"
						bind:value={settingsPrivClone.interval}
					/>
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
		<div class="card-actions justify-end mt-4">
			<button class="btn btn-primary" type="submit" bind:this={saveBtn}>Save</button>
		</div>
	</form>

	<div class="container mx-auto text-center">
		{#if version === undefined}
			UpSnap version: (untracked)
		{:else}
			UpSnap version: <a
				href="https://github.com/seriousm4x/UpSnap/releases/tag/{version}"
				class="text-reset">{version}</a
			>
		{/if}
	</div>
{/if}
