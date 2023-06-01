<script>
    import { goto } from '$app/navigation';
    import { pocketbase, devices } from '@stores/pocketbase';
    import Fa from 'svelte-fa/src/fa.svelte';
    import { faEye, faEyeSlash, faPlus, faTrashCan } from '@fortawesome/free-solid-svg-icons';

    export let device;
    export let mode;

    let timeout;
    let button = {
        state: 'none',
        error: ''
    };
    let deleteButton = {
        state: 'none',
        error: ''
    };
    let newPort = {
        name: '',
        number: null
    };
    let passwordShow = false;
    $: passwordType = passwordShow ? 'text' : 'password';

    async function addOrUpdateDevice() {
        button.state = 'waiting';
        try {
            // validate password length
            if (
                device.password.length != 0 &&
                device.password.length != 4 &&
                device.password.length != 6
            ) {
                throw 'Password must be 0, 4 or 6 characters long';
            }

            // add ports not in db
            for (let i = 0; i < device.expand.ports.length; i++) {
                const port = device.expand.ports[i];
                if (!port.id) {
                    const result = await $pocketbase.collection('ports').create(port);
                    device.ports = [...device.ports, result.id];
                } else {
                    await $pocketbase.collection('ports').update(port.id, port);
                }
            }

            // create or update device
            if (mode === 'add') {
                await $pocketbase.collection('devices').create(device);
            } else {
                await $pocketbase.collection('devices').update(device.id, device);
            }

            // show button with timeout
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                button.state = 'none';
            }, 3000);
            button.state = 'success';
        } catch (error) {
            clearTimeout(timeout);
            setTimeout(() => {
                button.error = '';
                button.state = 'none';
            }, 3000);
            button.error = error;
            button.state = 'failed';
        }
    }

    async function deleteDevice() {
        deleteButton.state = 'waiting';
        try {
            await $pocketbase.collection('devices').delete(device.id);
            device.ports.forEach(async (port) => {
                await $pocketbase.collection('ports').delete(port);
            });
            delete $devices[device.id];
            goto('/');
        } catch (error) {
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                deleteButton.error = '';
                deleteButton.state = 'none';
            }, 3000);
            deleteButton.error = error;
            deleteButton.state = 'failed';
        }
    }

    async function deletePort(idx) {
        // delete port from db if it has an id
        // ports with id exist in db, ports without id are not created yet
        const port = device.expand.ports[idx];
        if (port.id) {
            await $pocketbase.collection('ports').delete(port.id);
            const i = device.ports.indexOf(port.id);
            if (i !== -1) {
                device.ports.splice(i, 1);
                device.ports = device.ports;
            }
        }
        device.expand.ports.splice(idx, 1);
        device.expand.ports = device.expand.ports;
    }

    function addPort() {
        device.expand.ports = [...device.expand.ports, JSON.parse(JSON.stringify(newPort))];
        newPort.name = '';
        newPort.number = null;
    }

    function onPasswordInput(event) {
        device.password = event.target.value;
    }
</script>

