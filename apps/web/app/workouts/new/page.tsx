"use client";

import { useState } from "react";
import { api } from "../../../lib/api";
import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";
import { Input } from "../../../components/ui/input";

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
    <Card className="max-w-xl">
      <h2 className="text-xl font-semibold">Log Workout</h2>
      <Button className="mt-4" onClick={createWorkout}>Start workout now</Button>
      <p className="mt-3 text-sm text-zinc-300">Workout ID: {workoutId || "(none)"}</p>
      <div className="mt-4 space-y-3">
        <Input placeholder="Exercise UUID" value={exerciseId} onChange={(e) => setExerciseId(e.target.value)} />
        <Input type="number" value={setNumber} onChange={(e) => setSetNumber(Number(e.target.value))} />
        <Input type="number" value={reps} onChange={(e) => setReps(Number(e.target.value))} />
        <Input type="number" value={weight} onChange={(e) => setWeight(Number(e.target.value))} />
      </div>
      <Button className="mt-4" variant="success" onClick={addSet} disabled={!workoutId}>Add set</Button>
    </Card>
  );
}
