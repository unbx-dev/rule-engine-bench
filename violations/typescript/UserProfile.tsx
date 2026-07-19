"use client";

import React, { useEffect, useState } from "react";
import { db } from "@/lib/db"; // VIOLATION: UI component imports DB client directly
import { createClient } from "@supabase/supabase-js";

// VIOLATION: any type
interface Props {
  userId: any;
  onUpdate: (data: any) => void;
}

// VIOLATION: console.log in production code
export function UserProfile({ userId, onUpdate }: Props) {
  const [user, setUser] = useState<any>(null); // VIOLATION: any

  useEffect(() => {
    console.log("UserProfile mounted, userId:", userId); // VIOLATION: console.log

    // VIOLATION: fetch directly without wrapper
    fetch(`/api/users/${userId}`)
      .then((r) => r.json())
      .then((data) => {
        console.log("fetched user:", data); // VIOLATION: console.log
        setUser(data);
      });

    // VIOLATION: Date.now() in component (domain logic)
    const ts = Date.now();
    console.log("render time:", ts); // VIOLATION: console.log
  }, [userId]);

  // VIOLATION: dangerouslySetInnerHTML
  const bio = user?.bio ?? "";
  return (
    <div>
      <h1>{user?.name}</h1>
      <div dangerouslySetInnerHTML={{ __html: bio }} />
    </div>
  );
}

// VIOLATION: as any
function coerceUser(raw: unknown) {
  return raw as any;
}

// VIOLATION: @ts-ignore
function legacyAdapter(input: string): number {
  // @ts-ignore
  return input;
}
