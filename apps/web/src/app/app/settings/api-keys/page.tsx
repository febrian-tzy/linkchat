"use client";

import { useEffect, useState } from "react";
import { apiKeysApi, type APIKey } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

export default function APIKeysPage() {
  const [keys, setKeys]           = useState<APIKey[]>([]);
  const [loading, setLoading]     = useState(true);
  const [showForm, setShowForm]   = useState(false);
  const [name, setName]           = useState("");
  const [saving, setSaving]       = useState(false);
  const [newKey, setNewKey]       = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const res = await apiKeysApi.list();
      setKeys(res.data ?? []);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault(); setSaving(true);
    try {
      const res = await apiKeysApi.create({ name });
      setNewKey((res as APIKey & { key: string }).key);
      setShowForm(false); setName("");
      load();
    } finally { setSaving(false); }
  }

  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900">API Keys</h1>
          <p className="text-sm text-gray-500">Authenticate programmatic access</p>
        </div>
        <Button onClick={() => setShowForm(true)}>+ New Key</Button>
      </div>

      {/* One-time key display */}
      {newKey && (
        <div className="bg-green-50 border border-green-200 rounded-xl p-4">
          <p className="text-sm font-semibold text-green-800 mb-1">⚠️ Copy this key now — it won&apos;t be shown again!</p>
          <code className="block bg-white border border-green-200 rounded-lg px-3 py-2 text-sm font-mono break-all text-gray-900">{newKey}</code>
          <Button variant="secondary" className="mt-2 text-xs" onClick={() => { navigator.clipboard.writeText(newKey); }}>Copy to Clipboard</Button>
          <button onClick={() => setNewKey(null)} className="ml-3 text-xs text-green-700 hover:text-green-900">Dismiss</button>
        </div>
      )}

      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
          <h2 className="font-semibold text-gray-800 mb-3">New API Key</h2>
          <form onSubmit={handleCreate} className="flex gap-2">
            <Input placeholder="Key name (e.g. Production)" value={name} onChange={(e) => setName(e.target.value)} className="flex-1" required />
            <Button type="submit" isLoading={saving}>Generate</Button>
            <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>Cancel</Button>
          </form>
        </div>
      )}

      <div className="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 overflow-hidden">
        {loading ? <p className="p-8 text-center text-gray-400">Loading...</p>
        : keys.length === 0 ? <p className="p-8 text-center text-gray-400">No API keys yet.</p>
        : keys.map(k => (
          <div key={k.id} className="flex items-center gap-4 px-4 py-3 border-b border-gray-100 last:border-0">
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900">{k.name}</p>
              <code className="text-xs text-gray-400">{k.key_prefix}••••••••••••••••</code>
            </div>
            <span className={`text-xs rounded-full px-2 py-0.5 font-medium ${k.is_active ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-500"}`}>
              {k.is_active ? "Active" : "Revoked"}
            </span>
            <p className="text-xs text-gray-400 hidden sm:block">{new Date(k.created_at).toLocaleDateString()}</p>
            <button onClick={() => apiKeysApi.delete(k.id).then(load)} className="text-xs text-red-500 hover:text-red-700">Revoke</button>
          </div>
        ))}
      </div>

      <div className="bg-gray-50 rounded-xl p-4 text-sm text-gray-600">
        <p className="font-medium text-gray-800 mb-1">How to use</p>
        <p>Include the key in the <code className="bg-gray-200 rounded px-1">Authorization</code> header:</p>
        <code className="block mt-1 bg-gray-800 text-green-400 rounded-lg px-3 py-2 text-xs">
          Authorization: Bearer lc_your_api_key_here
        </code>
      </div>
    </div>
  );
}
