<script>
    import { pocketbase, authorizedStore } from '@stores/pocketbase';

    let username;
    let password;
    let isAdmin;
    let error = {
        hidden: true,
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
                    error.msg = err;
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
                    error.msg = err;
                    error.hidden = false;
                });
        }
    }
</script>

<div class="container text-dark-emphasis h-100">
    <div class="row h-100 justify-content-center align-items-center">
        <div class="col-md-6 col-lg-5">
            <form class="w-100" on:submit|preventDefault={login}>
                <div class="mb-3">
                    <label for="username" class="form-label">Username or email address</label>
                    <input
                        type="text"
                        class="form-control"
                        id="username"
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
                        bind:value={password}
                        required
                    />
                </div>
                <div class="mb-3 form-check">
                    <input
                        type="checkbox"
                        class="form-check-input"
                        id="isAdminCheck"
                        bind:checked={isAdmin}
                    />
                    <label class="form-check-label" for="isAdminCheck">Admin Login</label>
                </div>
                <div class="callout callout-danger" class:d-none={error.hidden}>{error.msg}</div>
                <button class="btn btn-secondary w-100" type="submit">Login</button>
            </form>
        </div>
    </div>
</div>
