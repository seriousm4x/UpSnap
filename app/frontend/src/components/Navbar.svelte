<script>
    import store from '../store.js';
    import DarkToggle from "./DarkToggle.svelte";
    export let visitors;
    export let settings;

    let addDevice = {
        wake: {
            enabled: false,
            cron: ""
        },
        shutdown: {
            enabled: false,
            cron: "",
            command: ""
        }
    }

    function updateDevice(data) {
        if (Object.keys(data).length < 4) {
            return
        }
        store.sendMessage({
            type: "update_device",
            data: data
        })
    }

    function updateSettings() {
        store.sendMessage({
            type: "update_settings",
            data: settings
        })
        hideModal("settings");
    }

    function scanNetwork() {
        store.sendMessage({
            type: "scan_network"
        })
        const btnScan = document.querySelector("#btnScan");
        const btnScanSpinner = document.querySelector("#btnScanSpinner");
        const btnScanText = document.querySelector("#btnScanText");
        btnScan.disabled = true;
        btnScanSpinner.classList.remove("d-none");
        btnScanText.innerText = "Scanning...";
    }

    function addScan(i) {
        document.querySelector(`#btnAdd${i}`).disabled = true;
        const dev = settings.scan_network[i];
        updateDevice(dev);
    }

	const restoreFromFile = (e) => {
        let file = e.target.files[0];
        let reader = new FileReader();
        reader.readAsText(file);
        reader.onload = e => {
            let data = JSON.parse(e.target.result);
            if (Array.isArray(data)) {
                // v2 file restore
                data.forEach((device) => {
                    updateDevice(device)
                })
            } else {
                // v1 file restore
                for (const [key, value] of Object.entries(data)) {
                    value["mac"] = key
                    updateDevice(value)
                }
            }
        }
        hideModal("settings");
    }

    function backupToFile() {
        store.sendMessage({
            type: "backup"
        })
    }

    function hideModal(id) {
        const modalEl = document.querySelector(`#${id}`);
        const modal = bootstrap.Modal.getInstance(modalEl);
        modal.hide();
    }

</script>

<nav class="navbar navbar-expand-sm">
    <div class="container-fluid">
        <div class="navbar-brand" href="/">
            <img src="/favicon.png" alt="Logo" width="24" height="24" class="me-2">
            UpSnap
        </div>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
            <span class="ms-auto d-flex">
                <DarkToggle />
                <div class="dropdown">
                    <button class="btn btn-light dropdown-toggle px-3 me-2 py-2" type="button" id="dropdownMenuButton1" data-bs-toggle="dropdown" aria-expanded="false">
                        <i class="fa-solid fa-wrench me-2"></i>More
                    </button>
                    <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
                        <li>
                            <button class="dropdown-item" data-bs-toggle="modal" data-bs-target="#addDevice">
                                <i class="fa-solid fa-plus me-2"></i>Add device
                            </button>
                        </li>
                        <li>
                            <button class="dropdown-item" data-bs-toggle="modal" data-bs-target="#settings">
                                <i class="fa-solid fa-sliders me-2"></i>Settings
                            </button>
                        </li>
                    </ul>
                </div>
                <div class="btn btn-light px-3 me-2 py-2 pe-none">
                    <i class="me-2 fa-solid {visitors === 1 ? "fa-user" : "fa-user-group"}"></i>{visitors} {visitors === 1 ? "Visitor" : "Visitors"}
                </div>
            </span>
        </div>
    </div>
</nav>

