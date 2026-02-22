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
    <main>
      <h2>Login</h2>
      <input placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} /><br />
      <input placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} /><br />
      <button onClick={login}>Login</button>
      <p>{message}</p>
    </main>
  );
}
