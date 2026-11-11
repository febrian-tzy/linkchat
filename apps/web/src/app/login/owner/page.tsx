"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import { ownerApi, ownerTokenStorage } from "@/lib/api";

export default function LoginOwnerPage() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError]       = useState<string | null>(null);
  const [loading, setLoading]   = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (!username.trim() || !password) {
      setError("Username and password are required.");
      return;
    }

    setLoading(true);
    try {
      const result = await ownerApi.login(username, password);
      console.log("[Owner Login] success, token:", ownerTokenStorage.get()?.substring(0, 20));
      // Small delay to ensure localStorage is committed before navigation
      await new Promise((r) => setTimeout(r, 100));
      router.push("/owner/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-900 px-4 py-12">
      <div className="w-full max-w-sm">
        {/* Header */}
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-brand-600 shadow-lg">
            <span className="text-2xl font-extrabold text-white">L</span>
          </div>
          <h1 className="text-2xl font-extrabold text-white">Owner Portal</h1>
          <p className="mt-1 text-sm text-gray-400">Platform administration access</p>
        </div>

        {/* Card */}
        <div className="rounded-2xl bg-gray-800 p-8 shadow-xl ring-1 ring-gray-700">
          {error && (
            <div role="alert" className="mb-4 rounded-lg bg-red-900/50 border border-red-700 px-4 py-3 text-sm text-red-300">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4" noValidate>
            <div className="flex flex-col gap-1">
              <label htmlFor="username" className="text-sm font-medium text-gray-300">
                Username
              </label>
              <input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="owner"
                autoComplete="username"
                required
                className="block w-full rounded-lg border border-gray-600 bg-gray-700 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
              />
            </div>
            <div className="flex flex-col gap-1">
              <label htmlFor="password" className="text-sm font-medium text-gray-300">
                Password
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                autoComplete="current-password"
                required
                className="block w-full rounded-lg border border-gray-600 bg-gray-700 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 focus:ring-offset-gray-800 disabled:opacity-60 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
            >
              {loading && (
                <span className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
              )}
              Sign in as Owner
            </button>
          </form>
        </div>

        {/* Back link */}
        <p className="mt-4 text-center text-sm text-gray-500">
          Not an owner?{" "}
          <Link href="/login/tenant" className="text-brand-400 hover:text-brand-300 hover:underline">
            Tenant login →
          </Link>
        </p>
        <p className="mt-2 text-center">
          <Link href="/" className="text-xs text-gray-600 hover:text-gray-400">
            ← Back to home
          </Link>
        </p>
      </div>
    </div>
  );
}
