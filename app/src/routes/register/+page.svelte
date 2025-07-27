<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';

	let username = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	const handleRegister = async (e: Event) => {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await fetch('http://localhost:8080/api/user/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			if (!res.ok) {
				const err = await res.json();
				throw err;
			}
		} catch (e) {
			const err = e as Error;
			error = err.message;
		} finally {
			loading = false;
		}

		location.href = '/login';
	};
</script>

<form
	class="shadow-secondary flex w-1/3 max-w-lg flex-col gap-4 rounded-xl p-8 shadow-2xl"
	onsubmit={handleRegister}
>
	<h1 class="text-center text-lg font-bold">Register</h1>
	<div class="flex flex-col gap-2">
		<Label for="username">Your username</Label>
		<Input name="username" bind:value={username} />
	</div>
	<div class="flex flex-col gap-2">
		<Label for="password">Your password</Label>
		<Input name="password" type="password" bind:value={password} />
	</div>
	<Button type="submit" disabled={loading} class="mt-3">Register</Button>
	<p class="text-center text-sm">
		Already have an account? <a href="/login" class="text-blue-800 dark:text-yellow-300">Login</a> here
	</p>
	<p>{error}</p>
</form>
