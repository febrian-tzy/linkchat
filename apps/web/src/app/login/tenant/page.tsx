"use client";

import { useState, FormEvent, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import { authApi } from "@/lib/api";

function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const justRegistered = params.get("registered") === "1";

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (!username.trim() || !password) {
      setError("Username dan password wajib diisi.");
      return;
    }

    setLoading(true);
    try {
      await authApi.login({ username, password });
      router.push("/app");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Terjadi kesalahan.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12">
      <div className="w-full max-w-sm">
        {/* Header */}
        <div className="mb-8 text-center">
          <Link
            href="/"
            className="inline-block text-2xl font-extrabold text-brand-600"
          >
            LinkChat
          </Link>
          <h2 className="mt-2 text-xl font-bold text-gray-900">
            Login Tenant
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Masuk ke workspace Anda
          </p>
        </div>

        {/* Registered success banner */}
        {justRegistered && (
          <div className="mb-4 rounded-lg bg-green-50 px-4 py-3 text-sm text-green-700">
            ✅ Registrasi berhasil! Silakan login.
          </div>
        )}

        {/* Card */}
        <div className="rounded-2xl bg-white p-8 shadow-sm ring-1 ring-gray-200">
          {error && (
            <div
              role="alert"
              className="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700"
            >
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4" noValidate>
            <Input
              id="username"
              name="username"
              label="Username"
              placeholder="Masukkan username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
              required
            />
            <Input
              id="password"
              name="password"
              type="password"
              label="Password"
              placeholder="Masukkan password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              required
            />

            <Button type="submit" isLoading={loading} className="w-full">
              Masuk
            </Button>
          </form>
        </div>

        <p className="mt-4 text-center text-sm text-gray-500">
          Belum punya workspace?{" "}
          <Link
            href="/register"
            className="font-medium text-brand-600 hover:underline"
          >
            Daftar sekarang
          </Link>
        </p>
      </div>
    </div>
  );
}

// Wrap in Suspense because useSearchParams() requires it in Next.js App Router
export default function LoginTenantPage() {
  return (
    <Suspense>
      <LoginForm />
    </Suspense>
  );
}
