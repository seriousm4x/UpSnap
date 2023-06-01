<script>
    import { dev } from '$app/environment';
    import { pocketbase, settings_private, settings_public, devices } from '@stores/pocketbase';
    import UnauthorizedMsg from '@components/UnauthorizedMsg.svelte';
    import DeviceForm from '@components/DeviceForm.svelte';
    import Fa from 'svelte-fa/src/fa.svelte';
    import { faPlus } from '@fortawesome/free-solid-svg-icons';

    let version = import.meta.env.UPSNAP_VERSION;

    let iconFiles = [];
    let iconPreview;
    let timeout;
    let buttons = {
        settings: {
            state: 'none',
            error: ''
        },
        scan: {
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
    let scannedDevices = {};

    $: if (iconPreview && iconFiles.length > 0) {
        iconPreview.src = URL.createObjectURL(iconFiles[0]);
    } else if (iconPreview && $settings_public.favicon !== '') {
        iconPreview.src = `${dev ? $pocketbase.baseUrl : ''}/api/files/settings_public/${
            $settings_public.id
        }/${$settings_public.favicon}`;
    } else if (iconPreview) {
        iconPreview.src = '/gopher.svg';
    }

    async function saveSettings() {
        buttons.settings.state = 'waiting';
        try {
            // update settings private
            if ($settings_private.interval === '') {
                $settings_private.interval = '@every 3s';
            }
            await $pocketbase
                .collection('settings_private')
                .update($settings_private.id, $settings_private);

            // update settings public
            if (iconFiles.length > 0) {
                const formData = new FormData();
                formData.append('favicon', iconFiles[0]);
                await $pocketbase
                    .collection('settings_public')
                    .update($settings_public.id, formData);
            }
            await $pocketbase
                .collection('settings_public')
                .update($settings_public.id, $settings_public);

            //set button
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

    async function scanDevices() {
        buttons.scan.state = 'waiting';
        await $pocketbase
            .collection('settings_private')
            .update($settings_private.id, $settings_private)
            .catch((err) => {
                clearTimeout(timeout);
                timeout = setTimeout(() => {
                    buttons.scan.error = '';
                    buttons.scan.state = 'none';
                }, 3000);
                buttons.scan.error = err;
                buttons.scan.state = 'failed';
            });

        fetch(`${dev ? 'http://localhost:8090' : ''}/api/upsnap/scan`, {
            headers: {
                Authorization: $pocketbase.authStore.baseToken
            }
        })
            .then((res) => res.json())
            .then((data) => {
                if (data.message) {
                    clearTimeout(timeout);
                    timeout = setTimeout(() => {
                        buttons.scan.error = '';
                        buttons.scan.state = 'none';
                    }, 3000);
                    buttons.scan.error = data.message;
                    buttons.scan.state = 'failed';
                } else {
                    scannedDevices = data;
                    clearTimeout(timeout);
                    timeout = setTimeout(() => {
                        buttons.scan.state = 'none';
                    }, 3000);
                    buttons.scan.state = 'success';
                }
            })
            .catch((error) => {
                clearTimeout(timeout);
                timeout = setTimeout(() => {
                    buttons.scan.error = '';
                    buttons.scan.state = 'none';
                }, 3000);
                buttons.scan.error = error;
                buttons.scan.state = 'failed';
            });
    }

    async function addDevice(device) {
        await $pocketbase.collection('devices').create(device);
        $devices[device.id] = device;
    }
</script>

<div class="container">
    <UnauthorizedMsg />
    <section class="m-0 mt-4 p-4 shadow-sm">
        <form on:submit|preventDefault={saveSettings}>
            <div class="row">
                <div class="col-md-6">
                    <h3 class="mb-3 text-body-emphasis">Ping interval</h3>
                    <p>Sets the interval in which the devices are pinged.</p>
                    <div class="input-group mb-3">
                        <span class="input-group-text">Cron</span>
                        <input
                            class="form-control"
                            placeholder="e.g. '@every 5s' or '@every 1m'"
                            aria-label="Interval"
                            aria-describedby="addon-wrapping"
                            type="text"
                            bind:value={$settings_private.interval}
                        />
                    </div>
                    <h3 class="my-3 text-body-emphasis">Title</h3>
                    <p>Sets the title of the website and in the browser tab.</p>
                    <div class="input-group mb-3">
                        <span class="input-group-text" id="website-title">Title</span>
                        <input
                            type="text"
                            class="form-control"
                            placeholder="e.g. 'UpSnap'"
                            aria-label="Website title"
                            aria-describedby="website-title"
                            bind:value={$settings_public.website_title}
                        />
                    </div>
                    <h3 class="my-3 text-body-emphasis">Icon</h3>
                    <img
                        src=""
                        alt=""
                        class="img-fluid icon-preview mb-3"
                        bind:this={iconPreview}
                    />
                    <div class="mb-3">
                        <label for="iconInput" class="form-label">Set a custom favicon.</label>
                        <div class="input-group">
                            <input
                                class="form-control"
                                type="file"
                                id="iconInput"
                                accept=".ico,.png,.svg,.gif,.jpg,.jpeg"
                                bind:files={iconFiles}
                            />
                            <button
                                type="button"
                                class="btn btn-outline-danger"
                                on:click={() => {
                                    iconFiles = [];
                                    $settings_public.favicon = '';
                                }}>Back to Gopher</button
                            >
                        </div>
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="callout callout-info m-0">
                        <h5>Ping interval:</h5>
                        <p>
                            Leave blank to use default value of <span class="fw-bold"
                                >"@every 3s"</span
                            >. Learn more about the correct syntax for cron on
                            <a
                                href="https://en.wikipedia.org/wiki/Cron"
                                target="_blank"
                                rel="noreferrer">wikipedia</a
                            >
                            or refer to the
                            <a
                                target="_blank"
                                rel="noreferrer"
                                href="https://pkg.go.dev/github.com/robfig/cron"
                                >package documentation</a
                            >
                        </p>
                        <h5>Icon:</h5>
                        <p class="m-0">
                            Supported file types are
                            <span class="badge text-bg-secondary">.ico</span>,
                            <span class="badge text-bg-secondary">.png</span>,
                            <span class="badge text-bg-secondary">.svg</span>,
                            <span class="badge text-bg-secondary">.gif</span> and
                            <span class="badge text-bg-secondary">.jpg/.jpeg</span>.
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
    <DeviceForm bind:device={newDevice} mode="add" />
    <section class="m-0 my-4 p-4 shadow-sm">
        <div class="row">
            <div class="col-md-6">
                <h3 class="mb-3 text-body-emphasis">Network scan</h3>
                <p>Set the network address to scan.</p>
                <form on:submit|preventDefault={scanDevices}>
                    <div class="input-group mb-3">
                        <span class="input-group-text" id="ip-range">IP range</span>
                        <input
                            type="text"
                            class="form-control"
                            placeholder="e.g. '192.168.1.0/24'"
                            aria-label="ip-range"
                            aria-describedby="ip-range"
                            bind:value={$settings_private.scan_range}
                        />
                    </div>
                    <button
                        type="submit"
                        class="btn btn-secondary"
                        class:btn-success={buttons.scan.state === 'success' ? true : false}
                        class:btn-warning={buttons.scan.state === 'waiting' ? true : false}
                        class:btn-danger={buttons.scan.state === 'failed' ? true : false}
                        disabled={$settings_private.scan_range === '' ||
                        buttons.scan.state !== 'none'
                            ? true
                            : false}
                    >
                        {#if buttons.scan.state === 'none'}
                            Scan
                        {:else if buttons.scan.state === 'success'}
                            Scanned
                        {:else if buttons.scan.state === 'waiting'}
                            Scan running (might take some seconds)
                        {:else if buttons.scan.state === 'failed'}
                            Failed: {buttons.scan.error}
                        {/if}
                    </button>
                </form>
            </div>
            <div class="col-md-6">
                <div class="callout callout-info mb-0">
                    <p class="m-0">
                        For network scan to work, you need to run UpSnap as root/admin and have nmap
                        installed and available in your $PATH. (For docker users, thats already the
                        case and you don't need to do anything.)
                    </p>
                </div>
            </div>
        </div>
        {#if scannedDevices.devices}
            <div class="mt-3">
                <div class="row align-items-center my-1">
                    <div class="col-md fw-bold">Name</div>
                    <div class="col-md fw-bold">IP</div>
                    <div class="col-md fw-bold">MAC</div>
                    <div class="col-md fw-bold">Netmask</div>
                    <div class="col-md fw-bold">Add</div>
                </div>
                {#each scannedDevices?.devices as device}
                    <hr class="my-1" />
                    <div class="row align-items-center my-1">
                        <div class="col-md">
                            {device.name}
                        </div>
                        <div class="col-md">
                            {device.ip}
                        </div>
                        <div class="col-md">
                            {device.mac}
                        </div>
                        <div class="col-md">
                            {scannedDevices.netmask}
                        </div>
                        <div class="col-md">
                            <button
                                class="btn btn-sm btn-secondary py-0"
                                on:click={async (e) => {
                                    e.target.disabled = true;
                                    device.netmask = scannedDevices.netmask;
                                    await addDevice(device);
                                }}
                                on:keydown={async (e) => {
                                    e.target.disabled = true;
                                    device.netmask = scannedDevices.netmask;
                                    await addDevice(device);
                                }}
                            >
                                Add <Fa icon={faPlus} />
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </section>
    <p class="m-0 my-4 p-4 text-center text-muted">
        {#if version !== undefined}
            UpSnap version: <a
                href="https://github.com/seriousm4x/UpSnap/releases/tag/{version}"
                class="text-reset">{version}</a
            >
        {:else}
            UpSnap version: (untracked)
        {/if}
    </p>
</div>

<style lang="scss">
    .icon-preview {
        height: 130px;
    }
</style>
