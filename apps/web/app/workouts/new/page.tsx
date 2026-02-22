"use client";

import { useState } from "react";
import { z } from "zod";
import { api } from "../../../lib/api";
import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";
import { FieldError } from "../../../components/ui/field-error";
import { Input } from "../../../components/ui/input";
import { theme } from "../../../lib/theme";

const schema = z.object({
  exerciseId: z.string().uuid("Exercise UUID is invalid"),
  setNumber: z.number().int().min(1, "Set number must be at least 1"),
  reps: z.number().int().min(1, "Reps must be at least 1"),
  weight: z.number().min(0, "Weight must be 0 or higher"),
});

export default function NewWorkoutPage() {
  const [workoutId, setWorkoutId] = useState("");
  const [exerciseId, setExerciseId] = useState("");
  const [setNumber, setSetNumber] = useState(1);
  const [reps, setReps] = useState(10);
  const [weight, setWeight] = useState(20);
  const [errors, setErrors] = useState<Record<string, string | undefined>>({});

  async function createWorkout() {
    const now = new Date().toISOString();
    const data = await api("/workouts", { method: "POST", body: JSON.stringify({ performed_at: now }) });
    setWorkoutId(data.id);
  }

  async function addSet() {
    const parsed = schema.safeParse({ exerciseId, setNumber, reps, weight });
    if (!parsed.success) {
      const f = parsed.error.flatten().fieldErrors;
      setErrors({
        exerciseId: f.exerciseId?.[0],
        setNumber: f.setNumber?.[0],
        reps: f.reps?.[0],
        weight: f.weight?.[0],
      });
      return;
    }

    setErrors({});
    await api(`/workouts/${workoutId}/sets`, {
      method: "POST",
      body: JSON.stringify({ exercise_id: exerciseId, set_number: setNumber, reps, weight })
    });
    alert("Set saved");
  }

  return (
    <Card className="max-w-xl">
      <h2 className={theme.title}>Log Workout</h2>
      <Button className="mt-4" onClick={createWorkout}>Start workout now</Button>
      <p className="mt-3 text-sm text-zinc-300">Workout ID: {workoutId || "(none)"}</p>
      <div className="mt-4 space-y-3">
        <div>
          <Input placeholder="Exercise UUID" value={exerciseId} onChange={(e) => setExerciseId(e.target.value)} />
          <FieldError message={errors.exerciseId} />
        </div>
        <div>
          <Input type="number" value={setNumber} onChange={(e) => setSetNumber(Number(e.target.value))} />
          <FieldError message={errors.setNumber} />
        </div>
        <div>
          <Input type="number" value={reps} onChange={(e) => setReps(Number(e.target.value))} />
          <FieldError message={errors.reps} />
        </div>
        <div>
          <Input type="number" value={weight} onChange={(e) => setWeight(Number(e.target.value))} />
          <FieldError message={errors.weight} />
        </div>
      </div>
      <Button className="mt-4" variant="success" onClick={addSet} disabled={!workoutId}>Add set</Button>
    </Card>
  );
}
