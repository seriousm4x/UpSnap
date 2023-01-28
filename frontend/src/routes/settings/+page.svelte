<script>
    import { onMount } from 'svelte';
    import { pocketbase } from '@stores/pocketbase';
    import DeviceForm from '@components/DeviceForm.svelte';

    let pb;
    let files;
    let settings = {};
    let timeout;
    let buttons = {
        settings: {
            state: 'none',
            error: ''
        },
        restore: {
            state: 'none',
            error: ''
        }
    };
    let newDevice = {
        name: '',
        ip: '',
        mac: '',
        netmask: '255.255.255.0',
        expand: {
            ports: []
        },
        ports: [],
        link: '',
        wake_cron: '',
        wake_cron_enabled: false,
        shutdown_cron: '',
        shutdown_cron_enabled: false,
        shutdown_cmd: '',
        password: ''
    };

    onMount(async () => {
        pocketbase.subscribe((conn) => {
            pb = conn;
        });

        const result = await pb.collection('settings').getList(1, 1);
        settings = result.items[0];
    });

    async function saveSettings() {
        buttons.settings.state = 'waiting';
        try {
            if (settings.interval === '') {
                settings.interval = '@every 3s';
            }
            await pb.collection('settings').update(settings.id, settings);
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                buttons.settings.state = 'none';
            }, 3000);
            buttons.settings.state = 'success';
        } catch (error) {
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                buttons.settings.error = '';
                buttons.settings.state = 'none';
            }, 3000);
            buttons.settings.error = error;
            buttons.settings.state = 'failed';
        }
    }

    async function restoreV2Backup() {
        buttons.restore.state = 'waiting';

        try {
            // delete all devices in pocketbase
            let result = await pb.collection('devices').getFullList();
            for (let index = 0; index < result.length; index++) {
                pb.collection('devices').delete(result[index].id);
            }
            // delete all ports in pocketbase
            result = await pb.collection('ports').getFullList();
            for (let index = 0; index < result.length; index++) {
                await pb.collection('ports').delete(result[index].id);
            }

            let reader = new FileReader();
            reader.readAsText(files[0]);
            reader.onload = async (e) => {
                // parse uploaded file
                let data = JSON.parse(e.target.result);
                if (!Array.isArray(data)) {
                    return;
                }

                // loop devices
                for (let index = 0; index < data.length; index++) {
                    const device = data[index];

                    // create ports
                    let thisDevicePorts = [];
                    for (let index = 0; index < device.ports.length; index++) {
                        const port = device.ports[index];
                        const record = await pb.collection('ports').create({
                            name: port.name,
                            number: port.number
                        });
                        thisDevicePorts.push(record.id);
                    }

                    // create device
                    await pb.collection('devices').create({
                        name: device.name,
                        ip: device.ip,
                        mac: device.mac,
                        netmask: device.netmask,
                        ports: thisDevicePorts,
                        link: device.link,
                        wake_cron: device.wake.cron,
                        wake_cron_enabled: device.wake.enabled,
                        shutdown_cron: device.shutdown.cron,
                        shutdown_cron_enabled: device.shutdown.cron.enabled,
                        shutdown_cmd: device.shutdown.command
                    });
                }
            };
            setTimeout(() => {
                buttons.restore.state = 'none';
            }, 3000);
            buttons.restore.state = 'success';
        } catch (error) {
            setTimeout(() => {
                buttons.restore.error = '';
                buttons.restore.state = 'none';
            }, 3000);
            buttons.restore.error = error;
            buttons.restore.state = 'failed';
        }
    }
</script>

<div class="container">
    <DeviceForm bind:device={newDevice} mode="add" />
    <section class="m-0 mt-4 p-4 shadow-sm">
        <form on:submit|preventDefault={saveSettings}>
            <h3 class="mb-3">Ping interval</h3>
            <div class="row">
                <div class="col-md-6">
                    <p>Sets the interval in which the devices are pinged.</p>
                    <div class="input-group mb-3">
                        <span class="input-group-text">Cron</span>
                        <input
                            class="form-control"
                            placeholder="e.g. '@every 5s' or '@every 1m'"
                            aria-label="Interval"
                            aria-describedby="addon-wrapping"
                            type="text"
                            bind:value={settings.interval}
                        />
                    </div>
                    <h3 class="my-3">Website title</h3>
                    <p>Set the website title in the navbar.</p>
                    <div class="input-group mb-3">
                        <span class="input-group-text" id="website-title">Title</span>
                        <input
                            type="text"
                            class="form-control"
                            placeholder="e.g. 'UpSnap'"
                            aria-label="Website title"
                            aria-describedby="website-title"
                            bind:value={settings.website_title}
                        />
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="callout callout-info m-0">
                        <p class="m-0">
                            Leave blank to use default value of <span class="fw-bold"
                                >"@every 3s"</span
                            >.
                        </p>
                        <p class="m-0">
                            Read more about valid cron syntax at <a
                                target="_blank"
                                rel="noreferrer"
                                href="https://pkg.go.dev/github.com/robfig/cron"
                                >pkg.go.dev/github.com/robfig/cron</a
                            >
                        </p>
                    </div>
                </div>
            </div>
            <button
                type="submit"
                class="btn btn-secondary"
                class:btn-success={buttons.settings.state === 'success' ? true : false}
                class:btn-warning={buttons.settings.state === 'waiting' ? true : false}
                class:btn-danger={buttons.settings.state === 'failed' ? true : false}
                disabled={buttons.settings.state !== 'none' ? true : false}
            >
                {#if buttons.settings.state === 'none'}
                    Save
                {:else if buttons.settings.state === 'success'}
                    Saved
                {:else if buttons.settings.state === 'waiting'}
                    Waiting
                {:else if buttons.settings.state === 'failed'}
                    Failed: {buttons.settings.error}
                {/if}
            </button>
        </form>
    </section>
    <section class="m-0 my-4 p-4 shadow-sm">
        <h3 class="mb-3">Restore</h3>
        <p>
            If you had UpSnap v2.3.2 (or higher) running before, you can restore your devices here.
        </p>
        <div class="callout callout-danger fw-bold">This will wipe the existing database!</div>
        <input
            class="form-control mb-3"
            aria-label="Restore"
            aria-describedby="Restore"
            type="file"
            accept=".json"
            bind:files
        />
        <button
            type="button"
            class="btn btn-secondary"
            class:btn-success={buttons.restore.state === 'success' ? true : false}
            class:btn-warning={buttons.restore.state === 'waiting' ? true : false}
            class:btn-danger={buttons.restore.state === 'failed' ? true : false}
            disabled={files ? false : true}
            on:click={() => restoreV2Backup()}
        >
            {#if buttons.restore.state === 'none'}
                Restore
            {:else if buttons.restore.state === 'success'}
                Restored
            {:else if buttons.restore.state === 'waiting'}
                Waiting
            {:else if buttons.restore.state === 'failed'}
                Failed: {buttons.restore.error}
            {/if}
        </button>
    </section>
</div>
