<script>
    export let visitors;

    let label = visitors === 1 ? "Visitor" : "Visitors";
    let icon = visitors === 1 ? "fa-user" : "fa-user-group";

</script>

<nav class="navbar navbar-expand-lg bg-light">
<div class="container-fluid">
    <a class="navbar-brand" href="/">
        <img src="/favicon.png" alt="Logo" width="24" height="24">
    </a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
    <div class="navbar-nav">
        <a class="nav-link link-dark" href="/">Devices</a>
        <a class="nav-link link-dark" href="/settings">Settings</a>
    </div>
    </div>
    <span class="navbar-text">
        <div class="card border-0">
            <div class="card-body py-2 px-3">
                {visitors} {label}<i class="ms-1 fa-solid {icon} text-muted"></i>
            </div>
        </div>
    </span>
</div>
</nav>

<style lang="scss">
    .card {
        border-radius: 1em;
    }
</style>
