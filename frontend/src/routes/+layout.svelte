<script>
	import Navbar from '@components/Navbar.svelte';
</script>

<svelte:head>
    <script>
        if (document) {
            theme = window.matchMedia('(prefers-color-scheme: dark)').matches
                ? 'dark'
                : 'light';
            document.documentElement.setAttribute('data-bs-theme', theme);
        }
    </script>
</svelte:head>

<Navbar />

<slot />

<style lang="scss" global>
    @import '../scss/main.scss';
</style>
