"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import { authApi } from "@/lib/api";

interface FormState {
  workspace_name: string;
  full_name: string;
  email: string;
  username: string;
  password: string;
  confirm_password: string;
}

const initialState: FormState = {
  workspace_name: "",
  full_name: "",
  email: "",
  username: "",
  password: "",
  confirm_password: "",
};

export default function RegisterPage() {
  const router = useRouter();
  const [form, setForm] = useState<FormState>(initialState);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  }

  // Client-side quick validation before hitting the API
  function validate(): string | null {
    if (!form.workspace_name.trim()) return "Nama workspace wajib diisi.";
    if (!form.full_name.trim()) return "Nama lengkap wajib diisi.";
    if (!form.email.trim()) return "Email wajib diisi.";
    if (form.username.trim().length < 3)
      return "Username minimal 3 karakter.";
    if (form.password.length < 8)
      return "Password minimal 8 karakter.";
    if (form.password !== form.confirm_password)
      return "Password dan konfirmasi password tidak cocok.";
    return null;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    const validationError = validate();
    if (validationError) {
      setError(validationError);
      return;
    }

    setLoading(true);
    try {
      await authApi.register({
        workspace_name: form.workspace_name,
        full_name: form.full_name,
        email: form.email,
        username: form.username,
        password: form.password,
        confirm_password: form.confirm_password,
      });
      router.push("/login/tenant?registered=1");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Terjadi kesalahan.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="mb-8 text-center">
          <Link
            href="/"
            className="inline-block text-2xl font-extrabold text-brand-600"
          >
            LinkChat
          </Link>
          <h2 className="mt-2 text-xl font-bold text-gray-900">
            Daftar Workspace Baru
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Isi form berikut untuk membuat workspace dan akun admin pertama
          </p>
        </div>

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
              id="workspace_name"
              name="workspace_name"
              label="Nama Workspace"
              placeholder="Contoh: Damian Network"
              value={form.workspace_name}
              onChange={handleChange}
              autoComplete="organization"
              required
            />
            <Input
              id="full_name"
              name="full_name"
              label="Nama Lengkap"
              placeholder="Nama lengkap Anda"
              value={form.full_name}
              onChange={handleChange}
              autoComplete="name"
              required
            />
            <Input
              id="email"
              name="email"
              type="email"
              label="Email"
              placeholder="email@perusahaan.com"
              value={form.email}
              onChange={handleChange}
              autoComplete="email"
              required
            />
            <Input
              id="username"
              name="username"
              label="Username"
              placeholder="username (min. 3 karakter)"
              value={form.username}
              onChange={handleChange}
              autoComplete="username"
              required
            />
            <Input
              id="password"
              name="password"
              type="password"
              label="Password"
              placeholder="Min. 8 karakter"
              value={form.password}
              onChange={handleChange}
              autoComplete="new-password"
              required
            />
            <Input
              id="confirm_password"
              name="confirm_password"
              type="password"
              label="Konfirmasi Password"
              placeholder="Ulangi password"
              value={form.confirm_password}
              onChange={handleChange}
              autoComplete="new-password"
              required
            />

            <Button
              type="submit"
              isLoading={loading}
              className="mt-2 w-full"
            >
              Buat Workspace
            </Button>
          </form>
        </div>

        <p className="mt-4 text-center text-sm text-gray-500">
          Sudah punya akun?{" "}
          <Link
            href="/login/tenant"
            className="font-medium text-brand-600 hover:underline"
          >
            Login di sini
          </Link>
        </p>
      </div>
    </div>
  );
}
