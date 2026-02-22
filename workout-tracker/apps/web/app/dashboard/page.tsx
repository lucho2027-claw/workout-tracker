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
    <main>
      <h2>Dashboard</h2>
      <p>Add exercise</p>
      <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Bench Press" />
      <button onClick={addExercise}>Add</button>
      <h3>Exercises</h3>
      <ul>
        {exercises.map((e) => <li key={e.id}>{e.name}</li>)}
      </ul>
    </main>
  );
}
