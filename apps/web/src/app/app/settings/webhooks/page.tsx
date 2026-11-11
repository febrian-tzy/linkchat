"use client";

import { useEffect, useState } from "react";
import { webhooksApi, type Webhook } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

const AVAILABLE_EVENTS = [
  "conversation.created", "conversation.resolved", "conversation.assigned",
  "message.created", "contact.created", "broadcast.sent",
];

export default function WebhooksPage() {
  const [webhooks, setWebhooks]   = useState<Webhook[]>([]);
  const [loading, setLoading]     = useState(true);
  const [showForm, setShowForm]   = useState(false);
  const [form, setForm]           = useState({ url: "", secret: "", events: [] as string[] });
  const [saving, setSaving]       = useState(false);
  const [error, setError]         = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const res = await webhooksApi.list();
      setWebhooks(res.data ?? []);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  function toggleEvent(ev: string) {
    setForm(p => ({
      ...p,
      events: p.events.includes(ev) ? p.events.filter(e => e !== ev) : [...p.events, ev],
    }));
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault(); setError(null); setSaving(true);
    try {
      await webhooksApi.create(form);
      setShowForm(false);
      setForm({ url: "", secret: "", events: [] });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally { setSaving(false); }
  }

  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900">Webhooks</h1>
          <p className="text-sm text-gray-500">Receive real-time event notifications</p>
        </div>
        <Button onClick={() => setShowForm(true)}>+ Add Webhook</Button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
          <h2 className="font-semibold text-gray-800 mb-3">New Webhook</h2>
          {error && <p className="text-sm text-red-600 mb-2">{error}</p>}
          <form onSubmit={handleCreate} className="space-y-3">
            <Input label="Endpoint URL *" type="url" value={form.url} onChange={(e) => setForm(p => ({ ...p, url: e.target.value }))} placeholder="https://your-server.com/webhook" required />
            <Input label="Secret (optional)" value={form.secret} onChange={(e) => setForm(p => ({ ...p, secret: e.target.value }))} placeholder="Used to verify payload signature" />
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-2">Events ({form.events.length} selected)</label>
              <div className="grid grid-cols-2 gap-1.5">
                {AVAILABLE_EVENTS.map(ev => (
                  <label key={ev} className="flex items-center gap-2 text-sm text-gray-700 cursor-pointer hover:text-brand-600">
                    <input type="checkbox" checked={form.events.includes(ev)} onChange={() => toggleEvent(ev)} className="rounded" />
                    {ev}
                  </label>
                ))}
              </div>
            </div>
            <div className="flex gap-2 justify-end">
              <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>Cancel</Button>
              <Button type="submit" isLoading={saving} disabled={form.events.length === 0}>Save Webhook</Button>
            </div>
          </form>
        </div>
      )}

      <div className="space-y-2">
        {loading ? <p className="text-center text-gray-400 py-8">Loading...</p>
        : webhooks.length === 0 ? (
          <div className="bg-white rounded-xl p-8 text-center text-gray-400 shadow-sm ring-1 ring-gray-200">No webhooks configured.</div>
        ) : webhooks.map(wh => (
          <div key={wh.id} className="bg-white rounded-xl p-4 shadow-sm ring-1 ring-gray-200">
            <div className="flex items-start justify-between gap-4">
              <div className="min-w-0">
                <p className="text-sm font-mono text-gray-800 truncate">{wh.url}</p>
                <div className="flex flex-wrap gap-1 mt-2">
                  {(wh.events ?? []).map(ev => (
                    <span key={ev} className="text-xs bg-brand-50 text-brand-700 rounded px-1.5 py-0.5">{ev}</span>
                  ))}
                </div>
              </div>
              <button onClick={() => webhooksApi.delete(wh.id).then(load)} className="text-xs text-red-500 hover:text-red-700 shrink-0">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
