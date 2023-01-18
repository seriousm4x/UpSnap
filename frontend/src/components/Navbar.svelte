<script>
	import { page } from '$app/stores';
</script>

<nav class="navbar navbar-expand">
	<div class="container-fluid">
		<a class="navbar-brand" href="/">
			<img src="favicon.png" alt="UpSnap" width="30" height="30" />
		</a>
		<div class="collapse navbar-collapse" id="navbarNav">
			<ul class="navbar-nav">
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/' ? true : false} href="/"
						>Home</a
					>
				</li>
				<li class="nav-item">
					<a
						class="nav-link active"
						class:active={$page.url.pathname === '/settings' ? true : false}
						href="/settings">Settings</a
					>
				</li>
			</ul>
		</div>
	</div>
</nav>
