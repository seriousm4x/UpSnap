<script>
    import socketStore from '../socketStore.js';
    export let device;

    let modalDevice = JSON.parse(JSON.stringify(device));
    let customPort = {}

	function wake(id) {
        socketStore.sendMessage({
            type: "wake",
            id: id
        })
	}

    function shutdown(id) {
        socketStore.sendMessage({
            type: "shutdown",
            id: id
        })
	}

	function deleteDevice() {
        socketStore.sendMessage({
            type: "delete_device",
            id: modalDevice.id
        })
	}

    function updateDevice() {
        device = modalDevice;
        socketStore.sendMessage({
            type: "update_device",
            data: modalDevice
        })
    }

    function updatePort() {
        if (!customPort.number) {
            return
        }
        const index = modalDevice.ports.findIndex(x => x.number == customPort.number);
        if (customPort.name) {
            // add port
            if (index === -1) {
                customPort.checked = true;
                modalDevice.ports.push(JSON.parse(JSON.stringify(customPort)));
            } else {
                customPort.checked = modalDevice.ports[index].checked;
                modalDevice.ports[index] = JSON.parse(JSON.stringify(customPort));
            }
        } else {
            // delete port
            if (index >= 0) {
                modalDevice.ports.splice(index, 1)
            }
        }
        modalDevice = modalDevice;
        // send to backend
        socketStore.sendMessage({
            type: "update_port",
            data: customPort
        })
    }

    function openModal() {
        modalDevice = Object.assign({}, device)
    }

    function validatePort() {
        if (typeof customPort.number != "number") {
            customPort.number = 1
        } else if (customPort.number > 65535) {
            customPort.number = 65535
        }
    }
</script>

