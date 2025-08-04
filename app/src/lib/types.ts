export interface User {
  id: string,
  username: string,
}

export interface Message {
  id: string;
  from: string;
  fromName: string;
  to: string;
  content: string;
  timestamp: string;
}

