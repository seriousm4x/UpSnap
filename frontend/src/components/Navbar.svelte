<script>
    import store from '../store.js';
    export let visitors;

    let label = visitors === 1 ? "Visitor" : "Visitors";
    let icon = visitors === 1 ? "fa-user" : "fa-user-group";

    let addDevice = {}

    function updateDevice() {
        if (Object.keys(addDevice).length < 4) {
            return
        }
        store.sendMessage({
            type: "update_device",
            data: addDevice
        })
    }

</script>

<nav class="navbar navbar-expand-sm navbar-light bg-light">
<div class="container-fluid">
    <a class="navbar-brand" href="/">
        <img src="/favicon.png" alt="Logo" width="24" height="24">
    </a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
        <div class="navbar-nav me-auto">
            <a class="nav-link {window.location.pathname == "/" ? 'active' : ''}" href="/">Devices</a>
            <a class="nav-link {window.location.pathname == "/settings" ? 'active' : ''}" href="/settings">Settings</a>
        </div>
        <span class="d-flex">
            <button type="button" class="btn btn-light py-2 px-3 me-2" data-bs-toggle="modal" data-bs-target="#addDevice">
                Add device<i class="ms-1 fa-solid fa-plus"></i>
            </button>
            <div class="card border-0">
                <div class="card-body py-2 px-3">
                    {visitors} {label}<i class="ms-1 fa-solid {icon} text-muted"></i>
                </div>
            </div>
        </span>
    </div>
</nav>

<div class="modal fade" id="addDevice" tabindex="-1" aria-labelledby="addDeviceLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title fw-bold" id="addDeviceLabel">Add device</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <form id="addForm" on:submit|preventDefault={updateDevice}>
                    <div class="row">
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputNameAddDevice" class="form-label">Device name</label>
                                <input type="text" class="form-control" id="inputNameAddDevice" bind:value="{addDevice.name}" required>
                            </div>
                        </div>
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputMacAddDevice" class="form-label">Mac address</label>
                                <input type="text" class="form-control" id="inputMacAddDevice" bind:value="{addDevice.mac}" pattern="^([0-9A-Fa-f]{'{'}2{'}'}[:-]){'{'}5{'}'}([0-9A-Fa-f]{'{'}2{'}'})|([0-9a-fA-F]{'{'}4{'}'}\\.[0-9a-fA-F]{'{'}4{'}'}\\.[0-9a-fA-F]{'{'}4{'}'})$" required>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputIpAddDevice" class="form-label">IP address</label>
                                <input type="text" class="form-control" id="inputIpAddDevice" bind:value="{addDevice.ip}" pattern="^([01]?\d\d?|2[0-4]\d|25[0-5])(?:\.(?:[01]?\d\d?|2[0-4]\d|25[0-5])){'{'}3{'}'}$" required>
                            </div>
                        </div>
                        <div class="col-sm">
                            <div class="mb-3">
                                <label for="inputNetmaskAddDevice" class="form-label">Netmask</label>
                                <input type="text" class="form-control" id="inputNetmaskAddDevice" bind:value="{addDevice.netmask}" pattern="^(((255\.){'{'}3{'}'}(255|254|252|248|240|224|192|128|0+))|((255\.){'{'}2{'}'}(255|254|252|248|240|224|192|128|0+)\.0)|((255\.)(255|254|252|248|240|224|192|128|0+)(\.0+){'{'}2{'}'})|((255|254|252|248|240|224|192|128|0+)(\.0+){'{'}3{'}'}))$" required>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="submit" form="addForm" class="btn btn-outline-success">Add device</button>
            </div>
        </div>
    </div>
</div>

<style lang="scss">
    .card {
        border-radius: 1em;
    }

    .btn-light {
        background-color: white;
        border-radius: 1em;

    }
</style>
