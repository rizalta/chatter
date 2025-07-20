<script lang="ts">
	import { onMount } from 'svelte';
	import type { Message } from '$lib/types';

	let messages = $state<Message[]>([]);
	let newMessage = $state('');
	let socket: WebSocket;

	onMount(() => {
		socket = new WebSocket('ws://localhost:8080/ws');

		socket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data) as Message;
				messages = [...messages, data];
			} catch (err) {
				console.log(err);
			}
		};

		socket.onclose = () => {
			messages = [...messages, { content: 'Connection closed', to: 'chatroom' }];
		};
	});

	function send() {
		if (newMessage && socket.readyState == WebSocket.OPEN) {
			const message = { content: newMessage, to: 'chatroom' };
			socket.send(JSON.stringify(message));
			newMessage = '';
		}
	}
</script>

<h1>Chatter</h1>
<ul>
	{#each messages as message, i (i)}
		<li>{message.content}</li>
	{/each}
</ul>

<form onsubmit={send}>
	<input type="text" bind:value={newMessage} />
	<button type="submit">Send</button>
</form>
