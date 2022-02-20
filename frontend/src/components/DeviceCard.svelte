<script>
    import store from '../store.js';
    export let device;

    let modalDevice = JSON.parse(JSON.stringify(device));
    let customPort = {}

	function wake(id) {
        store.sendMessage({
            type: "wake",
            id: id
        })
	}

	function deleteDevice() {
        store.sendMessage({
            type: "delete_device",
            id: modalDevice.id
        })
	}

    function updateDevice() {
        device = modalDevice;
        store.sendMessage({
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
                modalDevice.ports.push(customPort);
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
        store.sendMessage({
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
    <div class="card border-0 p-3 pt-2" >
        <div class="card-body">
            <div class="row">
                <div class="col-auto me-auto">
                    <div id="spinner-{device.id}" class="spinner-border color-warning d-none" role="status"></div>
                    <div class="hover" on:click="{() => wake(device.id)}" role="button">
                        {#if device.up === true || device.up === false}
                            <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x {device.up ? 'success' : 'danger'}"></i>
                        {:else}
                            <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x text-muted"></i>
                        {/if}
                    </div>
                </div>
                <div class="col-auto hover" data-bs-toggle="modal" data-bs-target="#device-modal-{device.id}" role="button" on:click="{() => openModal()}">
                    <i class="fa-solid fa-ellipsis-vertical fa-2x"></i>
                </div>
            </div>
            <h5 class="card-title fw-bold my-2">{device.name}</h5>
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
                        <h5 class="fw-bold">General</h5>
                        <div class="row">
                            <div class="col-sm">
                                <div class="mb-3">
                                    <label for="inputName{modalDevice.id}" class="form-label">Device name</label>
                                    <input type="text" class="form-control" id="inputName{modalDevice.id}" bind:value="{modalDevice.name}" required>
                                </div>
                            </div>
                            <div class="col-sm">
                                <div class="mb-3">
                                    <label for="inputMac{modalDevice.id}" class="form-label">Mac address</label>
                                    <input type="text" class="form-control" id="inputMac{modalDevice.mac}" bind:value="{modalDevice.mac}" required>
                                </div>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col-sm">
                                <div class="mb-3">
                                    <label for="inputIp{modalDevice.id}" class="form-label">IP address</label>
                                    <input type="text" class="form-control" id="inputIp{modalDevice.id}" bind:value="{modalDevice.ip}" pattern="^([01]?\d\d?|2[0-4]\d|25[0-5])(?:\.(?:[01]?\d\d?|2[0-4]\d|25[0-5])){'{'}3{'}'}$" required>
                                </div>
                            </div>
                            <div class="col-sm">
                                <div class="mb-3">
                                    <label for="inputNetmask{modalDevice.id}" class="form-label">Netmask</label>
                                    <input type="text" class="form-control" id="inputNetmask{modalDevice.id}" bind:value="{modalDevice.netmask}" pattern="^(((255\.){'{'}3{'}'}(255|254|252|248|240|224|192|128|0+))|((255\.){'{'}2{'}'}(255|254|252|248|240|224|192|128|0+)\.0)|((255\.)(255|254|252|248|240|224|192|128|0+)(\.0+){'{'}2{'}'})|((255|254|252|248|240|224|192|128|0+)(\.0+){'{'}3{'}'}))$" required>
                                </div>
                            </div>
                        </div>
                        <h5 class="fw-bold">Ports</h5>
                        {#each modalDevice.ports as port}
                            <div class="form-check">
                                <input type="checkbox" class="form-check-input" id="{device.id}-port-{port.name}" bind:checked="{port['checked']}">
                                <label class="form-check-label" for="{device.id}-port-{port.name}">{port.name} <span class="text-muted">({port.number})</span></label>
                            </div>
                        {/each}
                        <label class="form-label mt-3" for="{device.id}-custom-port">Custom port</label>
                        <div class="input-group mb-2">
                            <input type="text" id="{device.id}-custom-port" class="form-control" placeholder="Name" aria-label="Name" aria-describedby="button-addon2" bind:value={customPort.name}>
                            <input type="number" min="1" max="65535" class="form-control" placeholder="Port" aria-label="Port" aria-describedby="button-addon2" bind:value={customPort.number} on:input={validatePort}>
                            <button class="btn btn-outline-secondary" type="button" id="button-addon2" on:click="{updatePort}">Update Port</button>
                        </div>
                        <p>Port must be between 1 and 65535. Enter the same port with a differen name to change it. Leave name empty to delete port.</p>
                        <h5 class="fw-bold">Scheduled wake</h5>
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
    @import "../variables.scss";

    .card {
        border-radius: 2em;
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
        color: $success;
    }

    i {
        &.success {
            color: $success;
        }

        &.danger {
            color: $danger;
        }
    }

    button {
        &.btn-outline-success {
            border-color: $success;
            &:hover {
                background-color: $success;
            }
        }

        &.btn-outline-danger {
            border-color: $danger;
            &:hover {
                background-color: $danger;
            }
        }
    }

    .color-warning {
        color: $warning;
    }

</style>
