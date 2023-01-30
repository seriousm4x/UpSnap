<script>
    import { onMount } from 'svelte';
    import DeviceCard from '@components/DeviceCard.svelte';
    import { pocketbase, devices } from '@stores/pocketbase';
    import { sortDevices } from '../sorts';

    onMount(async () => {
        // subscribe to database events
        $pocketbase.collection('devices').subscribe('*', async (e) => {
            if (e.action === 'create') {
                $devices[e.record.id] = e.record;
            } else if (e.action === 'update') {
                const device = await $pocketbase.collection('devices').getOne(e.record.id, {
                    expand: 'ports'
                });
                $devices[device.id] = device;
            } else if (e.action === 'delete') {
                delete $devices[e.record.id];
            }
        });
        $pocketbase.collection('ports').subscribe('*', async (e) => {
            if (e.action === 'update') {
                const device = Object.values($devices).find((dev) =>
                    dev.ports.includes(e.record.id)
                );
                if (!device) {
                    return;
                }

                // replace device.expand.ports with updated record
                const portIdx = device.expand.ports.findIndex((port) => port.id === e.record.id);
                $devices[device.id].expand.ports[portIdx] = e.record;
            }
        });
    });

    // update device date
    let now = Date.now();
    let interval;
    $: {
        clearInterval(interval);
        interval = setInterval(() => {
            now = Date.now();
        }, 1000);
    }
</script>

<div class="container text-body-emphasis mb-4">
    {#if Object.keys($devices).length > 0}
        <div class="row">
            {#each Object.values($devices).sort(sortDevices) as device}
                <DeviceCard {device} {now} />
            {/each}
        </div>
    {:else}
        <div class="text-center">
            <h4 class="text-muted">No devices</h4>
            <p>
                <a class="text-muted" href="/settings">Go to settings to add devices...</a>
            </p>
        </div>
    {/if}
</div>
