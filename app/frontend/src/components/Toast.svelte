<script>
    export let toast;
</script>

<div class="position-fixed bottom-0 end-0 p-3 toast-container">
    <div class="toast fade {toast.show ? 'show' : 'hide'} {toast.color}-bg" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="toast-header">
            <strong id="toast-title" class="me-auto">{toast.title}</strong>
            <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
        <div class="toast-body">
            {toast.message}
        </div>
    </div>
</div>

<style lang="scss">
    @import "../variables.scss";

    .toast-container {
		z-index: 11;
		border-radius: $toast-border-radius;

        div>.toast-body {
            color: $dark-lighter;
        }
	}

	.success-bg {
        background-color: var(--success);
    }

	.danger-bg {
        background-color: var(--danger);
    }
</style>
