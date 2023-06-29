<script>
    import { dev } from '$app/environment';
    import { pocketbase, authorizedStore, settings_public } from '@stores/pocketbase';

    let username;
    let password;
    let passwordConfirm;
    let isAdmin;
    let callout = {
        hidden: true,
        level: 'danger',
        title: '',
        msg: ''
    };

    $: if ($settings_public.setup_completed === false) {
        callout.title = '';
        callout.msg = 'No admin account has been set up. Please create one.';
        callout.level = 'info';
        callout.hidden = false;
    }

    async function login() {
        if (isAdmin) {
            $pocketbase.admins
                .authWithPassword(username, password)
                .then(() => {
                    authorizedStore.set(true);
                })
                .catch((err) => {
                    authorizedStore.set(false);
                    callout.level = 'danger';
                    callout.title = err.status;
                    callout.msg = err.message;
                    callout.hidden = false;
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
                    callout.level = 'danger';
                    callout.title = err.status;
                    callout.msg = err.message;
                    callout.hidden = false;
                });
        }
    }

    async function register() {
        if (password !== passwordConfirm) {
            callout.level = 'danger';
            callout.title = "Passwords don't match";
            callout.msg = 'Please try again';
            callout.hidden = false;
            return;
        } else {
            callout.hidden = true;
        }
        $pocketbase.admins
            .create({
                email: username,
                password: password,
                passwordConfirm: passwordConfirm
            })
            .then(() => {
                $pocketbase.admins
                    .authWithPassword(username, password)
                    .then(() => {
                        authorizedStore.set(true);
                    })
                    .catch((err) => {
                        return err;
                    });
            })
            .catch((err) => {
                authorizedStore.set(false);
                callout.level = 'danger';
                callout.title = err.status;
                callout.msg = err.message;
                callout.hidden = false;
            });
    }
</script>

<div class="container text-dark-emphasis">
    <div class="row justify-content-center align-items-center">
        <div class="col text-center p-3 p-md-5">
            <img
                src={$settings_public.favicon !== ''
                    ? `${dev ? $pocketbase.baseUrl : ''}/api/files/settings_public/${
                          $settings_public.id
                      }/${$settings_public.favicon}`
                    : '/gopher.svg'}
                alt={$settings_public.website_title ? $settings_public.website_title : 'UpSnap'}
                class="logo pe-none user-select-none"
            />
            <h2 class="text-dark-emphasis fw-bold text-center mt-3">
                {$settings_public.setup_completed === true ? 'Login' : 'Register'}
            </h2>
        </div>
    </div>
    <div class="row justify-content-center align-items-center p-2">
        <div class="col-md-6 col-lg-5 p-4 login rounded-5 shadow-sm">
            <div class="callout callout-{callout.level} mt-0" class:d-none={callout.hidden}>
                {#if callout.title !== ''}
                    <h5>{callout.title}</h5>
                {/if}
                {#if callout.msg !== ''}
                    <p class="m-0">{callout.msg}</p>
                {/if}
            </div>
            {#if $settings_public.setup_completed === true}
                <form class="w-100" on:submit|preventDefault={login}>
                    <div class="mb-3">
                        <label for="username" class="form-label">Username or email</label>
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
                        <label class="form-check-label" for="flexSwitchCheckChecked"
                            >Admin login</label
                        >
                    </div>
                    <button class="btn btn-secondary w-100" type="submit">Login</button>
                </form>
                <p class="text-center m-0 mt-3">
                    <!-- hacky way of linking a non existing route in svelte -->
                    <a href="/_/" rel="external" class="text-muted cursor-pointer">Manage users</a>
                </p>
            {:else}
                <form class="w-100" on:submit|preventDefault={register}>
                    <div class="mb-3">
                        <label for="username" class="form-label">Email</label>
                        <input
                            type="email"
                            class="form-control"
                            id="username"
                            placeholder="user@example.com"
                            bind:value={username}
                            required
                        />
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label"
                            >Password (min. 10 characters)</label
                        >
                        <input
                            type="password"
                            class="form-control"
                            id="password"
                            placeholder="secret-password"
                            minlength="10"
                            bind:value={password}
                            required
                        />
                    </div>
                    <div class="mb-3">
                        <label for="passwordConfirm" class="form-label">Password confirm</label>
                        <input
                            type="password"
                            class="form-control"
                            id="passwordConfirm"
                            placeholder="secret-password"
                            minlength="10"
                            bind:value={passwordConfirm}
                            required
                        />
                    </div>
                    <button class="btn btn-secondary w-100" type="submit">Create and login</button>
                </form>
            {/if}
        </div>
    </div>
</div>

<style lang="scss">
    .logo {
        width: 40%;
        min-width: 50px;
        max-width: 200px;
    }
</style>
