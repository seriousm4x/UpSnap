<script lang="ts">
	import { scale } from 'svelte/transition';
	import { tweened } from 'svelte/motion';

	import Fa from 'svelte-fa';
	import type { IconDefinition } from '@fortawesome/free-solid-svg-icons';

	export let color: string;
	export let message: string;
	export let icon: IconDefinition;
	export let customClasses = '';
	export let duration = 0;

	const value = tweened(0, {
		duration: duration
	});
	value.set(duration);
</script>

<div class="alert alert-{color} {customClasses}" transition:scale={{ delay: 0, duration: 200 }}>
	<div class="flex flex-col text-left">
		{#if icon}
			<Fa {icon} />
		{/if}
		<span>{message}</span>
		{#if duration !== 0}
			<progress class="progress" value={duration - $value} max={duration} />
		{/if}
	</div>
</div>
