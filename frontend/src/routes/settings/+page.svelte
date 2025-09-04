<script lang="ts">
	import { goto } from '$app/navigation';
	import { asset, resolve } from '$app/paths';
	import { page } from '$app/state';
	import { PUBLIC_VERSION } from '$env/static/public';
	import PageLoading from '$lib/components/PageLoading.svelte';
	import { nextCronDate, parseCron } from '$lib/helpers/cron';
	import { m } from '$lib/paraglide/messages';
	import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
	import { settingsPriv, settingsPub } from '$lib/stores/settings';
	import type { SettingsPrivate, SettingsPublic } from '$lib/types/settings';
	import { faSave } from '@fortawesome/free-solid-svg-icons';
	import cronParser from 'cron-parser';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let settingsPubClone: SettingsPublic | undefined;
	let settingsPrivClone: SettingsPrivate | undefined;
	let faviconPreview: HTMLImageElement;
	let faviconInputElement: HTMLInputElement;

	onMount(() => {
		if (!$pocketbase.authStore.isSuperuser) {
			toast(m.toasts_no_permission({ url: page.url.pathname }), {
				icon: '⛔'
			});
			goto(resolve('/'));
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

		try {
			cronParser.parse(settingsPrivClone.interval);
		} catch {
			toast.error(m.settings_invalid_cron());
			throw new Error(`failed to parse cron: ` + settingsPrivClone.interval);
		}

		if (settingsPrivClone.interval === '' || !(await parseCron(settingsPrivClone.interval))) {
			toast.error(m.settings_invalid_cron());
			throw new Error('ping_interval not valid');
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

		if (faviconInputElement.files !== null && faviconInputElement.files?.length > 0) {
			let form = new FormData();
			form.append('favicon', faviconInputElement.files[0]);
			const res = await $pocketbase.collection('settings_public').update(settingsPubClone.id, form);
			settingsPub.set(res as SettingsPublic);
		}

		await $pocketbase
			.collection('settings_private')
			.update(settingsPrivClone.id, settingsPrivClone)
			.then((res) => {
				toast.success(m.toasts_settings_saved());
				settingsPriv.set(res as SettingsPrivate);
				goto(resolve('/'));
			})
			.catch((err) => {
				toast.error(err.message);
			});
	}
</script>

{#if settingsPubClone === undefined || settingsPrivClone === undefined}
	<PageLoading />
{:else}
	<h1 class="mb-8 text-3xl font-bold">{m.settings_page_title()}</h1>
	<form
		onsubmit={(e) => {
			e.preventDefault();
			save();
		}}
		class="flex w-full flex-col"
	>
		<div class="card bg-base-200 mt-6 w-full shadow-sm">
			<div class="card-body">
				<h2 class="card-title">{m.settings_ping_interval_title()}</h2>
				<p class="mt-2">
					<!-- eslint-disable svelte/no-at-html-tags -->
					{@html m.settings_ping_interval_desc1()}
				</p>
				<p>
					<!-- Usage of  @html is not a nice way to do this, but it's safe in this case. If someones finds
						a good solution to avoid @html but still be able to use multiple links in a single sentence,
						please create a pr/issue. -->
					<!-- eslint-disable svelte/no-at-html-tags -->
					{@html m.settings_ping_interval_desc2()}
				</p>
				<div class="mt-2 w-full">
					<input
						type="text"
						placeholder="e.g. '*/5 * * * * *' or '* */1 * * * *'"
						class="input"
						bind:value={settingsPrivClone.interval}
					/>
				</div>
				<p class="fieldset-label">
					{#await parseCron(settingsPrivClone.interval)}
						<span class="loading loading-spinner loading-xs"></span>
					{:then valid}
						{valid ? '✅ ' + nextCronDate(settingsPrivClone.interval) : m.settings_invalid_cron()}
					{/await}
				</p>
				<pre class="bg-base-100 text-base-content/80 w-fit rounded p-2"><code
						>* * * * * *
| | | | | |
| | | | | day of the week (0–6)
| | | | month (1–12)
| | | day of the month (1–31)
| | hour (0–23)
| minute (0–59)
second (0–59, optional)
</code></pre>
				<h2 class="card-title mt-2">{m.settings_lazy_ping_title()}</h2>
				<p class="mt-2">
					{m.settings_lazy_ping_desc()}
				</p>
				<div class="w-fit">
					<label class="label cursor-pointer gap-2">
						<input type="checkbox" class="checkbox" bind:checked={settingsPrivClone.lazy_ping} />
						<span>{m.settings_lazy_ping_enable()}</span>
					</label>
				</div>
			</div>
		</div>
		<div class="card bg-base-200 mt-6 w-full shadow-sm">
			<div class="card-body">
				<h2 class="card-title">{m.settings_website_title_title()}</h2>
				<p class="my-2">{m.settings_website_title_desc()}</p>
				<div class="w-full">
					<input
						type="text"
						placeholder="e.g. 'UpSnap'"
						class="input"
						bind:value={settingsPubClone.website_title}
					/>
				</div>
			</div>
		</div>
		<div class="card bg-base-200 mt-6 w-full shadow-sm">
			<div class="card-body">
				<h2 class="card-title">{m.settings_icon_title()}</h2>
				<p class="my-2">
					{m.settings_icon_desc()}
					<span class="badge">.ico</span> <span class="badge">.png</span>
					<span class="badge">.svg</span> <span class="badge">.gif</span>
					<span class="badge">.jpg/.jpeg</span>
				</p>
				<div>
					<img
						class="h-36"
						src={$settingsPub.favicon !== ''
							? `${backendUrl}api/files/settings_public/${$settingsPub.id}/${$settingsPub.favicon}`
							: asset('/gopher.svg')}
						alt="Favicon preview"
						bind:this={faviconPreview}
					/>
				</div>
				<div class="flex w-full max-w-md flex-row gap-4">
					<input
						type="file"
						class="file-input"
						accept=".ico,.png,.svg,.gif,.jpg,.jpeg"
						bind:this={faviconInputElement}
					/>
					<button
						class="btn btn-error"
						onclick={() => resetFavicon()}
						onkeydown={() => resetFavicon()}>{m.buttons_reset()}</button
					>
				</div>
			</div>
		</div>
		<div class="card-actions mt-6 justify-end">
			<button class="btn btn-success" type="submit"><Fa icon={faSave} />{m.buttons_save()}</button>
		</div>
	</form>
	<div class="container mx-auto mt-6 text-center">
		{#if PUBLIC_VERSION === ''}
			{m.settings_upsnap_version()}: (untracked)
		{:else}
			{m.settings_upsnap_version()}:
			<a href="https://github.com/seriousm4x/UpSnap/releases/tag/{PUBLIC_VERSION}" class="link"
				>{PUBLIC_VERSION}</a
			>
		{/if}
	</div>
{/if}
