<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import { ModeWatcher } from 'mode-watcher';
	import '../app.css';
	import { DarkMode } from '@/components';

	let { children } = $props();
</script>

<ModeWatcher />
<div class="flex h-screen w-full flex-col">
	<nav class="flex w-full items-center justify-between border-b px-4 py-2">
		<Button variant="ghost" class="text-lg font-bold" href="/">Chatter</Button>
		<DarkMode />
	</nav>

	<main class="flex flex-1 items-center justify-center overflow-y-auto p-4">
		{@render children()}
	</main>
</div>
