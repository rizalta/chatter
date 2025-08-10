<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { Input } from '@/components/ui/input';
	import { Button } from '@/components/ui/button';
	import { ChatBubble } from '@/components/ui/chat_bubble';
	import { auth, type AuthState } from '@/stores/auth';
	import type { Message, User, WSMessage } from '@/types';
	import { onDestroy, tick } from 'svelte';

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
	let activeUsers = $state<User[]>([]);
	let chatRoom = $state<HTMLDivElement>();

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
				const data = JSON.parse(event.data) as WSMessage;
				switch (data.type) {
					case 'chat': {
						messages = [...messages, data.data];
						break;
					}
					case 'presence': {
						const presence = data.data;
						if (presence.status === 'joined') {
							activeUsers = [...activeUsers, presence.user];
						} else if (presence.status === 'left') {
							activeUsers = activeUsers.filter((u) => u.id !== presence.user.id);
						}
						break;
					}
					case 'user_list': {
						activeUsers = data.data;
						break;
					}
				}
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
		} catch (err) {
			error = err instanceof Error ? err.message : 'Something went wrong';
		} finally {
			newMessage = '';
		}
	};

	const scrollToBottom = () => {
		if (chatRoom) chatRoom.scrollTop = chatRoom?.scrollHeight;
	};

	$effect(() => {
		if (messages.length > 0) {
			tick().then(() => {
				scrollToBottom();
			});
		}
	});
</script>

<main class="flex h-full w-4/5 lg:w-3/4">
	<span>{status} {error}</span>
	<div class="flex w-full flex-col justify-between gap-2 p-3 sm:w-3/4">
		<div class="flex w-full flex-col gap-2 overflow-auto px-3" bind:this={chatRoom}>
			{#each messages as message (message.id)}
				<ChatBubble isUser={authState.user?.id === message.from} {message} />
			{/each}
		</div>
		<form class="flex w-full gap-1" onsubmit={sendMessage}>
			<Input name="message" bind:value={newMessage} />
			<Button type="submit">Send</Button>
		</form>
	</div>
	<div class="bg-primary-foreground hidden w-1/4 rounded-lg p-3 shadow-2xl sm:block">
		<h1 class="text-secondary-foreground text-shadow-accent text-center text-lg font-semibold">
			Active Users
		</h1>
		<div class="flex flex-col gap-2 overflow-auto">
			{#each activeUsers as user (user.id)}
				{#if user.id !== authState.user?.id}
					<span class="bg-secondary rounded-md py-1 pl-5">
						<a href={`/message/${user.id}`} class="hover:underline">{user.username}</a>
					</span>
				{/if}
			{/each}
		</div>
	</div>
</main>
