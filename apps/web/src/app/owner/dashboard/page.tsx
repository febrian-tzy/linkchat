"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ownerApi, ownerTokenStorage, type OwnerMe, type PlatformStats, type Tenant } from "@/lib/api";

export default function OwnerDashboardPage() {
  const router  = useRouter();
  const [me, setMe]           = useState<OwnerMe | null>(null);
  const [stats, setStats]     = useState<PlatformStats | null>(null);
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = ownerTokenStorage.get();
    console.log("[Owner Dashboard] token:", token ? token.substring(0, 30) + "..." : "NULL");

    if (!token) {
      console.log("[Owner Dashboard] no token, redirecting");
      router.replace("/login/owner");
      return;
    }

    ownerApi.me()
      .then((meData) => {
        console.log("[Owner Dashboard] me:", meData);
        setMe(meData);
        return Promise.all([ownerApi.stats(), ownerApi.listTenants()]);
      })
      .then(([statsData, tenantsRes]) => {
        setStats(statsData);
        setTenants(tenantsRes.data ?? []);
      })
      .catch((err) => {
        console.error("[Owner Dashboard] error:", err?.message);
        router.replace("/login/owner");
      })
      .finally(() => setLoading(false));
  }, [router]);

  async function handleLogout() {
    await ownerApi.logout();
    router.push("/login/owner");
  }

  if (loading || !me || !stats) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-900">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-brand-600 border-t-transparent" />
      </div>
    );
  }

  const statCards = [
    { label: "Total Tenants",       value: stats.total_tenants,       color: "text-blue-400",   bg: "bg-blue-900/30"   },
    { label: "Active Tenants",      value: stats.active_tenants,      color: "text-green-400",  bg: "bg-green-900/30"  },
    { label: "Total Users",         value: stats.total_users,         color: "text-purple-400", bg: "bg-purple-900/30" },
    { label: "Total Contacts",      value: stats.total_contacts,      color: "text-orange-400", bg: "bg-orange-900/30" },
    { label: "Total Conversations", value: stats.total_conversations, color: "text-pink-400",   bg: "bg-pink-900/30"   },
  ];

  const statusColor: Record<string, string> = {
    active:    "bg-green-900/40 text-green-400",
    inactive:  "bg-gray-700 text-gray-400",
    suspended: "bg-red-900/40 text-red-400",
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* Top Nav */}
      <header className="sticky top-0 z-10 bg-gray-800 border-b border-gray-700">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <div className="h-8 w-8 rounded-lg bg-brand-600 flex items-center justify-center font-bold text-sm">L</div>
            <div>
              <p className="text-sm font-semibold">LinkChat Owner Portal</p>
              <p className="text-xs text-gray-400">Platform Administration</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <span className="hidden sm:block text-sm text-gray-300">
              {me.full_name}
            </span>
            <span className="rounded-full bg-brand-900/50 border border-brand-700 px-2.5 py-0.5 text-xs font-medium text-brand-300">
              Owner
            </span>
            <button
              onClick={handleLogout}
              className="rounded-lg border border-gray-600 px-3 py-1.5 text-xs text-gray-300 hover:border-gray-400 hover:text-white transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 py-8 space-y-8">

        {/* Welcome */}
        <div>
          <h1 className="text-2xl font-bold text-white">Platform Overview</h1>
          <p className="text-sm text-gray-400 mt-1">
            Welcome back, <span className="text-white font-medium">{me.full_name}</span>
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
          {statCards.map((s) => (
            <div key={s.label} className={`rounded-xl p-4 ring-1 ring-gray-700 ${s.bg}`}>
              <p className={`text-2xl font-bold ${s.color}`}>{s.value.toLocaleString()}</p>
              <p className="text-xs text-gray-400 mt-1">{s.label}</p>
            </div>
          ))}
        </div>

        {/* Tenants Table */}
        <div>
          <h2 className="text-lg font-semibold text-white mb-3">
            All Tenants ({tenants.length})
          </h2>
          <div className="rounded-xl bg-gray-800 ring-1 ring-gray-700 overflow-hidden">
            {tenants.length === 0 ? (
              <p className="p-8 text-center text-gray-500">No tenants registered yet.</p>
            ) : (
              <table className="w-full text-sm">
                <thead className="bg-gray-750 border-b border-gray-700">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wide">ID</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wide">Workspace</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wide hidden sm:table-cell">Slug</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wide">Status</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wide hidden md:table-cell">Created</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-700">
                  {tenants.map((t) => (
                    <tr key={t.id} className="hover:bg-gray-750 transition-colors">
                      <td className="px-4 py-3 text-gray-400 text-xs">#{t.id}</td>
                      <td className="px-4 py-3 font-medium text-white">{t.name}</td>
                      <td className="px-4 py-3 text-gray-400 font-mono text-xs hidden sm:table-cell">/{t.slug}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${statusColor[t.status] ?? "bg-gray-700 text-gray-400"}`}>
                          {t.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-gray-400 text-xs hidden md:table-cell">
                        {new Date(t.created_at).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>

        {/* Owner Info */}
        <div className="rounded-xl bg-gray-800 ring-1 ring-gray-700 p-5">
          <h2 className="text-sm font-semibold text-gray-300 mb-3">Owner Account</h2>
          <div className="grid gap-3 sm:grid-cols-3">
            {[
              { label: "Full Name", value: me.full_name },
              { label: "Username",  value: me.username },
              { label: "Email",     value: me.email },
            ].map((f) => (
              <div key={f.label} className="rounded-lg bg-gray-700/50 px-3 py-2">
                <p className="text-xs text-gray-400">{f.label}</p>
                <p className="text-sm font-medium text-white mt-0.5">{f.value}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Footer note */}
        <p className="text-center text-xs text-gray-600 pb-4">
          LinkChat Owner Portal · Platform Administration ·{" "}
          <Link href="/" className="hover:text-gray-400">Back to Home</Link>
        </p>
      </main>
    </div>
  );
}
