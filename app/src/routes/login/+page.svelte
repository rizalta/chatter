<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';
	import { auth, type LoginCredentials } from '@/stores/auth';

	let username = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	const handleLogin = async (e: Event) => {
		e.preventDefault();
		loading = true;
		error = '';

		const credentials: LoginCredentials = { username: username, password: password };

		try {
			const result = await auth.login(credentials);

			if (result.success) {
				username = '';
				password = '';

				await goto('/');
			} else {
				error = result.error || 'Something went wrong';
			}
		} catch (err) {
			const errMessage = err instanceof Error ? err.message : 'Something went wrong';
			error = errMessage;
		} finally {
			loading = false;
		}
	};
</script>

<form
	onsubmit={handleLogin}
	class="shadow-secondary flex w-1/3 max-w-lg flex-col gap-4 rounded-xl p-8 shadow-2xl"
>
	<h1 class="text-center text-lg font-bold">Login</h1>
	<div class="flex flex-col gap-2">
		<Label for="username">Your username</Label>
		<Input name="username" bind:value={username} required disabled={loading} />
	</div>
	<div class="flex flex-col gap-2">
		<Label for="password">Your password</Label>
		<Input name="password" type="password" bind:value={password} required disabled={loading} />
	</div>
	<Button type="submit" disabled={loading} class="mt-3">Login</Button>
	<p class="text-center text-sm">
		Don't have an account? <a href="/register" class="text-blue-800 dark:text-yellow-300"
			>Register</a
		>
		here
	</p>
	<p>{error}</p>
</form>
