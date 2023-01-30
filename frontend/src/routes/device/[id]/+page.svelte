<script>
    import { devices } from '@stores/pocketbase';
    import DeviceForm from '@components/DeviceForm.svelte';

    export let data;
    let device = $devices[data.params.id];
    if (!device?.expand?.ports) {
        device.expand.ports = [];
    }
</script>

{#if device}
    <div class="container">
        <DeviceForm bind:device mode="edit" />
    </div>
{/if}
