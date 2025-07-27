<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import { ModeWatcher } from 'mode-watcher';
	import '../app.css';
	import { DarkMode } from '@/components';
	import { onMount } from 'svelte';

	import { auth, type AuthState } from '@/stores/auth';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let { children } = $props();

	onMount(() => {
		auth.init();
	});

	let authState: AuthState = $state({
		isAuthenticated: false,
		user: null,
		token: null,
		loading: true
	});

	auth.subscribe((state) => {
		authState = state;
	});

	const protectedRoutes = ['/chat'];
	const publicOnlyRoutes = ['/login', '/register'];

	$effect(() => {
		if (typeof window !== 'undefined' && !authState.loading) {
			const currentPath = page.url.pathname;

			if (protectedRoutes.some((route) => currentPath.startsWith(route)) || currentPath === '/') {
				if (!authState.isAuthenticated) goto('/login');
			}

			if (publicOnlyRoutes.includes(currentPath)) {
				if (authState.isAuthenticated) goto('/');
			}
		}
	});

	const handleLogout = () => {
		auth.logout();
	};
</script>

<ModeWatcher />
<div class="flex h-screen w-full flex-col">
	<nav class="flex w-full items-center justify-between border-b px-4 py-2">
		<Button variant="ghost" class="text-lg font-bold" href="/">Chatter</Button>
		<div class="flex gap-2">
			<DarkMode />
			{#if authState.isAuthenticated}
				<Button onclick={handleLogout}>Signout</Button>
			{/if}
		</div>
	</nav>

	<main class="flex flex-1 items-center justify-center overflow-y-auto p-4">
		{@render children()}
	</main>
</div>
