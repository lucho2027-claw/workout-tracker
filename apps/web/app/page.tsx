import Link from "next/link";

export default function Home() {
  return (
    <main>
      <p>Track exercises, sets, reps, and weight.</p>
      <ul>
        <li><Link href="/login">Login</Link></li>
        <li><Link href="/dashboard">Dashboard</Link></li>
        <li><Link href="/workouts/new">Log workout</Link></li>
      </ul>
    </main>
  );
}
