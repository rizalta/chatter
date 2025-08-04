<script lang="ts">
	import type { Message } from '@/types';
	import { cn } from '@/utils';
	import type { HTMLAttributes } from 'svelte/elements';
	import { Avatar, AvatarFallback } from '../avatar';

	export type ChatBubbleProps = {
		isUser: boolean;
		message: Message;
	} & HTMLAttributes<HTMLDivElement>;

	let { isUser, message, ...rest }: ChatBubbleProps = $props();

	const getInitials = (name: string) => {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.substring(0, 2)
			.toUpperCase();
	};

	export function formatTime(timestamp: string) {
		const date = new Date(timestamp);
		const now = new Date();

		const isToday =
			date.getFullYear() === now.getFullYear() &&
			date.getMonth() === now.getMonth() &&
			date.getDate() === now.getDate();

		if (isToday) {
			return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
		} else {
			return date.toLocaleDateString(); // localized date format
		}
	}
</script>

<div class={cn('flex items-center gap-2', isUser && 'flex-row-reverse')} {...rest}>
	<Avatar class="h-9 w-9">
		<AvatarFallback>
			{getInitials(message.fromName)}
		</AvatarFallback>
	</Avatar>
	<div class="flex flex-col gap-1">
		{#if !isUser}
			<p class="text-foreground/80 mb-1 text-xs font-semibold">{message.fromName}</p>
		{/if}
		<div
			class={cn(
				'group relative max-w-sm rounded-lg px-3 py-2 shadow-sm',
				isUser ? 'bg-primary text-primary-foreground rounded-br-none' : 'bg-muted rounded-bl-none'
			)}
		>
			<p class="text-sm whitespace-pre-wrap">{message.content}</p>
		</div>
		<span class={cn('text-muted-foreground text-sm', isUser ? 'text-end' : 'text-start')}
			>{formatTime(message.timestamp)}</span
		>
	</div>
</div>
