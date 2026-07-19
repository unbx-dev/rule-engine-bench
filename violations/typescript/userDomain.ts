// VIOLATION: domain layer imports react
import { useCallback } from "react";

export interface User {
  id: string;
  name: string;
  email: string;
  createdAt: number;
}

// VIOLATION: Date.now() in domain layer
export function createUser(name: string, email: string): User {
  return {
    id: crypto.randomUUID(),
    name,
    email,
    createdAt: Date.now(),
  };
}

// VIOLATION: domain uses React hook
export function useUserValidator(user: User) {
  return useCallback(() => {
    return user.email.includes("@");
  }, [user]);
}
