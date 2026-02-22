import { ReactNode } from "react";

export function Card({ children, className = "" }: { children: ReactNode; className?: string }) {
  return <section className={`rounded-xl border border-zinc-800 bg-zinc-900/50 p-6 ${className}`}>{children}</section>;
}
