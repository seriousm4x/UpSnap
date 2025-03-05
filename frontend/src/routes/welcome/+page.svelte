<script lang="ts">
	import { goto } from '$app/navigation';
	import { toggleVisibility } from '$lib/helpers/forms';
	import { m } from '$lib/paraglide/messages.js';
	import { pocketbase } from '$lib/stores/pocketbase';
	import { settingsPub } from '$lib/stores/settings';
	import type { SettingsPublic } from '$lib/types/settings';
	import { faArrowRight, faEye } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	let stepsCompleted = 0;
	let inputPassword: HTMLInputElement;
	let inputConfirm: HTMLInputElement;
	let form = {
		email: '',
		password: '',
		confirm: ''
	};

	onMount(() => {
		if ($settingsPub && !$settingsPub.setup_completed && $pocketbase.authStore.isValid) {
			$pocketbase.authStore.clear();
		}
	});

	async function register() {
		$pocketbase
			.send('/api/upsnap/init-superuser', {
				method: 'POST',
				body: {
					email: form.email,
					password: form.password,
					password_confirm: form.confirm
				}
			})
			.then(() => {
				$pocketbase
					.collection('settings_public')
					.getFirstListItem('')
					.then((data) => {
						settingsPub.set(data as SettingsPublic);
					});
				$pocketbase
					.collection('_superusers')
					.authWithPassword(form.email, form.password)
					.then(() => {
						stepsCompleted = 2;
					})
					.catch((err) => {
						if (err.data?.data?.identity?.message) {
							toast.error(err.data.data.identity.message);
							return;
						}
					});
			})
			.catch((err) => {
				if (err.data?.data?.passwordConfirm?.message) {
					toast.error(m.toasts_passwords_missmatch());
				} else if (err.data?.data?.email?.message) {
					toast.error(err.data.data.email.message);
				} else {
					toast.error(err.message || err);
				}
			});
	}
</script>

<div class="mt-10 flex items-center justify-center">
	<div class="my-4 flex w-screen max-w-lg flex-col gap-16">
		<div class="card bg-base-200 shadow-sm">
			{#if stepsCompleted === 0 && $settingsPub?.setup_completed}
				<figure class="mx-auto w-72 pt-6"><img src="/gopher.svg" alt="Gopher" /></figure>
				<div class="card-body">
					<h2 class="card-title">{m.welcome_not_expected_title()}</h2>
					<p>{m.welcome_not_expected_desc()}</p>
					<div class="card-actions justify-end">
						<button class="btn btn-success" on:click={() => goto('/')}
							>{m.welcome_not_expected_back()}</button
						>
					</div>
				</div>
			{:else if stepsCompleted === 0}
				<figure class="mx-auto w-44 pt-6"><img src="/gopher.svg" alt="Gopher" /></figure>
				<div class="card-body">
					<h2 class="card-title">{m.welcome_step1_page_title()}</h2>
					<p>{m.welcome_step1_setup_desc()}</p>
					<div class="card-actions justify-end">
						<button class="btn btn-success" on:click={() => (stepsCompleted = 1)}
							>{m.welcome_step1_setup_btn_next()} <Fa icon={faArrowRight} /></button
						>
					</div>
				</div>
			{:else if stepsCompleted === 1}
				<div class="card-body">
					<div class="flex flex-row gap-4">
						<figure class="w-16"><img src="/gopher.svg" alt="Gopher" /></figure>
						<h2 class="card-title">{m.welcome_step2_page_title()}</h2>
					</div>
					<form class="w-full" on:submit|preventDefault={register}>
						<label class="label" for="email">
							<span>{m.welcome_step2_label_email()}</span>
						</label>
						<input id="email" type="email" class="input w-full" required bind:value={form.email} />
						<label class="label" for="password">
							<span>{m.welcome_step2_label_password()}</span>
							<span>{m.welcome_step2_label_min_chars()}</span>
						</label>
						<label class="relative block">
							<div
								class="absolute top-1/2 right-4 -translate-y-1/2 cursor-pointer"
								role="none"
								on:click={() => toggleVisibility(inputPassword)}
								on:keydown={() => toggleVisibility(inputPassword)}
							>
								<Fa icon={faEye} />
							</div>
							<input
								id="password"
								type="password"
								class="input w-full"
								minlength="10"
								maxlength="72"
								required
								bind:value={form.password}
								bind:this={inputPassword}
							/>
						</label>
						<label class="label" for="passwordConfirm">
							<span>{m.welcome_step2_label_password_confirm()}</span>
						</label>
						<label class="relative block">
							<div
								class="absolute top-1/2 right-4 -translate-y-1/2 cursor-pointer"
								role="none"
								on:click={() => toggleVisibility(inputConfirm)}
								on:keydown={() => toggleVisibility(inputConfirm)}
							>
								<Fa icon={faEye} />
							</div>
							<input
								id="confirm"
								type="password"
								class="input w-full"
								minlength="10"
								maxlength="72"
								required
								bind:value={form.confirm}
								bind:this={inputConfirm}
							/>
						</label>
						<div class="card-actions mt-4 justify-end">
							<button class="btn btn-success" type="submit"
								>{m.welcome_step2_btn_create()} <Fa icon={faArrowRight} /></button
							>
						</div>
					</form>
				</div>
			{:else if stepsCompleted === 2}
				<figure class="mx-auto w-72 pt-6"><img src="/gopher.svg" alt="Gopher" /></figure>
				<div class="card-body">
					<h2 class="card-title">{m.welcome_step3_page_title()}</h2>
					<p>{m.welcome_step3_page_desc()}</p>
					<div class="card-actions justify-end">
						<button class="btn btn-success" on:click={() => goto('/')}
							>{m.welcome_step3_btn_done()}</button
						>
					</div>
				</div>
			{/if}
		</div>
		{#if $settingsPub && !$settingsPub.setup_completed}
			<ul class="steps steps-horizontal">
				<li class="step step-success">{m.welcome_progress_step1()}</li>
				<li class="step" class:step-success={stepsCompleted > 0}>{m.welcome_progress_step2()}</li>
				<li class="step" class:step-success={stepsCompleted > 1}>{m.welcome_progress_step3()}</li>
			</ul>
		{/if}
	</div>
</div>
