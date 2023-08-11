<script lang="ts">
	import { pocketbase } from '$lib/stores/pocketbase';
</script>

{$pocketbase.authStore.model?.email || $pocketbase.authStore.model?.username}
