"use client";

import { useEffect, useState } from "react";
import { contactsApi, type Contact } from "@/lib/api";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

export default function ContactsPage() {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [total, setTotal]       = useState(0);
  const [search, setSearch]     = useState("");
  const [loading, setLoading]   = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm]         = useState({ name: "", phone: "", email: "", notes: "" });
  const [saving, setSaving]     = useState(false);
  const [error, setError]       = useState<string | null>(null);

  async function load(q = search) {
    setLoading(true);
    try {
      const res = await contactsApi.list({ search: q, limit: 20 });
      setContacts(res.data ?? []);
      setTotal(res.total);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, []);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setError(null); setSaving(true);
    try {
      await contactsApi.create(form);
      setShowForm(false);
      setForm({ name: "", phone: "", email: "", notes: "" });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error");
    } finally { setSaving(false); }
  }

  async function handleDelete(id: number) {
    if (!confirm("Delete this contact?")) return;
    await contactsApi.delete(id);
    load();
  }

  return (
    <div className="max-w-4xl mx-auto space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900">Contacts</h1>
          <p className="text-sm text-gray-500">{total} total contacts</p>
        </div>
        <Button onClick={() => setShowForm(true)}>+ Add Contact</Button>
      </div>

      {/* Search */}
      <Input placeholder="Search by name, phone, or email..." value={search}
        onChange={(e) => { setSearch(e.target.value); load(e.target.value); }} />

      {/* Form */}
      {showForm && (
        <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
          <h2 className="font-semibold text-gray-800 mb-3">New Contact</h2>
          {error && <p className="text-sm text-red-600 mb-3">{error}</p>}
          <form onSubmit={handleCreate} className="grid gap-3 sm:grid-cols-2">
            <Input label="Name *" value={form.name} onChange={(e) => setForm(p => ({ ...p, name: e.target.value }))} required />
            <Input label="Phone *" value={form.phone} onChange={(e) => setForm(p => ({ ...p, phone: e.target.value }))} required />
            <Input label="Email" type="email" value={form.email} onChange={(e) => setForm(p => ({ ...p, email: e.target.value }))} />
            <Input label="Notes" value={form.notes} onChange={(e) => setForm(p => ({ ...p, notes: e.target.value }))} />
            <div className="sm:col-span-2 flex gap-2 justify-end">
              <Button variant="secondary" type="button" onClick={() => setShowForm(false)}>Cancel</Button>
              <Button type="submit" isLoading={saving}>Save</Button>
            </div>
          </form>
        </div>
      )}

      {/* List */}
      <div className="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-gray-400">Loading...</div>
        ) : contacts.length === 0 ? (
          <div className="p-8 text-center text-gray-400">No contacts found.</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-xs font-medium text-gray-500 uppercase tracking-wide">
              <tr>
                <th className="px-4 py-3 text-left">Name</th>
                <th className="px-4 py-3 text-left">Phone</th>
                <th className="px-4 py-3 text-left hidden sm:table-cell">Email</th>
                <th className="px-4 py-3 text-left hidden md:table-cell">Status</th>
                <th className="px-4 py-3"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {contacts.map((c) => (
                <tr key={c.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{c.name}</td>
                  <td className="px-4 py-3 text-gray-600">{c.phone}</td>
                  <td className="px-4 py-3 text-gray-500 hidden sm:table-cell">{c.email || "—"}</td>
                  <td className="px-4 py-3 hidden md:table-cell">
                    <span className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${c.is_blocked ? "bg-red-100 text-red-700" : "bg-green-100 text-green-700"}`}>
                      {c.is_blocked ? "Blocked" : "Active"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button onClick={() => handleDelete(c.id)} className="text-xs text-red-500 hover:text-red-700">Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
