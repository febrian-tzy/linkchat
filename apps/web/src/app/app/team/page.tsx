"use client";

import { useEffect, useState } from "react";
import { usersApi, type TeamUser } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

const ROLES = ["admin", "agent", "supervisor"];

export default function TeamPage() {
  const [users, setUsers]       = useState<TeamUser[]>([]);
  const [loading, setLoading]   = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm]         = useState({ full_name: "", email: "", username: "", password: "", role: "agent" });
  const [saving, setSaving]     = useState(false);
  const [error, setError]       = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const res = await usersApi.list();
      setUsers(res.data ?? []);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  async function handleInvite(e: React.FormEvent) {
    e.preventDefault(); setError(null); setSaving(true);
    try {
      await usersApi.invite(form);
      setShowForm(false);
      setForm({ full_name: "", email: "", username: "", password: "", role: "agent" });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally { setSaving(false); }
  }

  const roleColor: Record<string, string> = {
    admin: "bg-purple-100 text-purple-700", agent: "bg-blue-100 text-blue-700", supervisor: "bg-orange-100 text-orange-700",
  };

  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900">Team Members</h1>
          <p className="text-sm text-gray-500">{users.length} members</p>
        </div>
        <Button onClick={() => setShowForm(true)}>+ Invite Member</Button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
          <h2 className="font-semibold text-gray-800 mb-3">Invite Team Member</h2>
          {error && <p className="text-sm text-red-600 mb-2">{error}</p>}
          <form onSubmit={handleInvite} className="grid gap-3 sm:grid-cols-2">
            <Input label="Full Name" value={form.full_name} onChange={(e) => setForm(p => ({ ...p, full_name: e.target.value }))} />
            <Input label="Email" type="email" value={form.email} onChange={(e) => setForm(p => ({ ...p, email: e.target.value }))} />
            <Input label="Username *" value={form.username} onChange={(e) => setForm(p => ({ ...p, username: e.target.value }))} required />
            <Input label="Password *" type="password" value={form.password} onChange={(e) => setForm(p => ({ ...p, password: e.target.value }))} required />
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-1">Role</label>
              <select value={form.role} onChange={(e) => setForm(p => ({ ...p, role: e.target.value }))}
                className="w-full rounded-lg border border-gray-300 px-3 py-2.5 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500">
                {ROLES.map(r => <option key={r} value={r}>{r.charAt(0).toUpperCase() + r.slice(1)}</option>)}
              </select>
            </div>
            <div className="sm:col-span-2 flex gap-2 justify-end">
              <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>Cancel</Button>
              <Button type="submit" isLoading={saving}>Invite</Button>
            </div>
          </form>
        </div>
      )}

      <div className="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 overflow-hidden">
        {loading ? <p className="p-8 text-center text-gray-400">Loading...</p>
        : users.map(u => (
          <div key={u.id} className="flex items-center gap-4 px-4 py-3 border-b border-gray-100 last:border-0 hover:bg-gray-50">
            <div className="h-9 w-9 rounded-full bg-brand-100 flex items-center justify-center text-sm font-bold text-brand-700">
              {u.full_name ? u.full_name.charAt(0).toUpperCase() : u.username.charAt(0).toUpperCase()}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900">{u.full_name || u.username}</p>
              <p className="text-xs text-gray-400">@{u.username} · {u.email}</p>
            </div>
            <div className="flex items-center gap-2">
              <span className={`text-xs rounded-full px-2 py-0.5 font-medium ${roleColor[u.role] ?? ""}`}>{u.role}</span>
              <span className={`text-xs rounded-full px-2 py-0.5 font-medium ${u.is_active ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-500"}`}>
                {u.is_active ? "Active" : "Inactive"}
              </span>
              <button onClick={() => usersApi.toggle(u.id, !u.is_active).then(load)}
                className="text-xs text-gray-400 hover:text-gray-700">{u.is_active ? "Deactivate" : "Activate"}</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
