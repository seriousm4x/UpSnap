<script>
    import { onMount } from 'svelte';
    import DeviceCard from '@components/DeviceCard.svelte';
    import { pocketbase } from '@stores/pocketbase';

    let devices = {};

    onMount(async () => {
        let pb;
        pocketbase.subscribe((conn) => {
            pb = conn;
        });

        // get all devices in pocketbase
        const result = await pb.collection('devices').getFullList(200, {
            expand: 'ports',
            sort: 'name'
        });
        result.forEach((device) => {
            devices[device.id] = device;
        });

        // subscribe to database events
        pb.collection('devices').subscribe('*', async (e) => {
            if (e.action === 'create') {
                devices[e.record.id] = e.record;
            } else if (e.action === 'update') {
                const device = await pb.collection('devices').getOne(e.record.id, {
                    expand: 'ports'
                });
                devices[device.id] = device;
            } else if (e.action === 'delete') {
                delete devices[e.record.id];
            }
        });
        pb.collection('ports').subscribe('*', async (e) => {
            if (e.action === 'update') {
                const device = await pb
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
    let clear;
    $: {
        clearInterval(clear);
        clear = setInterval(() => {
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
