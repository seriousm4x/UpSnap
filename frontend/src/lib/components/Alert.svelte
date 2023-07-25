<script lang="ts">
	import { scale } from 'svelte/transition';

	import Fa from 'svelte-fa';
	import type { IconDefinition } from '@fortawesome/free-solid-svg-icons';

	export let color: string;
	export let message: string;
	export let icon: IconDefinition;
	export let customClasses = '';
</script>

<div class="alert alert-{color} {customClasses}" transition:scale={{ delay: 0, duration: 200 }}>
	{#if icon}
		<Fa {icon} />
	{/if}
	<span>{message}</span>
</div>
