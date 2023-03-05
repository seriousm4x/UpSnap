<script>
    import { pocketbase, authorizedStore } from '@stores/pocketbase';

    let username;
    let password;
    let isAdmin;
    let error = {
        hidden: true,
        title: '',
        msg: ''
    };

    async function login() {
        if (isAdmin) {
            $pocketbase.admins
                .authWithPassword(username, password)
                .then(() => {
                    authorizedStore.set(true);
                })
                .catch((err) => {
                    authorizedStore.set(false);
                    error.title = err.status;
                    error.msg = err.message;
                    error.hidden = false;
                });
        } else {
            $pocketbase
                .collection('users')
                .authWithPassword(username, password)
                .then(() => {
                    authorizedStore.set(true);
                })
                .catch((err) => {
                    authorizedStore.set(false);
                    error.title = err.status;
                    error.msg = err.message;
                    error.hidden = false;
                });
        }
    }
</script>

<div class="container text-dark-emphasis">
    <div class="row justify-content-center align-items-center">
        <div class="col text-center p-3 p-md-5">
            <img src="/favicon.png" alt="Logo" class="logo pe-none user-select-none" />
            <h2 class="text-dark-emphasis fw-bold text-center mt-3">Login</h2>
        </div>
    </div>
    <div class="row justify-content-center align-items-center p-2">
        <div class="col-md-6 col-lg-5 p-4 login rounded-5 shadow-sm">
            <form class="w-100" on:submit|preventDefault={login}>
                <div class="mb-3">
                    <label for="username" class="form-label">Username or email address</label>
                    <input
                        type="text"
                        class="form-control"
                        id="username"
                        placeholder="user@example.com"
                        bind:value={username}
                        required
                    />
                </div>
                <div class="mb-3">
                    <label for="password" class="form-label">Password</label>
                    <input
                        type="password"
                        class="form-control"
                        id="password"
                        placeholder="secret-password"
                        bind:value={password}
                        required
                    />
                </div>
                <div class="form-check form-switch mb-3">
                    <input
                        class="form-check-input"
                        type="checkbox"
                        role="switch"
                        id="flexSwitchCheckChecked"
                        bind:checked={isAdmin}
                    />
                    <label class="form-check-label" for="flexSwitchCheckChecked">Admin login</label>
                </div>
                <div class="callout callout-danger" class:d-none={error.hidden}>
                    <h5>Error {error.title}</h5>
                    <p class="m-0">{error.msg}</p>
                </div>
                <button class="btn btn-secondary w-100" type="submit">Login</button>
            </form>
            <p class="text-center m-0 mt-3">
                <!-- hacky way of linking a non existing route in svelte -->
                <a href="/_/" rel="external" class="text-muted cursor-pointer">Manage users</a>
            </p>
        </div>
    </div>
</div>

<style lang="scss">
    .logo {
        width: 100%;
        min-width: 50px;
        max-width: 100px;
    }

    .cursor-pointer {
        cursor: pointer;
    }
</style>
