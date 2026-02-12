<script lang="ts">
	import { fly } from 'svelte/transition';
	let { url, children } = $props();
</script>

{#key url}
	<div in:fly={{ y: 50, duration: 300 }}>
		{@render children()}
	</div>
{/key}
