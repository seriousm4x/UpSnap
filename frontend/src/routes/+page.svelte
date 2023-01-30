<script>
    import { onMount } from 'svelte';
    import DeviceCard from '@components/DeviceCard.svelte';
    import { pocketbase } from '@stores/pocketbase';

    let devices = {};

    onMount(async () => {
        // get all devices in pocketbase
        const result = await $pocketbase.collection('devices').getFullList(200, {
            expand: 'ports',
            sort: 'name'
        });
        result.forEach((device) => {
            devices[device.id] = device;
        });

        // subscribe to database events
        $pocketbase.collection('devices').subscribe('*', async (e) => {
            if (e.action === 'create') {
                devices[e.record.id] = e.record;
            } else if (e.action === 'update') {
                const device = await $pocketbase.collection('devices').getOne(e.record.id, {
                    expand: 'ports'
                });
                devices[device.id] = device;
            } else if (e.action === 'delete') {
                delete devices[e.record.id];
            }
        });
        $pocketbase.collection('ports').subscribe('*', async (e) => {
            if (e.action === 'update') {
                const device = await $pocketbase
                    .collection('devices')
                    .getFirstListItem(`ports.id ?= "${e.record.id}"`, {
                        expand: 'ports'
                    });
                devices[device.id] = device;
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

<div class="container text-body-emphasis">
    {#if Object.keys(devices).length > 0}
        <div class="row">
            {#each Object.entries(devices) as [_, device]}
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