<div id="device-col-{device.id}" class="col-xs-12 col-sm-6 col-md-4 col-lg-3 g-4">
    <div class="card border-0 p-3 pt-2">
        <div class="card-body">
            <div class="row">
                <div class="col-auto me-auto">
                    <div id="spinner-{device.id}" class="spinner-border warning d-none" role="status"></div>
                    {#if device.up === true}
                        {#if device.shutdown.command}
                        <div class="hover" on:click="{() => shutdown(device.id)}" data-bs-toggle="tooltip" title="Shutdown command: {device.shutdown.command}" role="button">
                            <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x success"></i>
                        </div>
                        {:else}
                            <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x success"></i>
                        {/if}
                    {:else if device.up === false}
                        <div class="hover" on:click="{() => wake(device.id)}" role="button">
                            <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x danger"></i>
                        </div>
                    {:else}
                        <div class="hover" role="button">
                            <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x text-muted"></i>
                        </div>
                    {/if}
                </div>
                {#if device.wake.enabled}
                    <div class="col-auto px-2" data-bs-toggle="tooltip" title="Wake cron: {device.wake.cron}">
                        <i class="fa-solid fa-circle-play fa-2x text-muted"></i>
                    </div>
                {/if}
                {#if device.shutdown.enabled}
                    <div class="col-auto px-2" data-bs-toggle="tooltip" title="Shutdown cron: {device.wake.cron}">
                        <i class="fa-solid fa-circle-stop fa-2x text-muted"></i>
                    </div>
                {/if}
                <div class="col-auto hover" data-bs-toggle="modal" data-bs-target="#device-modal-{device.id}" role="button" on:click="{() => openModal()}">
                    <i class="fa-solid fa-ellipsis-vertical fa-2x"></i>
                </div>
            </div>
            {#if device.link}
                <h5 class="card-title fw-bold my-2"><a class="inherit-color" href={device.link}>{device.name}</a></h5>
            {:else}
                <h5 class="card-title fw-bold my-2">{device.name}</h5>
            {/if}
            <h6 class="card-subtitle mb-2 text-muted">{device.ip}</h6>
            <ul class="list-group">
                {#each device.ports as port}
                    {#if port.checked === true}
                        <li class="list-group-item">
                            <i class="fa-solid fa-circle align-middle {port.open ? 'success' : 'danger'}"></i>
                            {port.name} <span class="text-muted">({port.number})</span>
                        </li>
                    {/if}
                {/each}
            </ul>
        </div>
    </div>

    <div class="modal fade" id="device-modal-{device.id}" tabindex="-1" aria-labelledby="device-modal-{device.id}-label" aria-hidden="true">
        <div class="modal-dialog modal-dialog-scrollable">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title fw-bold" id="device-modal-{modalDevice.id}-label">{modalDevice.name}</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <form id="form-{modalDevice.id}" on:submit|preventDefault={updateDevice}>
                        <!-- general -->
                        <h5 class="fw-bold">General</h5>
                        <div class="row mb-2">
                            <div class="col-sm">
                                <label for="inputName{modalDevice.id}" class="form-label">Device name</label>
                                <input type="text" class="form-control" id="inputName{modalDevice.id}" bind:value="{modalDevice.name}" required>
                            </div>
                            <div class="col-sm">
                                <label for="inputMac{modalDevice.id}" class="form-label">Mac address</label>
                                <input type="text" class="form-control" id="inputMac{modalDevice.mac}" bind:value="{modalDevice.mac}" required>
                            </div>
                        </div>
                        <div class="row mb-2">
                            <div class="col-sm">
                                <label for="inputIp{modalDevice.id}" class="form-label">IP address</label>
                                <input type="text" class="form-control" id="inputIp{modalDevice.id}" bind:value="{modalDevice.ip}" pattern="^([01]?\d\d?|2[0-4]\d|25[0-5])(?:\.(?:[01]?\d\d?|2[0-4]\d|25[0-5])){'{'}3{'}'}$" required>
                            </div>
                            <div class="col-sm">
                                <label for="inputNetmask{modalDevice.id}" class="form-label">Netmask</label>
                                <input type="text" class="form-control" id="inputNetmask{modalDevice.id}" bind:value="{modalDevice.netmask}" pattern="^(((255\.){'{'}3{'}'}(255|254|252|248|240|224|192|128|0+))|((255\.){'{'}2{'}'}(255|254|252|248|240|224|192|128|0+)\.0)|((255\.)(255|254|252|248|240|224|192|128|0+)(\.0+){'{'}2{'}'})|((255|254|252|248|240|224|192|128|0+)(\.0+){'{'}3{'}'}))$" required>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col">
                                <div class="mb-3">
                                    <label for="inputLinkAddDevice" class="form-label">Web link</label>
                                    <input type="text" class="form-control" id="inputILinkAddDevice" placeholder="http://...." bind:value="{modalDevice.link}">
                                </div>
                            </div>
                        </div>
                        <!-- ports -->
                        <h5 class="fw-bold mt-4">Ports</h5>
                        <p class="mb-2">Select ports to check if they are open.</p>
                        {#if modalDevice.ports.length === 0}
                            <p class="mb-0">No ports available. Add ports below.</p>
                        {/if}
                        {#each modalDevice.ports as port}
                            <div class="form-check">
                                <input type="checkbox" class="form-check-input" id="{device.id}-port-{port.number}" bind:checked="{port['checked']}">
                                <label class="form-check-label" for="{device.id}-port-{port.number}">{port.name} <span class="text-muted">({port.number})</span></label>
                            </div>
                        {/each}
                        <label class="form-label mt-2" for="{device.id}-custom-port">Custom port</label>
                        <div class="input-group mb-2">
                            <input type="text" id="{device.id}-custom-port" class="form-control rounded-0 rounded-start" placeholder="Name" aria-label="Name" aria-describedby="button-addon2" bind:value={customPort.name}>
                            <input type="number" min="1" max="65535" class="form-control rounded-0" placeholder="Port" aria-label="Port" aria-describedby="button-addon2" bind:value={customPort.number} on:input={validatePort}>
                            <button class="btn btn-secondary" type="button" id="button-addon2" on:click="{updatePort}">Update Port</button>
                        </div>
                        <button class="btn btn-secondary mt-2" type="button" data-bs-toggle="collapse" data-bs-target="#info-ports" aria-expanded="false" aria-controls="info-ports">
                            <i class="fa-solid fa-angle-down me-2"></i>How to use
                        </button>
                        <div class="collapse mt-3" id="info-ports">
                            <div class="callout callout-info">
                                <p class="mb-0">Ports must be between 1 and 65535. Enter the same port with a differen name to change it. Leave name empty to delete port.</p>
                            </div>
                        </div>
                        <!-- scheduled wake -->
                        <h5 class="fw-bold mt-4">Scheduled wake</h5>
                        <p class="mb-2">Wake your device at a given time.</p>
                        <div class="form-check">
                            <input class="form-check-input" type="radio" name="wake-disable" id="wake-radio-disabled-{modalDevice.id}"
                                bind:group={modalDevice["wake"]["enabled"]} value={false} checked={!modalDevice["wake"]["enabled"]}>
                            <label class="form-check-label" for="wake-radio-disabled-{modalDevice.id}">
                                Disabled
                            </label>
                        </div>
                        <div class="form-check">
                            <input class="form-check-input" type="radio" name="wake-cron" id="wake-radio-cron-{modalDevice.id}"
                                bind:group={modalDevice["wake"]["enabled"]} value={true} checked={modalDevice["wake"]["enabled"]}>
                            <label class="form-check-label" for="wake-radio-cron-{modalDevice.id}">
                                Enabled
                            </label>
                        </div>
                        <div class="input-group my-1" hidden={!modalDevice["wake"]["enabled"]}>
                            <span class="input-group-text rounded-0 rounded-start" id="wake-cron-{modalDevice.id}">Cron</span>
                            <input type="text" class="form-control rounded-0 rounded-end" placeholder="* /4 * * *" aria-label="Crontab" aria-describedby="wake-cron-{modalDevice.id}" bind:value={modalDevice["wake"]["cron"]}>
                        </div>
                        <button class="btn btn-secondary mt-2" type="button" data-bs-toggle="collapse" data-bs-target="#info-wake" aria-expanded="false" aria-controls="info-wake">
                            <i class="fa-solid fa-angle-down me-2"></i>How to use
                        </button>
                        <div class="collapse mt-3" id="info-wake">
                            <div class="callout callout-info">
                                <p class="mb-2">Cron is a syntax describing a time pattern when to execute jobs. The above field uses common cron syntax. Examples:</p>
                                <pre class="mb-2">Minute Hour DayOfMonth Month DayOfWeek
                                    * /4 * * * (Wake every 4 hours)
                                    0 9 * * 1-5 (Wake from Mo-Fr at 9 a.m.)
                                </pre>
                                <p class="mb-0">Read more about <a href="https://linux.die.net/man/5/crontab" target="_blank">valid syntax here</a> or <a href="https://crontab.guru/" target="_blank">use a generator</a>. Expressions starting with "@..." are not supported.</p>
                            </div>
                        </div>
                        <!-- scheduled shutdown -->
                        <h5 class="fw-bold mt-4">Shutdown</h5>
                        <p class="mb-2">Set the shutdown command here. This shell command will be executed when clicking the power button on the device card. You can use cron below, which will then execute the command at the given time.</p>
                        <div class="input-group">
                            <span class="input-group-text rounded-0 rounded-start" id="shutdown-command-{modalDevice.id}">Command</span>
                            <input type="text" class="form-control rounded-0 rounded-end" placeholder="sshpass -p your_password ssh -o 'StrictHostKeyChecking=no' user@hostname 'sudo shutdown'" aria-label="Ccommand" aria-describedby="shutdown-command-{modalDevice.id}" bind:value={modalDevice["shutdown"]["command"]}>
                        </div>
                        <button class="btn btn-secondary mt-2" type="button" data-bs-toggle="collapse" data-bs-target="#info-shutdown" aria-expanded="false" aria-controls="info-shutdown">
                            <i class="fa-solid fa-angle-down me-2"></i>How to use
                        </button>
                        <div class="collapse mt-3" id="info-shutdown">
                            <div class="callout callout-info">
                                <p class="mb-2">This field takes a shell command to trigger the shutdown. You can use <code>sshpass</code> for Linux or <code>net rpc</code> for Windows hosts.</p>
                                <div class="callout callout-danger mb-2">Note: This command is safed as cleartext. Meaning, passwords are clearly visible in the database.</div>
                                <p class="mb-2">Examples:</p>
                                <pre class="mb-2"># wake linux hosts
                                    sshpass -p your_password ssh -o 'StrictHostKeyChecking=no' user@hostname 'sudo shutdown'

                                    # wake windows hosts
                                    net rpc shutdown --ipaddress 192.168.1.1 --user user%password
                                </pre>
                                <p class="mb-0">Read more about <a href="https://linux.die.net/man/1/sshpass" target="_blank">sshpass</a> or <a href="https://linux.die.net/man/8/net" target="_blank">net rpc</a>.</p>
                            </div>
                        </div>
                        <div>
                            <label class="form-label mt-2" for="{device.id}-shutdown-cron">Use cron</label>
                        </div>
                        <div class="form-check" id="{device.id}-shutdown-cron">
                            <input class="form-check-input" type="radio" name="shutdown-disable" id="shutdown-radio-disabled-{modalDevice.id}"
                                bind:group={modalDevice["shutdown"]["enabled"]} value={false} checked={!modalDevice["shutdown"]["enabled"]}>
                            <label class="form-check-label" for="shutdown-radio-disabled-{modalDevice.id}">
                                Disabled
                            </label>
                        </div>
                        <div class="form-check">
                            <input class="form-check-input" type="radio" name="shutdown-enable" id="shutdown-radio-enabled-{modalDevice.id}"
                                bind:group={modalDevice["shutdown"]["enabled"]} value={true} checked={modalDevice["shutdown"]["enabled"]}>
                            <label class="form-check-label" for="shutdown-radio-enabled-{modalDevice.id}">
                                Enabled
                            </label>
                        </div>
                        <div class="input-group my-1" hidden={!modalDevice["shutdown"]["enabled"]}>
                            <span class="input-group-text rounded-0 rounded-start" id="shutdown-cron-{modalDevice.id}">Cron</span>
                            <input type="text" class="form-control rounded-0 rounded-end" placeholder="* /4 * * *" aria-label="Crontab" aria-describedby="shutdown-cron-{modalDevice.id}" bind:value={modalDevice["shutdown"]["cron"]}>
                        </div>
                    </form>
                </div>
                <div class="modal-footer justify-content-between">
                    <button type="button" class="btn btn-outline-danger" data-bs-dismiss="modal" on:click={deleteDevice}>Delete</button>
                    <button type="submit" form="form-{modalDevice.id}" class="btn btn-outline-success">Save changes</button>
                </div>
            </div>
        </div>
    </div>
</div>

<style lang="scss">
    .success {
        color: var(--success);
    }

    .warning {
        color: var(--warning);
    }

    .danger {
        color: var(--danger);
    }

    pre {
        white-space: pre-line;
        background-color: var(--color-bg);
        padding: 1em;
        border-radius: 1em;
    }

    .card {
        border-radius: 2em;
		background-color: var(--bg-lighter);
    }

    .spinner-border {
        width: 1.5em;
        height: 1.5em;
    }

    .hover {
        &:hover {
            text-shadow: 0px 0px 20px rgb(155, 155, 155);
            transition: all 0.1s;
        }
    }

    .fa-2x {
        font-size: 1.5em;
    }

    .fa-power-off {
        border-radius: 1em;
    }

    .fa-circle {
        font-size: 0.8em;
    }

    .list-group-item {
        color: var(--color-text);
    }

    .inherit-color {
        color: inherit;
    }
</style>
