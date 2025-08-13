export interface User {
  id: string,
  username: string,
}

export interface Message {
  id: string;
  from: string;
  fromName: string;
  to?: string;
  content: string;
  timestamp: string;
}

export type Status = 'joined' | 'left'

export type WSMessage =
  | { type: "chat", data: Message }
  | { type: "presence", data: { status: Status, user: User } }
  | { type: "user_list", data: User[] }
  | { type: "history", data: Message[] }
