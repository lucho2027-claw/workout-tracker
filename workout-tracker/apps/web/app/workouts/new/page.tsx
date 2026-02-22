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
    <main>
      <h2>Log Workout</h2>
      <button onClick={createWorkout}>Start workout now</button>
      <p>Workout ID: {workoutId || "(none)"}</p>
      <input placeholder="Exercise UUID" value={exerciseId} onChange={(e) => setExerciseId(e.target.value)} /><br />
      <input type="number" value={setNumber} onChange={(e) => setSetNumber(Number(e.target.value))} /><br />
      <input type="number" value={reps} onChange={(e) => setReps(Number(e.target.value))} /><br />
      <input type="number" value={weight} onChange={(e) => setWeight(Number(e.target.value))} /><br />
      <button onClick={addSet} disabled={!workoutId}>Add set</button>
    </main>
  );
}
