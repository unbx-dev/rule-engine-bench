"use client";

// VIOLATION: Client Component imports server-only module
import "server-only";
import { getSecretKey } from "@/lib/secrets";

export function ApiKeyDisplay() {
  const key = getSecretKey();
  return <span>{key}</span>;
}
