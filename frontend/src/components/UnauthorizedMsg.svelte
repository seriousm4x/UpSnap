<script>
    import { pocketbase } from '@stores/pocketbase';

    let isUser = false;
    $: isUser = $pocketbase.authStore.baseModel?.collectionName === 'users';
</script>

{#if isUser}
    <div class="row justify-content-center align-items-center">
        <div class="col-auto">
            <div class="callout callout-danger m-0">
                Logged in as read only user. To make changes, please log in as admin.
            </div>
        </div>
    </div>
{/if}