<section class="m-0 mt-4 p-4 shadow-sm">
    <h3 class="mb-3 text-body-emphasis">{mode === 'add' ? 'Add new device' : device.name}</h3>
    <div class="row">
        <div class="col-md-6">
            <form on:submit|preventDefault={addOrUpdateDevice}>
                <h5>Required:</h5>
                <div class="input-group mb-1">
                    <span class="input-group-text">Name</span>
                    <input
                        class="form-control"
                        placeholder="Office Pc"
                        aria-label="Name"
                        aria-describedby="addon-wrapping"
                        type="text"
                        required
                        bind:value={device.name}
                    />
                </div>
                <div class="input-group mb-1">
                    <span class="input-group-text">IP</span>
                    <input
                        class="form-control"
                        placeholder="192.168.1.34"
                        aria-label="IP"
                        aria-describedby="addon-wrapping"
                        type="text"
                        required
                        bind:value={device.ip}
                    />
                </div>
                <div class="input-group mb-1">
                    <span class="input-group-text">MAC</span>
                    <input
                        class="form-control"
                        placeholder="aa:bb:cc:dd:ee:ff"
                        aria-label="MAC"
                        aria-describedby="addon-wrapping"
                        type="text"
                        required
                        bind:value={device.mac}
                    />
                </div>
                <div class="input-group">
                    <span class="input-group-text">Netmask</span>
                    <input
                        class="form-control"
                        placeholder="Most likely 255.255.255.0 or 255.255.0.0"
                        aria-label="Netmask"
                        aria-describedby="addon-wrapping"
                        type="text"
                        required
                        bind:value={device.netmask}
                    />
                </div>
                <h5 class="mt-4">Optional:</h5>
                {#if device?.expand?.ports}
                    {#each device.expand.ports as port, idx}
                        <div class="input-group mb-1">
                            <span class="input-group-text">Port</span>
                            <input
                                type="text"
                                aria-label="Name"
                                class="form-control"
                                placeholder="Name"
                                bind:value={port.name}
                            />
                            <input
                                type="number"
                                min="0"
                                max="65535"
                                aria-label="Number"
                                class="form-control"
                                placeholder="Number"
                                bind:value={port.number}
                            />
                            <button
                                type="button"
                                class="btn btn-outline-secondary"
                                on:click={() => deletePort(idx)}
                            >
                                <Fa icon={faTrashCan} />
                            </button>
                        </div>
                    {/each}
                {/if}
                <div class="input-group mb-3">
                    <span class="input-group-text">Port</span>
                    <input
                        type="text"
                        aria-label="Name"
                        class="form-control"
                        placeholder="Name"
                        bind:value={newPort.name}
                    />
                    <input
                        type="number"
                        min="0"
                        max="65535"
                        aria-label="Number"
                        class="form-control"
                        placeholder="Number"
                        bind:value={newPort.number}
                    />
                    <button
                        type="button"
                        class="btn btn-outline-secondary"
                        on:click={() => addPort()}
                    >
                        <Fa icon={faPlus} />
                    </button>
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text">Link</span>
                    <input
                        class="form-control"
                        placeholder="Clickable link on device card"
                        aria-label="Link"
                        aria-describedby="addon-wrapping"
                        type="url"
                        bind:value={device.link}
                    />
                </div>
                <div class="input-group mb-1">
                    <span class="input-group-text">Wake Cron<sup>(1)</sup></span>
                    <input
                        class="form-control"
                        placeholder="Automatically wake device with cron"
                        aria-label="Wake Cron"
                        aria-describedby="addon-wrapping"
                        type="text"
                        bind:value={device.wake_cron}
                    />
                    <div class="input-group-text form-switch">
                        <label class="form-check-label" for="wake-cron-enabled">Enable</label>
                        <input
                            id="wake-cron-enabled"
                            class="form-check-input mt-0 ms-1"
                            type="checkbox"
                            role="switch"
                            aria-label="Enable/Disable wake cron"
                            bind:checked={device.wake_cron_enabled}
                        />
                    </div>
                </div>
                <div class="input-group mb-1">
                    <span class="input-group-text">Shutdown Cron<sup>(1)</sup></span>
                    <input
                        class="form-control"
                        placeholder="Automatically shutdown device with cron"
                        aria-label="Shutdown Cron"
                        aria-describedby="addon-wrapping"
                        type="text"
                        bind:value={device.shutdown_cron}
                    />
                    <div class="input-group-text form-switch">
                        <label class="form-check-label" for="wake-shutdown-enabled">Enable</label>
                        <input
                            id="wake-shutdown-enabled"
                            class="form-check-input mt-0 ms-1"
                            type="checkbox"
                            role="switch"
                            aria-label="Enable/Disable wake cron"
                            bind:checked={device.shutdown_cron_enabled}
                        />
                    </div>
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text">Shutdown Cmd<sup>(2)</sup></span>
                    <input
                        class="form-control"
                        placeholder="Command to shutdown device"
                        aria-label="Shutdown Cmd"
                        aria-describedby="addon-wrapping"
                        type="text"
                        bind:value={device.shutdown_cmd}
                    />
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text">Password<sup>(3)</sup></span>
                    <input
                        class="form-control"
                        placeholder="BIOS password for wol"
                        aria-label="IP"
                        aria-describedby="addon-wrapping"
                        type={passwordType}
                        value={device.password}
                        maxlength="6"
                        on:input={onPasswordInput}
                    />
                    <button
                        class="btn btn-outline-secondary"
                        type="button"
                        on:click={() => (passwordShow = !passwordShow)}
                    >
                        <Fa icon={passwordShow ? faEyeSlash : faEye} />
                    </button>
                </div>
                <button
                    type="submit"
                    class="btn btn-secondary"
                    class:btn-success={button.state === 'success' ? true : false}
                    class:btn-warning={button.state === 'waiting' ? true : false}
                    class:btn-danger={button.state === 'failed' ? true : false}
                    disabled={button.state !== 'none' ? true : false}
                >
                    {#if button.state === 'none'}
                        {mode === 'add' ? 'Add device' : 'Save device'}
                    {:else if button.state === 'success'}
                        {mode === 'add' ? 'Added successfully' : 'Saved successfully'}
                    {:else if button.state === 'waiting'}
                        Waiting
                    {:else if button.state === 'failed'}
                        Failed: {button.error}
                    {/if}
                </button>
            </form>
        </div>
        <div class="col-md-6 d-flex align-items-start flex-column">
            <div class="callout callout-info mb-auto">
                <h5>Optional:</h5>
                <p class="m-0">
                    (1) Learn more about the correct syntax for cron on
                    <a href="https://en.wikipedia.org/wiki/Cron" target="_blank" rel="noreferrer"
                        >wikipedia</a
                    >
                    or refer to the
                    <a
                        target="_blank"
                        rel="noreferrer"
                        href="https://pkg.go.dev/github.com/robfig/cron">package documentation</a
                    >
                </p>
                <p class="m-0">
                    (2) Shell command to be executed. "net rpc", "sshpass" and "curl" are available.
                    e.g.:
                </p>
                <ul>
                    <li>
                        Windows: <code>net rpc shutdown -I 192.168.1.13 -U "user%password"</code>
                    </li>
                    <li>
                        Linux: <code
                            >sshpass -p your_password ssh -o "StrictHostKeyChecking=no"
                            user@hostname "sudo poweroff"</code
                        >
                    </li>
                </ul>
                <p class="m-0">
                    (3) Some network cards have the option to set a password for magic packets, also
                    called "SecureON". Password can only be 0, 4 or 6 characters in length.
                </p>
            </div>
            {#if mode === 'edit'}
                <div class="text-end mt-3 align-self-end">
                    <button
                        class="btn btn-danger"
                        on:click={() => deleteDevice()}
                        on:keydown={() => deleteDevice()}
                    >
                        {#if deleteButton.state === 'none'}
                            Delete device
                        {:else if deleteButton.state === 'waiting'}
                            Waiting
                        {:else if deleteButton.state === 'failed'}
                            Failed: {deleteButton.error}
                        {/if}
                    </button>
                </div>
            {/if}
        </div>
    </div>
</section>
