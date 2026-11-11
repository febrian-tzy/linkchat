"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
import { authApi, tokenStorage } from "@/lib/api";
import type { MeResponse } from "@/types/auth";

const nav = [
  { href: "/app",              label: "Dashboard",    icon: "🏠" },
  { href: "/app/conversations",label: "Conversations",icon: "💬" },
  { href: "/app/contacts",     label: "Contacts",     icon: "👥" },
  { href: "/app/inbox",        label: "Inbox",        icon: "📥" },
  { href: "/app/broadcasts",   label: "Broadcast",    icon: "📢" },
  { href: "/app/automations",  label: "Automations",  icon: "⚡" },
  { href: "/app/team",         label: "Team",         icon: "🧑‍🤝‍🧑" },
  { href: "/app/settings/webhooks", label: "Webhooks", icon: "🔗" },
  { href: "/app/settings/api-keys", label: "API Keys", icon: "🔑" },
];

export default function AppShell({ children }: { children: React.ReactNode }) {
  const router   = useRouter();
  const pathname = usePathname();
  const [me, setMe]           = useState<MeResponse | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    if (!tokenStorage.get()) { router.replace("/login/tenant"); return; }
    authApi.me().then(setMe).catch(() => router.replace("/login/tenant"));
  }, [router]);

  async function handleLogout() {
    await authApi.logout();
    router.push("/login/tenant");
  }

  if (!me) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-brand-600 border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="flex h-screen overflow-hidden bg-gray-50">
      {/* Sidebar */}
      <aside className={`fixed inset-y-0 left-0 z-40 w-60 bg-gray-900 text-white flex flex-col transform transition-transform duration-200 md:relative md:translate-x-0 ${sidebarOpen ? "translate-x-0" : "-translate-x-full"}`}>
        {/* Logo */}
        <div className="flex items-center gap-2 px-4 py-4 border-b border-gray-700">
          <div className="h-8 w-8 rounded-lg bg-brand-600 flex items-center justify-center font-bold text-sm">
            {me.tenant.name.charAt(0)}
          </div>
          <div className="min-w-0">
            <p className="text-sm font-semibold truncate">{me.tenant.name}</p>
            <p className="text-xs text-gray-400 truncate">/{me.tenant.slug}</p>
          </div>
        </div>

        {/* Nav */}
        <nav className="flex-1 overflow-y-auto py-3 px-2 space-y-0.5">
          {nav.map((item) => {
            const active = item.href === "/app"
              ? pathname === "/app"
              : pathname.startsWith(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setSidebarOpen(false)}
                className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors ${
                  active
                    ? "bg-brand-600 text-white"
                    : "text-gray-400 hover:bg-gray-800 hover:text-white"
                }`}
              >
                <span>{item.icon}</span>
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>

        {/* User */}
        <div className="px-3 py-3 border-t border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <div className="h-8 w-8 rounded-full bg-gray-700 flex items-center justify-center text-xs font-bold">
              {me.full_name.charAt(0)}
            </div>
            <div className="min-w-0">
              <p className="text-xs font-medium truncate">{me.full_name}</p>
              <p className="text-xs text-gray-400 capitalize">{me.role}</p>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="w-full text-left text-xs text-gray-400 hover:text-white px-2 py-1 rounded hover:bg-gray-800 transition-colors"
          >
            🚪 Logout
          </button>
        </div>
      </aside>

      {/* Overlay (mobile) */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-30 bg-black/50 md:hidden" onClick={() => setSidebarOpen(false)} />
      )}

      {/* Main */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Mobile header */}
        <header className="md:hidden flex items-center gap-3 px-4 py-3 bg-white border-b border-gray-200">
          <button onClick={() => setSidebarOpen(true)} className="text-gray-600">
            ☰
          </button>
          <span className="text-sm font-semibold text-brand-600">LinkChat</span>
        </header>

        <main className="flex-1 overflow-y-auto p-4 md:p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
