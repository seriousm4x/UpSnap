<script>
    import { devices } from '@stores/pocketbase';
    import DeviceForm from '@components/DeviceForm.svelte';
    import UnauthorizedMsg from '@components/UnauthorizedMsg.svelte';

    export let data;
    let device;

    $: {
        device = $devices[data.params.id];
        if (device && !device?.expand?.ports) {
            device.expand.ports = [];
        }
    }
</script>

{#if device}
    <div class="container">
        <UnauthorizedMsg />
        <DeviceForm bind:device mode="edit" />
    </div>
{/if}
