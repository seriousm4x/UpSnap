<script>
    import { dev } from '$app/environment';
    import { parseISO, formatDistance } from 'date-fns';
    import {
        faPowerOff,
        faEllipsisVertical,
        faCircle,
        faCircleUp,
        faCircleDown,
        faLock
    } from '@fortawesome/free-solid-svg-icons';
    import Fa from 'svelte-fa';
    import { sortPorts } from '../sorts';

    export let device;
    export let now;

    function shutdown() {
        fetch(`${dev ? 'http://localhost:8090' : ''}/api/upsnap/shutdown/${device.id}`);
    }

    function wake() {
        fetch(`${dev ? 'http://localhost:8090' : ''}/api/upsnap/wake/${device.id}`);
    }
</script>

<div class="col-xs-12 col-sm-6 col-md-4 col-lg-3 g-4">
    <div
        class="card border-0 rounded-5 px-3 py-2 shadow-sm"
        class:offline={device.status == 'offline' ? true : false}
        class:online={device.status == 'online' ? true : false}
        class:pending={device.status == 'pending' ? true : false}
    >
        <div class="card-body">
            <div class="row">
                <div class="col-auto me-auto">
                    {#if device.status === 'offline'}
                        <span class="text-danger" on:click={() => wake()} on:keydown={() => wake()}>
                            <div role="button">
                                <Fa icon={faPowerOff} class="fs-4 power-hover" />
                            </div>
                        </span>
                    {:else if device.status === 'online'}
                        {#if device.shutdown_cmd !== ''}
                            <span
                                class="text-success"
                                on:click={() => shutdown()}
                                on:keydown={() => shutdown()}
                            >
                                <div role="button">
                                    <Fa icon={faPowerOff} class="fs-4 power-hover" />
                                </div>
                            </span>
                        {:else}
                            <span class="text-success">
                                <div>
                                    <Fa icon={faPowerOff} class="fs-4" />
                                </div>
                            </span>
                        {/if}
                    {:else if device.status === 'pending'}
                        <div
                            class="spinner-border text-warning"
                            style="width: 23px;height: 23px;"
                            role="status"
                        >
                            <span class="visually-hidden">Loading...</span>
                        </div>
                    {/if}
                </div>
                <div class="col-auto fs-5">
                    <a
                        href="/device/{device.id}"
                        class="text-reset text-center ellipsis power-hover"
                    >
                        <Fa icon={faEllipsisVertical} />
                    </a>
                </div>
            </div>
            {#if device.link}
                <p class="m-0 fw-bold fs-5">
                    <a class="text-reset" target="_blank" rel="noreferrer" href={device.link}
                        >{device.name}</a
                    >
                </p>
            {:else}
                <p class="m-0 fw-bold fs-5">{device.name}</p>
            {/if}
            <p class="text-muted mb-2">{device.ip}</p>
            {#if device?.expand?.ports}
                <div class="mb-2">
                    {#each device.expand.ports.sort(sortPorts) as port}
                        <p class="m-0">
                            <Fa
                                icon={faCircle}
                                class="me-1 {port.status ? 'port-up' : 'port-down'}"
                            />{port.name}
                            <span class="text-muted">({port.number})</span>
                        </p>
                    {/each}
                </div>
            {/if}
            {#if device.wake_cron !== ''}
                <span
                    class="badge rounded-pill {device.wake_cron_enabled
                        ? 'text-bg-success'
                        : 'text-bg-danger'}"
                    data-toggle="tooltip"
                    title="Wake cron {device.wake_cron_enabled ? 'enabled' : 'disabled'}"
                >
                    <Fa icon={faCircleUp} class="me-1" />
                    {device.wake_cron}
                </span>
            {/if}
            {#if device.shutdown_cron !== ''}
                <span
                    class="badge rounded-pill {device.shutdown_cron_enabled
                        ? 'text-bg-success'
                        : 'text-bg-danger'}"
                    data-toggle="tooltip"
                    title="Shutdown cron {device.shutdown_cron_enabled ? 'enabled' : 'disabled'}"
                >
                    <Fa icon={faCircleDown} class="me-1" />
                    {device.shutdown_cron}
                </span>
            {/if}
            {#if device.password !== ''}
                <span
                    class="badge rounded-pill text-bg-secondary"
                    data-toggle="tooltip"
                    title="Wake password set"
                >
                    <Fa icon={faLock} class="mx-1" />
                </span>
            {/if}
            <p
                class="text-muted m-0"
                data-toggle="tooltip"
                title="Last status change: {device.updated}"
            >
                {formatDistance(parseISO(device.updated), now, {
                    includeSeconds: true,
                    addSuffix: true
                })}
            </p>
        </div>
    </div>
</div>
