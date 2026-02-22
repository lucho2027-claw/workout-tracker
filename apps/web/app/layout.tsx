import "./globals.css";
import Link from "next/link";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <div className="mx-auto min-h-screen max-w-4xl px-4 py-8">
          <header className="mb-8 rounded-xl border border-zinc-800 bg-zinc-900/60 p-4">
            <h1 className="text-2xl font-bold">Workout Tracker</h1>
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
