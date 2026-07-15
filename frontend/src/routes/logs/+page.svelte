<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages';
	import { backendUrl, pocketbase } from '$lib/stores/pocketbase';
	import { faArrowsRotate } from '@fortawesome/free-solid-svg-icons';
	import { onDestroy, onMount } from 'svelte';
	import Fa from 'svelte-fa';
	import toast from 'svelte-french-toast';

	type LogEntry = {
		time: string;
		level: string;
		message: string;
	};

	let logs = $state([] as LogEntry[]);
	let loading = $state(true);
	let interval: ReturnType<typeof setInterval>;

	const levelBadge: Record<string, string> = {
		INFO: 'badge-info',
		DEBUG: 'badge-neutral',
		WARNING: 'badge-warning',
		ERROR: 'badge-error'
	};

	function getLogs() {
		fetch(`${backendUrl}api/upsnap/logs`, {
			headers: {
				Authorization: $pocketbase.authStore.token
			}
		})
			.then(async (resp) => {
				if (resp.ok) {
					return resp.json();
				}
				return Promise.reject(await resp.json());
			})
			.then((data) => {
				logs = (data as LogEntry[]).slice().reverse();
			})
			.catch((err) => {
				toast.error(err.message ?? String(err));
			})
			.finally(() => (loading = false));
	}

	onMount(() => {
		if (!$pocketbase.authStore.isSuperuser) {
			toast(m.toasts_no_permission({ url: page.url.pathname }), {
				icon: '⛔'
			});
			goto(resolve('/'));
			return;
		}

		getLogs();
		interval = setInterval(getLogs, 5000);
	});

	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<svelte:head>
	<title>{m.logs_page_title()} | UpSnap</title>
</svelte:head>

<div class="container mx-auto p-4">
	<div class="mb-4 flex items-center justify-between">
		<h1 class="text-2xl font-bold">{m.logs_page_title()}</h1>
		<button class="btn btn-ghost btn-sm" onclick={getLogs}>
			<Fa icon={faArrowsRotate} />
			{m.logs_refresh()}
		</button>
	</div>

	{#if loading}
		<span class="loading loading-spinner loading-lg"></span>
	{:else if logs.length === 0}
		<p>{m.logs_empty()}</p>
	{:else}
		<div class="overflow-x-auto">
			<table class="table-zebra table table-sm">
				<thead>
					<tr>
						<th class="w-48">{m.logs_column_time()}</th>
						<th class="w-24">{m.logs_column_level()}</th>
						<th>{m.logs_column_message()}</th>
					</tr>
				</thead>
				<tbody>
					{#each logs as entry, i (i)}
						<tr>
							<td class="whitespace-nowrap">{new Date(entry.time).toLocaleString()}</td>
							<td>
								<span class="badge {levelBadge[entry.level] ?? 'badge-neutral'}">{entry.level}</span
								>
							</td>
							<td class="font-mono text-xs break-all">{entry.message}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
