<script>
    import { onMount } from 'svelte';
    import DeviceCard from '@components/DeviceCard.svelte';
    import { pocketbase, devices } from '@stores/pocketbase';
    import { sortDevices, sortGroups } from '../sorts';

    onMount(async () => {
        // subscribe to database events
        $pocketbase.collection('devices').subscribe('*', async (e) => {
            if (e.action === 'create') {
                $devices[e.record.id] = e.record;
            } else if (e.action === 'update') {
                const device = await $pocketbase.collection('devices').getOne(e.record.id, {
                    expand: 'ports,groups'
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

    let devicesWithoutGroups = [];
    let devicesWithGroup = {};

    devices.subscribe((d) => {
        // sort devices into their groups
        devicesWithoutGroups = [];
        devicesWithGroup = {};
        Object.values(d).forEach((device) => {
            if (device.groups.length === 0) {
                devicesWithoutGroups = [...devicesWithoutGroups, device];
                return;
            }
            device.expand.groups.forEach((grp) => {
                if (!devicesWithGroup.hasOwnProperty(grp.name)) {
                    devicesWithGroup[grp.name] = [];
                }
                devicesWithGroup[grp.name] = [...devicesWithGroup[grp.name], device];
            });
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
        {#if devicesWithoutGroups.length > 0}
            <div class="row">
                {#each devicesWithoutGroups.sort(sortDevices) as device}
                    <DeviceCard {device} {now} />
                {/each}
            </div>
        {/if}

        {#each Object.keys(devicesWithGroup).sort(sortGroups) as grp}
            <h2 class="mt-4 mb-0">{grp}</h2>
            <div class="row">
                {#each devicesWithGroup[grp] as device}
                    <DeviceCard {device} {now} />
                {/each}
            </div>
        {/each}
    {:else}
        <div class="text-center">
            <h4 class="text-muted">No devices</h4>
            <p>
                <a class="text-muted" href="/settings">Go to settings to add devices...</a>
            </p>
        </div>
    {/if}
</div>
