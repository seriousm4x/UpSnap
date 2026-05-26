<script lang="ts">
	import { goto } from '$app/navigation';
	import { asset, resolve } from '$app/paths';
	import { m } from '$lib/paraglide/messages';
	import { pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import type { SettingsPublic } from '$lib/types/settings';

	let notReady = $state(false);

	async function checkAndContinue() {
		const data = await $pocketbase.collection('settings_public').getFirstListItem('');
		settingsPub.set(data as SettingsPublic);
		if (data.setup_completed) {
			goto(resolve('/login'));
		} else {
			notReady = true;
		}
	}
</script>

<div class="mt-10 flex items-center justify-center">
	<div class="my-4 flex w-screen max-w-lg flex-col gap-16">
		<div class="card bg-base-200 shadow-sm">
			{#if !$settingsPub?.setup_completed}
				<figure class="mx-auto w-44 pt-6"><img src={asset('/gopher.svg')} alt="Gopher" /></figure>
				<div class="card-body">
					<h2 class="card-title">{m.welcome_title()}</h2>
					<p>{m.welcome_desc()}</p>
					{#if notReady}
						<p class="text-error text-sm">{m.welcome_not_ready()}</p>
					{/if}
					<div class="card-actions justify-end">
						<button class="btn btn-success" onclick={checkAndContinue}
							>{m.welcome_btn_done()}</button
						>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
