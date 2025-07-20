<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';
</script>

<form class="shadow-secondary flex w-1/3 max-w-lg flex-col gap-4 rounded-xl p-8 shadow-2xl">
	<h1 class="text-center text-lg font-bold">Login</h1>
	<div class="flex flex-col gap-2">
		<Label for="username">Your username</Label>
		<Input name="username" />
	</div>
	<div class="flex flex-col gap-2">
		<Label for="password">Your password</Label>
		<Input name="password" type="password" />
	</div>
	<Button type="submit" class="mt-3">Login</Button>
	<p class="text-center text-sm">
		Don't have an account? <a href="/signup" class="text-blue-800 dark:text-yellow-300">Signup</a> here
	</p>
</form>
