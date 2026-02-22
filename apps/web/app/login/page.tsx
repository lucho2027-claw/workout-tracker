"use client";

import { useState } from "react";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { Input } from "../../components/ui/input";

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
    <Card className="max-w-md">
      <h2 className="mb-4 text-xl font-semibold">Login</h2>
      <div className="space-y-3">
        <Input placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <Input placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <Button onClick={login}>Login</Button>
      </div>
      <p className="mt-3 text-sm text-zinc-300">{message}</p>
    </Card>
  );
}
