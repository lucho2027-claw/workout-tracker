"use client";

import { useEffect, useState } from "react";
import { api } from "../../lib/api";

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
    <main className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-6">
      <h2 className="text-xl font-semibold">Dashboard</h2>
      <div className="mt-4 flex gap-2">
        <input className="flex-1 rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" value={name} onChange={(e) => setName(e.target.value)} placeholder="Bench Press" />
        <button className="rounded-lg bg-emerald-600 px-4 py-2 font-medium hover:bg-emerald-500" onClick={addExercise}>Add</button>
      </div>
      <h3 className="mt-6 text-lg font-medium">Exercises</h3>
      <ul className="mt-2 space-y-2">
        {exercises.map((e) => <li className="rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2" key={e.id}>{e.name}</li>)}
      </ul>
    </main>
  );
}
