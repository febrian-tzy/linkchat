"use client";

import { useEffect, useState } from "react";
import { broadcastsApi, contactsApi, inboxesApi, type Broadcast, type Contact, type Inbox } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

const STATUS_COLORS: Record<string, string> = {
  draft: "bg-gray-100 text-gray-600", scheduled: "bg-blue-100 text-blue-700",
  sending: "bg-yellow-100 text-yellow-700", sent: "bg-green-100 text-green-700",
  failed: "bg-red-100 text-red-700",
};

export default function BroadcastsPage() {
  const [broadcasts, setBroadcasts] = useState<Broadcast[]>([]);
  const [loading, setLoading]       = useState(true);
  const [showForm, setShowForm]     = useState(false);
  const [contacts, setContacts]     = useState<Contact[]>([]);
  const [inboxes, setInboxes]       = useState<Inbox[]>([]);
  const [form, setForm]             = useState({ name: "", message: "", inbox_id: 0, contact_ids: [] as number[] });
  const [saving, setSaving]         = useState(false);
  const [error, setError]           = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const res = await broadcastsApi.list();
      setBroadcasts(res.data ?? []);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  async function openForm() {
    const [c, i] = await Promise.all([contactsApi.list({ limit: 200 }), inboxesApi.list()]);
    setContacts(c.data ?? []);
    setInboxes(i.data ?? []);
    setShowForm(true);
  }

  function toggleContact(id: number) {
    setForm(p => ({
      ...p,
      contact_ids: p.contact_ids.includes(id)
        ? p.contact_ids.filter(c => c !== id)
        : [...p.contact_ids, id],
    }));
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault(); setError(null); setSaving(true);
    try {
      await broadcastsApi.create({ ...form });
      setShowForm(false);
      setForm({ name: "", message: "", inbox_id: 0, contact_ids: [] });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally { setSaving(false); }
  }

  return (
    <div className="max-w-4xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-900">Broadcasts</h1>
        <Button onClick={openForm}>+ New Broadcast</Button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200 space-y-4">
          <h2 className="font-semibold text-gray-800">New Broadcast</h2>
          {error && <p className="text-sm text-red-600">{error}</p>}
          <form onSubmit={handleCreate} className="space-y-3">
            <Input label="Broadcast Name *" value={form.name} onChange={(e) => setForm(p => ({ ...p, name: e.target.value }))} required />
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-1">Inbox *</label>
              <select value={form.inbox_id} onChange={(e) => setForm(p => ({ ...p, inbox_id: +e.target.value }))}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500" required>
                <option value={0}>Select inbox...</option>
                {inboxes.map(i => <option key={i.id} value={i.id}>{i.name}</option>)}
              </select>
            </div>
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-1">Message *</label>
              <textarea value={form.message} onChange={(e) => setForm(p => ({ ...p, message: e.target.value }))} rows={3} required
                className="w-full rounded-lg border border-gray-300 px-3 py-2.5 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500" />
            </div>
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-2">
                Recipients ({form.contact_ids.length} selected)
              </label>
              <div className="max-h-40 overflow-y-auto border border-gray-200 rounded-lg divide-y divide-gray-100">
                {contacts.map(c => (
                  <label key={c.id} className="flex items-center gap-3 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                    <input type="checkbox" checked={form.contact_ids.includes(c.id)} onChange={() => toggleContact(c.id)} className="rounded" />
                    <span className="text-sm text-gray-800">{c.name}</span>
                    <span className="text-xs text-gray-400">{c.phone}</span>
                  </label>
                ))}
              </div>
            </div>
            <div className="flex gap-2 justify-end">
              <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>Cancel</Button>
              <Button type="submit" isLoading={saving} disabled={form.contact_ids.length === 0}>Create Broadcast</Button>
            </div>
          </form>
        </div>
      )}

      <div className="space-y-2">
        {loading ? <p className="text-center text-gray-400 py-8">Loading...</p>
        : broadcasts.length === 0 ? (
          <div className="bg-white rounded-xl p-8 text-center text-gray-400 shadow-sm ring-1 ring-gray-200">No broadcasts yet.</div>
        ) : broadcasts.map(b => (
          <div key={b.id} className="bg-white rounded-xl p-4 shadow-sm ring-1 ring-gray-200">
            <div className="flex items-start justify-between gap-4">
              <div className="min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <p className="font-semibold text-gray-900">{b.name}</p>
                  <span className={`text-xs rounded-full px-2 py-0.5 font-medium ${STATUS_COLORS[b.status] ?? ""}`}>{b.status}</span>
                </div>
                <p className="text-sm text-gray-600 line-clamp-2">{b.message}</p>
                <div className="flex gap-4 mt-2 text-xs text-gray-400">
                  <span>Total: {b.total}</span>
                  <span>Sent: {b.sent_count}</span>
                  <span>Failed: {b.failed_count}</span>
                </div>
              </div>
              <button onClick={() => broadcastsApi.delete(b.id).then(load)} className="text-xs text-red-500 hover:text-red-700 shrink-0">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
