"use client";

import { useEffect, useState } from "react";
import { z } from "zod";
import { api } from "../../lib/api";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { FieldError } from "../../components/ui/field-error";
import { Input } from "../../components/ui/input";
import { theme } from "../../lib/theme";

const schema = z.object({
  name: z.string().min(2, "Exercise name must be at least 2 characters"),
});

export default function Dashboard() {
  const [exercises, setExercises] = useState<any[]>([]);
  const [name, setName] = useState("");
  const [error, setError] = useState<string | undefined>();

  async function load() {
    try {
      const data = await api("/exercises");
      setExercises(data);
    } catch {
      setExercises([]);
    }
  }

  async function addExercise() {
    const parsed = schema.safeParse({ name });
    if (!parsed.success) {
      setError(parsed.error.flatten().fieldErrors.name?.[0]);
      return;
    }
    setError(undefined);
    await api("/exercises", { method: "POST", body: JSON.stringify({ name }) });
    setName("");
    await load();
  }

  useEffect(() => { load(); }, []);

  return (
    <Card>
      <h2 className={theme.title}>Dashboard</h2>
      <div className="mt-4 flex gap-2">
        <div className="flex-1">
          <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Bench Press" className="flex-1" />
          <FieldError message={error} />
        </div>
        <Button variant="success" onClick={addExercise}>Add</Button>
      </div>
      <h3 className="mt-6 text-lg font-medium tracking-tight">Exercises</h3>
      <ul className="mt-2 space-y-2">
        {exercises.map((e) => <li className="rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2" key={e.id}>{e.name}</li>)}
      </ul>
    </Card>
  );
}
