<script>
    import { onMount } from 'svelte';
    import { pocketbase } from '@stores/pocketbase';
    import DeviceForm from '@components/DeviceForm.svelte';

    export let data;
    const id = data.params.id;
    let pb;
    let device = {
        id: id
    };

    onMount(async () => {
        pocketbase.subscribe((conn) => {
            pb = conn;
        });

        const result = await pb.collection('devices').getOne(id, {
            expand: 'ports'
        });
        if (!result.expand.ports) {
            result.expand.ports = [];
        }
        device = result;
    });
</script>

<div class="container">
    <DeviceForm bind:device mode="edit" />
</div>
