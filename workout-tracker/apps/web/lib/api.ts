const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function api(path: string, init?: RequestInit) {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(`${API_URL}${path}`, { ...init, headers: { ...headers, ...(init?.headers as any) } });
  if (!res.ok) throw new Error(await res.text());
  return res.status === 204 ? null : res.json();
}
