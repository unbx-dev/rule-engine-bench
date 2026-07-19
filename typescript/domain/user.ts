import React from "react"; // ❌ domain層でreactをimport

export interface User {
  id: number;
  name: string;
  email: string;
}

// ❌ domain層でDate.nowを直接使う
export function getCurrentTimestamp(): number {
  return Date.now();
}

// ❌ domain層でReactのhookを使う
export function useSelectedUser(users: User[]) {
  const [selected, setSelected] = React.useState<User | null>(null);
  const timestamp = getCurrentTimestamp();
  return { selected, setSelected, timestamp };
}
