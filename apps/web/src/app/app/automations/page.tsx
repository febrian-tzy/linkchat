"use client";

import { useEffect, useState } from "react";
import { automationsApi, type Automation } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

const EVENTS = [
  "conversation_created", "message_received", "conversation_resolved",
  "conversation_assigned", "contact_created",
];

export default function AutomationsPage() {
  const [automations, setAutomations] = useState<Automation[]>([]);
  const [loading, setLoading]         = useState(true);
  const [showForm, setShowForm]       = useState(false);
  const [form, setForm]               = useState({ name: "", description: "", event: EVENTS[0], conditions: null, actions: [] });
  const [saving, setSaving]           = useState(false);
  const [error, setError]             = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const res = await automationsApi.list();
      setAutomations(res.data ?? []);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault(); setError(null); setSaving(true);
    try {
      await automationsApi.create({ ...form, actions: [{ type: "send_message", message: "Auto reply" }] });
      setShowForm(false);
      setForm({ name: "", description: "", event: EVENTS[0], conditions: null, actions: [] });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally { setSaving(false); }
  }

  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900">Automations</h1>
          <p className="text-sm text-gray-500">Automate actions based on events</p>
        </div>
        <Button onClick={() => setShowForm(true)}>+ New Automation</Button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
          <h2 className="font-semibold text-gray-800 mb-3">New Automation</h2>
          {error && <p className="text-sm text-red-600 mb-2">{error}</p>}
          <form onSubmit={handleCreate} className="space-y-3">
            <Input label="Name *" value={form.name} onChange={(e) => setForm(p => ({ ...p, name: e.target.value }))} required />
            <Input label="Description" value={form.description} onChange={(e) => setForm(p => ({ ...p, description: e.target.value }))} />
            <div>
              <label className="text-sm font-medium text-gray-700 block mb-1">Trigger Event *</label>
              <select value={form.event} onChange={(e) => setForm(p => ({ ...p, event: e.target.value }))}
                className="w-full rounded-lg border border-gray-300 px-3 py-2.5 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500">
                {EVENTS.map(ev => <option key={ev} value={ev}>{ev.replace(/_/g, " ")}</option>)}
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
        : automations.length === 0 ? (
          <div className="bg-white rounded-xl p-8 text-center text-gray-400 shadow-sm ring-1 ring-gray-200">No automations yet.</div>
        ) : automations.map(a => (
          <div key={a.id} className="bg-white rounded-xl p-4 shadow-sm ring-1 ring-gray-200 flex items-center justify-between gap-4">
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <p className="font-medium text-gray-900">{a.name}</p>
                <span className={`text-xs rounded-full px-2 py-0.5 font-medium ${a.is_active ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-500"}`}>
                  {a.is_active ? "Active" : "Inactive"}
                </span>
              </div>
              <p className="text-xs text-gray-400 mt-0.5">Trigger: <span className="text-gray-600">{a.event.replace(/_/g, " ")}</span> · Runs: {a.run_count}</p>
              {a.description && <p className="text-xs text-gray-500 mt-0.5">{a.description}</p>}
            </div>
            <div className="flex items-center gap-2 shrink-0">
              <button onClick={() => automationsApi.toggle(a.id, !a.is_active).then(load)}
                className={`text-xs px-2 py-1 rounded-lg border transition ${a.is_active ? "border-gray-200 text-gray-600 hover:bg-gray-50" : "border-green-200 text-green-600 hover:bg-green-50"}`}>
                {a.is_active ? "Disable" : "Enable"}
              </button>
              <button onClick={() => automationsApi.delete(a.id).then(load)} className="text-xs text-red-500 hover:text-red-700">Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
