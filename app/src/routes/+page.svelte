<script lang="ts">
	import { onMount } from 'svelte';

	let messages = $state<string[]>([]);
	let newMessage = $state('');
	let socket: WebSocket;

	onMount(() => {
		socket = new WebSocket('ws://localhost:8080/ws');

		socket.onmessage = (event) => {
			messages = [...messages, event.data];
		};

		socket.onclose = () => {
			messages = [...messages, 'Connection closed'];
		};
	});

	function send() {
		if (newMessage && socket.readyState == WebSocket.OPEN) {
			socket.send(newMessage);
			newMessage = '';
		}
	}
</script>

<h1>Chatter</h1>
<ul>
	{#each messages as message, i (i)}
		<li>{message}</li>
	{/each}
</ul>

<form onsubmit={send}>
	<input type="text" bind:value={newMessage} />
	<button type="submit">Send</button>
</form>
