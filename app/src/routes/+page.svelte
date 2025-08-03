<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { Input } from '@/components/ui/input';
	import { Button } from '@/components/ui/button';
	import { ChatBubble } from '@/components/ui/chat_bubble';
	import { auth, type AuthState } from '@/stores/auth';
	import type { Message } from '@/types';
	import { onDestroy } from 'svelte';

	type ConnectionStatus = 'connecting' | 'connected' | 'disconnected' | 'error';

	let authState = $state<AuthState>({
		isAuthenticated: false,
		user: null,
		token: null,
		loading: true
	});

	auth.subscribe((state) => {
		authState = state;
	});

	let ws: WebSocket | null = null;
	let status = $state<ConnectionStatus>('disconnected');
	let messages = $state<Message[]>([]);

	let newMessage = $state('');
	let error = $state('');

	const connectWS = () => {
		if (!authState.isAuthenticated) {
			console.error('No token available for ws connection');
			return;
		}

		if (ws && ws.readyState === WebSocket.OPEN) {
			console.log('ws already connected');
			return;
		}

		const wsURL = `${PUBLIC_API_URL}/chat/ws?token=${authState.token}`;

		ws = new WebSocket(wsURL);

		ws.onopen = () => {
			console.log('ws connected');
			status = 'connected';
		};

		ws.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data) as Message;
				messages = [...messages, data];
			} catch (error) {
				console.error('Error parsing ws message:', error);
			}
		};

		ws.onclose = () => {
			status = 'disconnected';
		};

		ws.onerror = (error) => {
			console.error('error connecting to ws: ', error);
			status = 'error';
		};
	};

	const disconnectWS = () => {
		if (ws) {
			ws.close();
			ws = null;
		}
	};

	$effect(() => {
		if (authState.isAuthenticated && !authState.loading) {
			connectWS();
		} else if (!authState.isAuthenticated) {
			disconnectWS();
		}
	});

	onDestroy(() => {
		disconnectWS();
	});

	const sendMessage = async (event: Event) => {
		event.preventDefault();
		error = '';
		try {
			const res = await fetch(`${PUBLIC_API_URL}/chat/chatroom`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${authState.token}`
				},
				body: JSON.stringify({ message: newMessage })
			});

			if (!res.ok) {
				throw new Error('Sending message failed');
			}
		} catch (error) {
			const err = error instanceof Error ? error.message : 'Something went wrong';
			console.log(err);
		}
	};
</script>

<main class="flex h-full w-4/5 lg:w-3/4">
	<div class="flex w-full flex-col justify-between gap-2 sm:w-3/4">
		<div class="w-full">
			{#each messages as message (message.id)}
				<ChatBubble isUser={authState.user?.id === message.from} {message} />
			{/each}
		</div>
		<form class="flex w-full gap-1" onsubmit={sendMessage}>
			<Input name="message" bind:value={newMessage} />
			<Button type="submit">Send</Button>
		</form>
	</div>
	<div class="hidden w-1/4 bg-blue-400 sm:block"></div>
</main>
