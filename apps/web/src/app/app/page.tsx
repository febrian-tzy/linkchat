"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { contactsApi, conversationsApi, broadcastsApi, inboxesApi } from "@/lib/api";

export default function DashboardPage() {
  const [stats, setStats] = useState({ contacts: 0, conversations: 0, broadcasts: 0, inboxes: 0 });

  useEffect(() => {
    Promise.allSettled([
      contactsApi.list({ limit: 1 }),
      conversationsApi.list({ limit: 1 }),
      broadcastsApi.list(),
      inboxesApi.list(),
    ]).then(([c, cv, b, i]) => {
      setStats({
        contacts:      c.status === "fulfilled"  ? c.value.total           : 0,
        conversations: cv.status === "fulfilled" ? cv.value.total          : 0,
        broadcasts:    b.status === "fulfilled"  ? b.value.data.length     : 0,
        inboxes:       i.status === "fulfilled"  ? i.value.data.length     : 0,
      });
    });
  }, []);

  const cards = [
    { label: "Contacts",      value: stats.contacts,      href: "/app/contacts",      icon: "👥", color: "bg-blue-50 text-blue-700" },
    { label: "Conversations", value: stats.conversations, href: "/app/conversations", icon: "💬", color: "bg-green-50 text-green-700" },
    { label: "Broadcasts",    value: stats.broadcasts,    href: "/app/broadcasts",    icon: "📢", color: "bg-orange-50 text-orange-700" },
    { label: "Inboxes",       value: stats.inboxes,       href: "/app/inbox",         icon: "📥", color: "bg-purple-50 text-purple-700" },
  ];

  const quickLinks = [
    { href: "/app/conversations", label: "New Conversation", icon: "💬" },
    { href: "/app/contacts",      label: "Add Contact",      icon: "➕" },
    { href: "/app/broadcasts",    label: "New Broadcast",    icon: "📢" },
    { href: "/app/team",          label: "Manage Team",      icon: "🧑‍🤝‍🧑" },
    { href: "/app/automations",   label: "Automations",      icon: "⚡" },
    { href: "/app/settings/webhooks", label: "Webhooks",     icon: "🔗" },
  ];

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-sm text-gray-500 mt-1">Overview of your LinkChat workspace</p>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {cards.map((c) => (
          <Link key={c.href} href={c.href}
            className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200 hover:shadow-md transition group">
            <div className={`inline-flex h-10 w-10 items-center justify-center rounded-lg text-xl ${c.color}`}>
              {c.icon}
            </div>
            <p className="mt-3 text-2xl font-bold text-gray-900">{c.value}</p>
            <p className="text-sm text-gray-500">{c.label}</p>
          </Link>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="bg-white rounded-xl p-5 shadow-sm ring-1 ring-gray-200">
        <h2 className="text-sm font-semibold text-gray-700 mb-3">Quick Actions</h2>
        <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
          {quickLinks.map((l) => (
            <Link key={l.href} href={l.href}
              className="flex items-center gap-2 rounded-lg border border-gray-200 px-3 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:border-brand-300 transition">
              <span>{l.icon}</span>
              <span>{l.label}</span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
