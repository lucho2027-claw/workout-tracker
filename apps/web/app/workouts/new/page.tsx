"use client";

import { useState } from "react";
import { api } from "../../../lib/api";

export default function NewWorkoutPage() {
  const [workoutId, setWorkoutId] = useState("");
  const [exerciseId, setExerciseId] = useState("");
  const [setNumber, setSetNumber] = useState(1);
  const [reps, setReps] = useState(10);
  const [weight, setWeight] = useState(20);

  async function createWorkout() {
    const now = new Date().toISOString();
    const data = await api("/workouts", { method: "POST", body: JSON.stringify({ performed_at: now }) });
    setWorkoutId(data.id);
  }

  async function addSet() {
    await api(`/workouts/${workoutId}/sets`, {
      method: "POST",
      body: JSON.stringify({ exercise_id: exerciseId, set_number: setNumber, reps, weight })
    });
    alert("Set saved");
  }

  return (
    <main className="max-w-xl rounded-xl border border-zinc-800 bg-zinc-900/50 p-6">
      <h2 className="text-xl font-semibold">Log Workout</h2>
      <button className="mt-4 rounded-lg bg-blue-600 px-4 py-2 font-medium hover:bg-blue-500" onClick={createWorkout}>Start workout now</button>
      <p className="mt-3 text-sm text-zinc-300">Workout ID: {workoutId || "(none)"}</p>
      <div className="mt-4 space-y-3">
        <input className="w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" placeholder="Exercise UUID" value={exerciseId} onChange={(e) => setExerciseId(e.target.value)} />
        <input className="w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" type="number" value={setNumber} onChange={(e) => setSetNumber(Number(e.target.value))} />
        <input className="w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" type="number" value={reps} onChange={(e) => setReps(Number(e.target.value))} />
        <input className="w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2" type="number" value={weight} onChange={(e) => setWeight(Number(e.target.value))} />
      </div>
      <button className="mt-4 rounded-lg bg-emerald-600 px-4 py-2 font-medium hover:bg-emerald-500 disabled:opacity-50" onClick={addSet} disabled={!workoutId}>Add set</button>
    </main>
  );
}
