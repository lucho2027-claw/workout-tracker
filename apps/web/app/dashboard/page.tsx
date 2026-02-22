"use client";

import { useEffect, useState } from "react";
import { api } from "../../lib/api";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { Input } from "../../components/ui/input";

export default function Dashboard() {
  const [exercises, setExercises] = useState<any[]>([]);
  const [name, setName] = useState("");

  async function load() {
    try {
      const data = await api("/exercises");
      setExercises(data);
    } catch {
      setExercises([]);
    }
  }

  async function addExercise() {
    await api("/exercises", { method: "POST", body: JSON.stringify({ name }) });
    setName("");
    await load();
  }

  useEffect(() => { load(); }, []);

  return (
    <Card>
      <h2 className="text-xl font-semibold">Dashboard</h2>
      <div className="mt-4 flex gap-2">
        <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Bench Press" className="flex-1" />
        <Button variant="success" onClick={addExercise}>Add</Button>
      </div>
      <h3 className="mt-6 text-lg font-medium">Exercises</h3>
      <ul className="mt-2 space-y-2">
        {exercises.map((e) => <li className="rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2" key={e.id}>{e.name}</li>)}
      </ul>
    </Card>
  );
}
