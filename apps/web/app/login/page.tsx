"use client";

import { useState } from "react";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  async function login() {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password })
    });
    if (!res.ok) {
      setMessage("Login failed");
      return;
    }
    const data = await res.json();
    localStorage.setItem("token", data.token);
    setMessage("Logged in ✅");
  }

  return (
    <main className="max-w-md rounded-xl border border-zinc-800 bg-zinc-900/50 p-6">
      <h2 className="mb-4 text-xl font-semibold">Login</h2>
      <div className="space-y-3">
        <input className="w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input className="w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <button className="rounded-lg bg-blue-600 px-4 py-2 font-medium hover:bg-blue-500" onClick={login}>Login</button>
      </div>
      <p className="mt-3 text-sm text-zinc-300">{message}</p>
    </main>
  );
}
