import "./globals.css";
import Link from "next/link";
import { theme } from "../lib/theme";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <div className={theme.container}>
          <header className="mb-8 rounded-xl border border-zinc-800 bg-zinc-900/60 p-4">
            <h1 className="text-2xl font-bold tracking-tight">Workout Tracker</h1>
            <nav className="mt-3 flex gap-4 text-sm text-zinc-300">
              <Link className="hover:text-white" href="/">Home</Link>
              <Link className="hover:text-white" href="/login">Login</Link>
              <Link className="hover:text-white" href="/dashboard">Dashboard</Link>
              <Link className="hover:text-white" href="/workouts/new">Log workout</Link>
            </nav>
          </header>
          {children}
        </div>
      </body>
    </html>
  );
}
