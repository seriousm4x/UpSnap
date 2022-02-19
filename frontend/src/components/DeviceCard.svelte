<script>
    import store from '../store.js';
    export let device;
    export let ports;

    let modalDevice = JSON.parse(JSON.stringify(device));
    // for (const [key, _] of Object.entries(modalDevice.ports)) {
    //     modalDevice.ports[key]["checked"] = true;
    // }

	function wake(id) {
        store.sendMessage({
            type: "wake",
            id: id
        })
	}

    function openModal() {
        modalDevice = Object.assign({}, device)

        store.sendMessage({
            type: "get_ports"
        })
    }

	function del(id) {
        store.sendMessage({
            type: "delete",
            id: id
        })
	}

    function update(data) {
        store.sendMessage({
            type: "update",
            data: data
        })
    }

</script>

<div id="device-col-{device.id}" class="col-xs-12 col-sm-6 col-md-4 col-lg-3 g-4">
    <div class="card border-0 p-3 pt-2" >
        <div class="card-body">
            <div class="row">
                <div class="col-auto me-auto">
                    <div id="spinner-{device.id}" class="spinner-border color-warning d-none" role="status"></div>
                    <div class="hover" on:click="{() => wake(device.id)}" role="button">
                        <i id="dot-{device.id}" class="fa-solid fa-power-off fa-2x {device.up ? 'success' : 'danger'}"></i>
                    </div>
                </div>
                <div class="col-auto hover" data-bs-toggle="modal" data-bs-target="#device-modal-{device.id}" role="button" on:click="{() => openModal()}">
                    <i class="fa-solid fa-ellipsis-vertical fa-2x"></i>
                </div>
            </div>
            <h5 class="card-title fw-bold my-2">{device.name}</h5>
            <h6 class="card-subtitle mb-2 text-muted">{device.ip}</h6>
            <ul class="list-group">
                {#each Object.entries(device.ports) as [_, port]}
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
        <div class="modal-dialog">
            <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title fw-bold" id="device-modal-{modalDevice.id}-label">{modalDevice.name}</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <form>
                    <h5 class="fw-bold">General</h5>
                    <div class="mb-3">
                        <label for="inputName{modalDevice.id}" class="form-label">Device name</label>
                        <input type="text" class="form-control" id="inputName{modalDevice.id}" bind:value="{modalDevice.name}">
                    </div>
                    <div class="row">
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputIp{modalDevice.id}" class="form-label">IP address</label>
                                <input type="text" class="form-control" id="inputIp{modalDevice.id}" bind:value="{modalDevice.ip}">
                            </div>
                        </div>
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputNetmask{modalDevice.id}" class="form-label">Netmask</label>
                                <input type="text" class="form-control" id="inputNetmask{modalDevice.id}" bind:value="{modalDevice.netmask}">
                            </div>
                        </div>
                    </div>
                    <h5 class="fw-bold">Ports</h5>
                    {#each ports as port}
                        <div class="form-check">
                            <input type="checkbox" class="form-check-input" id="{device.id}-port-{port.fields.name}" bind:checked="{modalDevice.ports[port.fields.number]['checked']}" >
                            <label class="form-check-label" for="{device.id}-port-{port.fields.name}">{port.fields.name}</label>
                        </div>
                    {/each}
                </form>
            </div>
            <div class="modal-footer justify-content-between">
                <button type="button" class="btn btn-outline-danger" data-bs-dismiss="modal" on:click="{() => del(modalDevice.id)}">Delete</button>
                <button type="button" class="btn btn-outline-success" on:click="{() => update(modalDevice)}">Save changes</button>
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
