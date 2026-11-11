"use client";

import { useEffect, useState } from "react";
import { conversationsApi, contactsApi, inboxesApi, type Conversation, type Contact, type Inbox, type Message } from "@/lib/api";
import Button from "@/components/ui/Button";

const STATUS_TABS = ["", "open", "resolved", "pending"];
const STATUS_LABELS: Record<string, string> = { "": "All", open: "Open", resolved: "Resolved", pending: "Pending" };

export default function ConversationsPage() {
  const [convs, setConvs]         = useState<Conversation[]>([]);
  const [total, setTotal]         = useState(0);
  const [status, setStatus]       = useState("");
  const [loading, setLoading]     = useState(true);
  const [selected, setSelected]   = useState<Conversation | null>(null);
  const [messages, setMessages]   = useState<Message[]>([]);
  const [msgInput, setMsgInput]   = useState("");
  const [sending, setSending]     = useState(false);
  const [showNew, setShowNew]     = useState(false);
  const [contacts, setContacts]   = useState<Contact[]>([]);
  const [inboxes, setInboxes]     = useState<Inbox[]>([]);
  const [newForm, setNewForm]     = useState({ inbox_id: 0, contact_id: 0 });

  async function load(s = status) {
    setLoading(true);
    try {
      const res = await conversationsApi.list({ status: s, limit: 30 });
      setConvs(res.data ?? []);
      setTotal(res.total);
    } finally { setLoading(false); }
  }

  useEffect(() => { load(); }, []);

  async function openConv(cv: Conversation) {
    setSelected(cv);
    const res = await conversationsApi.listMessages(cv.id, { limit: 50 });
    setMessages(res.data ?? []);
  }

  async function sendMessage() {
    if (!selected || !msgInput.trim()) return;
    setSending(true);
    const msg = await conversationsApi.sendMessage(selected.id, { content: msgInput });
    setMessages((m) => [...m, msg]);
    setMsgInput("");
    setSending(false);
  }

  async function changeStatus(id: number, s: string) {
    await conversationsApi.updateStatus(id, { status: s });
    load();
    if (selected?.id === id) setSelected((prev) => prev ? { ...prev, status: s } : null);
  }

  async function openNewForm() {
    const [c, i] = await Promise.all([contactsApi.list({ limit: 100 }), inboxesApi.list()]);
    setContacts(c.data ?? []);
    setInboxes(i.data ?? []);
    setShowNew(true);
  }

  async function createConv() {
    await conversationsApi.create({ inbox_id: newForm.inbox_id, contact_id: newForm.contact_id });
    setShowNew(false); load();
  }

  const statusColor: Record<string, string> = {
    open: "bg-green-100 text-green-700", resolved: "bg-gray-100 text-gray-600",
    pending: "bg-yellow-100 text-yellow-700", snoozed: "bg-blue-100 text-blue-700",
  };

  return (
    <div className="h-[calc(100vh-8rem)] flex flex-col max-w-6xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <h1 className="text-xl font-bold text-gray-900">Conversations <span className="text-sm font-normal text-gray-400">({total})</span></h1>
        <Button onClick={openNewForm}>+ New</Button>
      </div>

      {/* Status tabs */}
      <div className="flex gap-1 mb-3">
        {STATUS_TABS.map((s) => (
          <button key={s} onClick={() => { setStatus(s); load(s); }}
            className={`px-3 py-1.5 rounded-lg text-xs font-medium transition ${status === s ? "bg-brand-600 text-white" : "bg-white text-gray-600 border border-gray-200 hover:border-brand-300"}`}>
            {STATUS_LABELS[s]}
          </button>
        ))}
      </div>

      <div className="flex gap-3 flex-1 min-h-0">
        {/* List */}
        <div className="w-72 shrink-0 bg-white rounded-xl shadow-sm ring-1 ring-gray-200 overflow-y-auto">
          {loading ? <p className="p-4 text-center text-gray-400 text-sm">Loading...</p>
          : convs.length === 0 ? <p className="p-4 text-center text-gray-400 text-sm">No conversations</p>
          : convs.map((cv) => (
            <button key={cv.id} onClick={() => openConv(cv)}
              className={`w-full text-left px-4 py-3 border-b border-gray-100 hover:bg-gray-50 transition ${selected?.id === cv.id ? "bg-brand-50" : ""}`}>
              <div className="flex items-start justify-between gap-2">
                <div className="min-w-0">
                  <p className="text-sm font-medium text-gray-900 truncate">{cv.contact?.name ?? `#${cv.id}`}</p>
                  <p className="text-xs text-gray-400 truncate">{cv.inbox?.name}</p>
                </div>
                <span className={`shrink-0 text-xs rounded-full px-1.5 py-0.5 font-medium ${statusColor[cv.status] ?? ""}`}>{cv.status}</span>
              </div>
            </button>
          ))}
        </div>

        {/* Chat */}
        {selected ? (
          <div className="flex-1 bg-white rounded-xl shadow-sm ring-1 ring-gray-200 flex flex-col min-h-0">
            {/* Chat header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100">
              <div>
                <p className="font-semibold text-gray-900">{selected.contact?.name ?? `#${selected.id}`}</p>
                <p className="text-xs text-gray-400">{selected.inbox?.name}</p>
              </div>
              <div className="flex gap-2">
                {["open", "resolved", "pending"].map((s) => (
                  <button key={s} onClick={() => changeStatus(selected.id, s)}
                    className={`text-xs px-2 py-1 rounded-lg border transition ${selected.status === s ? "bg-brand-600 text-white border-brand-600" : "border-gray-200 text-gray-600 hover:border-brand-300"}`}>
                    {s}
                  </button>
                ))}
              </div>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-4 space-y-2">
              {messages.length === 0 ? (
                <p className="text-center text-gray-400 text-sm py-8">No messages yet. Send the first one!</p>
              ) : messages.map((m) => (
                <div key={m.id} className={`flex ${m.sender_type === "agent" ? "justify-end" : "justify-start"}`}>
                  <div className={`max-w-xs rounded-2xl px-3 py-2 text-sm ${
                    m.sender_type === "agent"
                      ? "bg-brand-600 text-white rounded-br-sm"
                      : "bg-gray-100 text-gray-900 rounded-bl-sm"
                  }`}>
                    {m.content}
                    <p className={`text-xs mt-0.5 ${m.sender_type === "agent" ? "text-brand-200" : "text-gray-400"}`}>
                      {new Date(m.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                    </p>
                  </div>
                </div>
              ))}
            </div>

            {/* Input */}
            <div className="flex gap-2 p-3 border-t border-gray-100">
              <input value={msgInput} onChange={(e) => setMsgInput(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && !e.shiftKey && sendMessage()}
                placeholder="Type a message... (Enter to send)"
                className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500" />
              <Button onClick={sendMessage} isLoading={sending}>Send</Button>
            </div>
          </div>
        ) : (
          <div className="flex-1 bg-white rounded-xl shadow-sm ring-1 ring-gray-200 flex items-center justify-center">
            <p className="text-gray-400 text-sm">Select a conversation to start chatting</p>
          </div>
        )}
      </div>

      {/* New conversation modal */}
      {showNew && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl p-6 w-full max-w-sm shadow-xl">
            <h2 className="font-semibold text-gray-900 mb-4">New Conversation</h2>
            <div className="space-y-3">
              <div>
                <label className="text-sm font-medium text-gray-700 block mb-1">Contact</label>
                <select value={newForm.contact_id} onChange={(e) => setNewForm(p => ({ ...p, contact_id: +e.target.value }))}
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500">
                  <option value={0}>Select contact...</option>
                  {contacts.map((c) => <option key={c.id} value={c.id}>{c.name} — {c.phone}</option>)}
                </select>
              </div>
              <div>
                <label className="text-sm font-medium text-gray-700 block mb-1">Inbox</label>
                <select value={newForm.inbox_id} onChange={(e) => setNewForm(p => ({ ...p, inbox_id: +e.target.value }))}
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-brand-500">
                  <option value={0}>Select inbox...</option>
                  {inboxes.map((i) => <option key={i.id} value={i.id}>{i.name}</option>)}
                </select>
              </div>
            </div>
            <div className="flex gap-2 justify-end mt-4">
              <Button variant="secondary" onClick={() => setShowNew(false)}>Cancel</Button>
              <Button onClick={createConv} disabled={!newForm.contact_id || !newForm.inbox_id}>Create</Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