<div class="modal fade" id="addDevice" tabindex="-1" aria-labelledby="addDeviceLabel" aria-hidden="true">
    <div class="modal-dialog modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title fw-bold" id="addDeviceLabel">Add device</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <form id="addForm" on:submit|preventDefault={() => updateDevice(addDevice)}>
                    <div class="row">
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputNameAddDevice" class="form-label">Device name</label>
                                <input type="text" class="form-control" id="inputNameAddDevice" placeholder="Max PC" bind:value="{addDevice.name}" required>
                            </div>
                        </div>
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputMacAddDevice" class="form-label">Mac address</label>
                                <input type="text" class="form-control" id="inputMacAddDevice" placeholder="aa:aa:aa:aa:aa:aa" bind:value="{addDevice.mac}" pattern="^([0-9A-Fa-f]{'{'}2{'}'}[:-]){'{'}5{'}'}([0-9A-Fa-f]{'{'}2{'}'})|([0-9a-fA-F]{'{'}4{'}'}\\.[0-9a-fA-F]{'{'}4{'}'}\\.[0-9a-fA-F]{'{'}4{'}'})$" required>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputIpAddDevice" class="form-label">IP address</label>
                                <input type="text" class="form-control" id="inputIpAddDevice" placeholder="192.168.1.1" bind:value="{addDevice.ip}" pattern="^([01]?\d\d?|2[0-4]\d|25[0-5])(?:\.(?:[01]?\d\d?|2[0-4]\d|25[0-5])){'{'}3{'}'}$" required>
                            </div>
                        </div>
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputNetmaskAddDevice" class="form-label">Netmask</label>
                                <input type="text" class="form-control" id="inputNetmaskAddDevice" placeholder="255.255.255.0" bind:value="{addDevice.netmask}" pattern="^(((255\.){'{'}3{'}'}(255|254|252|248|240|224|192|128|0+))|((255\.){'{'}2{'}'}(255|254|252|248|240|224|192|128|0+)\.0)|((255\.)(255|254|252|248|240|224|192|128|0+)(\.0+){'{'}2{'}'})|((255|254|252|248|240|224|192|128|0+)(\.0+){'{'}3{'}'}))$" required>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-auto ms-auto">
                            <button type="submit" form="addForm" class="btn btn-outline-success">Add device</button>
                        </div>
                    </div>
                    <h5 class="fw-bold">Network discovery</h5>
                    {#if !settings.discovery}
                        <div class="callout callout-danger mb-2">
                            <p class="m-0">To enable this option, please enter your network address in the settings.</p>
                        </div>
                    {/if}
                    <button id="btnScan" class="btn btn-secondary" type="button" on:click={scanNetwork} disabled={!settings.discovery}>
                        <span id="btnScanSpinner" class="spinner-grow spinner-grow-sm d-none" role="status" aria-hidden="true"></span>
                        <span id="btnScanText">Scan</span>
                    </button>
                    {#if settings.scan_network?.length}
                        <table class="table">
                            <thead>
                                <tr>
                                    <td>Name</td>
                                    <td>IP</td>
                                    <td>Netmask</td>
                                    <td>Mac</td>
                                    <td>Add</td>
                                </tr>
                            </thead>
                            <tbody>
                                {#each settings.scan_network as device, i}
                                    <tr>
                                        <td>{device.name}</td>
                                        <td>{device.ip}</td>
                                        <td>{device.netmask}</td>
                                        <td>{device.mac}</td>
                                        <td>
                                            <button type="button" id="btnAdd{i}" class="btn btn-outline-secondary py-0" on:click={() => addScan(i)}>
                                                <i class="fa-solid fa-plus fa-sm"></i>
                                            </button>
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    {/if}
                </form>
            </div>
        </div>
    </div>
</div>

<div class="modal fade" id="settings" tabindex="-1" aria-labelledby="settingsLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title fw-bold" id="settingsLabel">Settings</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <form id="settingsForm" on:submit|preventDefault={updateSettings}>
                    <h5 class="fw-bold">General</h5>
                    <div class="row">
                        <div class="col-sm-8">
                            <div class="mb-3">
                                <label for="inputNetworkDiscovery" class="form-label">Network discovery address</label>
                                <input type="text" class="form-control" id="inputNetworkDiscovery" placeholder="192.168.1.0/24" bind:value="{settings.discovery}" pattern="^([01]?\d\d?|2[0-4]\d|25[0-5])(?:\.(?:[01]?\d\d?|2[0-4]\d|25[0-5])){'{'}2{'}'}(?:\.(?:0))(?:/[0-2]\d|/3[0-2])$">
                            </div>
                        </div>
                        <div class="col-sm-4">
                            <div class="mb-3">
                                <label for="inputIntervalSettings" class="form-label">Interval (seconds)</label>
                                <input type="number" class="form-control" id="inputIntervalSettings" min="5" bind:value="{settings.interval}" required>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm">
                            <div class="form-check">
                                <label class="form-check-label" for="flexCheckDefault">
                                    Enable notifications
                                </label>
                                <input class="form-check-input" type="checkbox" value="" id="flexCheckDefault" bind:checked={settings.notifications}>
                            </div>
                        </div>
                    </div>
                </form>
                <div class="row mb-3">
                    <div class="col-auto ms-auto">
                        <button type="submit" form="settingsForm" class="btn btn-outline-success">Save</button>
                    </div>
                </div>
                <h5 class="fw-bold">Backup/Restore</h5>
                <div class="callout callout-info mb-2">
                    <p class="m-0">Backup file structure has changed in v2. You can still restore both versions with this file upload.</p>
                </div>
                <div class="mb-3">
                    <label for="inputRestore" class="form-label">Restore from .json</label>
                    <input id="inputRestore" type="file" class="form-control" accept=".json" on:change={(e) => restoreFromFile(e)}>
                </div>
                <div class="mb-3">
                    <button type="button" class="btn btn-secondary" on:click={backupToFile}>
                        <i class="fa-solid fa-download me-2"></i>
                        Export .json
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>

<style lang="scss" global>
    .table {
        color: var(--color-text);
    }

    .navbar-brand {
        color: var(--color-text);
    }

    .dropdown-menu {
        background-color: var(--bg-lighter);

        li:first-child>button {
            border-radius: 0.5rem 0.5rem 0rem 0rem;
        }

        li:last-child>button {
            border-radius: 0rem 0rem 0.5rem 0.5rem;
        }

        li:only-child>button {
            border-radius: 0.5rem;
        }
    }

    .dropdown-item {
        color: var(--color-text);
        background-color: var(--bg-lighter);

        &:hover {
            color: var(--color-bg);
            background-color: var(--color-text);
        }
    }
</style>
