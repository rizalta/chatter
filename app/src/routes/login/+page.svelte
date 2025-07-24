<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';

	let username = $state('');
	let password = $state('');
	let error = $state('');

	const handleLogin = async (e: Event) => {
		e.preventDefault();

		error = '';

		try {
			const res = await fetch('http://localhost:8080/api/user/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			if (!res.ok) {
				const errData = await res.json();
				throw new Error(errData.message || 'Login failed');
			}

			console.log(res);

			const data = await res.json();
			localStorage.setItem('jwt_token', data.token);
			window.location.href = '/chat';
		} catch (e) {
			error = (e as Error).message;
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
		<Input name="username" bind:value={username} />
	</div>
	<div class="flex flex-col gap-2">
		<Label for="password">Your password</Label>
		<Input name="password" type="password" bind:value={password} />
	</div>
	<Button type="submit" class="mt-3">Login</Button>
	<p class="text-center text-sm">
		Don't have an account? <a href="/signup" class="text-blue-800 dark:text-yellow-300">Signup</a> here
	</p>
	<p>{error}</p>
</form>
