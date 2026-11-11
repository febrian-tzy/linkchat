"use client";

import { useEffect, useState } from "react";
import { inboxesApi, type Inbox } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

const CHANNEL_TYPES = ["whatsapp", "web", "email", "api"];

export default function InboxPage() {
  const [inboxes, setInboxes]   = useState<Inbox[]>([]);
  const [loading, setLoading]   = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [name, setName]         = useState("");
  const [type, setType]         = useState("whatsapp");
  const [saving, setSaving]     = useState(false);
  const [error, setError]       = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const res = await inboxesApi.list();
      setInboxes(res.data ?? []);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault(); setError(null); setSaving(true);
    try {
      await inboxesApi.create({ name, channel_type: type });
      setShowForm(false); setName(""); setType("whatsapp");
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally { setSaving(false); }
  }

  async function handleDelete(id: number) {
    if (!confirm("Delete this inbox?")) return;
    await inboxesApi.delete(id); load();
  }

  const channelIcon: Record<string, string> = {
    whatsapp: "💬", web: "🌐", email: "📧", api: "🔌",
  };

  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-900">Inboxes</h1>
        <Button onClick={() => setShowForm(true)}>+ New Inbox</Button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
          <h2 className="font-semibold text-gray-800 mb-3">New Inbox</h2>
          {error && <p className="text-sm text-red-600 mb-2">{error}</p>}
          <form onSubmit={handleCreate} className="space-y-3">
            <Input label="Inbox Name *" value={name} onChange={(e) => setName(e.target.value)} required />
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-1">Channel Type</label>
              <select value={type} onChange={(e) => setType(e.target.value)}
                className="w-full rounded-lg border border-gray-300 px-3 py-2.5 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500">
                {CHANNEL_TYPES.map((t) => (
                  <option key={t} value={t}>{channelIcon[t]} {t.charAt(0).toUpperCase() + t.slice(1)}</option>
                ))}
              </select>
            </div>
            <div className="flex gap-2 justify-end">
              <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>Cancel</Button>
              <Button type="submit" isLoading={saving}>Create</Button>
            </div>
          </form>
        </div>
      )}

      <div className="space-y-2">
        {loading ? <p className="text-center text-gray-400 py-8">Loading...</p>
        : inboxes.length === 0 ? <div className="bg-white rounded-xl p-8 text-center text-gray-400 shadow-sm ring-1 ring-gray-200">No inboxes yet. Create one to get started.</div>
        : inboxes.map((inbox) => (
          <div key={inbox.id} className="bg-white rounded-xl p-4 shadow-sm ring-1 ring-gray-200 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-brand-50 flex items-center justify-center text-xl">
                {channelIcon[inbox.channel_type] ?? "📥"}
              </div>
              <div>
                <p className="font-medium text-gray-900">{inbox.name}</p>
                <p className="text-xs text-gray-500 capitalize">{inbox.channel_type}</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <span className={`text-xs rounded-full px-2 py-0.5 font-medium ${inbox.is_enabled ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-500"}`}>
                {inbox.is_enabled ? "Active" : "Disabled"}
              </span>
              <button onClick={() => handleDelete(inbox.id)} className="text-xs text-red-500 hover:text-red-700">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
