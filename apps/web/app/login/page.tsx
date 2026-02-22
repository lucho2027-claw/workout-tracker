"use client";

import { useState } from "react";
import { z } from "zod";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { FieldError } from "../../components/ui/field-error";
import { Input } from "../../components/ui/input";
import { theme } from "../../lib/theme";

const schema = z.object({
  email: z.string().email("Please enter a valid email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

  async function login() {
    const parsed = schema.safeParse({ email, password });
    if (!parsed.success) {
      const f = parsed.error.flatten().fieldErrors;
      setErrors({ email: f.email?.[0], password: f.password?.[0] });
      return;
    }
    setErrors({});

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
      <h2 className={theme.title}>Login</h2>
      <div className={theme.stack + " " + theme.sectionGap}>
        <div>
          <Input placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          <FieldError message={errors.email} />
        </div>
        <div>
          <Input placeholder="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
          <FieldError message={errors.password} />
        </div>
        <Button onClick={login}>Login</Button>
      </div>
      <p className="mt-3 text-sm text-zinc-300">{message}</p>
    </Card>
  );
}
