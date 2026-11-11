import Link from "next/link";

export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-brand-900 via-brand-700 to-brand-500 px-4">
      {/* Logo / Brand */}
      <div className="mb-10 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-white shadow-lg">
          <span className="text-3xl font-extrabold text-brand-600">L</span>
        </div>
        <h1 className="text-4xl font-extrabold tracking-tight text-white">
          LinkChat
        </h1>
        <p className="mt-2 text-brand-100">
          Platform komunikasi pelanggan untuk bisnis Anda
        </p>
      </div>

      {/* Action Cards */}
      <div className="w-full max-w-sm space-y-4">
        <Link
          href="/login/tenant"
          className="flex w-full items-center justify-center rounded-xl bg-white px-6 py-4 text-base font-semibold text-brand-700 shadow-md transition hover:bg-brand-50 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-brand-700"
        >
          🔐 Login Tenant
        </Link>

        <Link
          href="/register"
          className="flex w-full items-center justify-center rounded-xl bg-brand-500 px-6 py-4 text-base font-semibold text-white shadow-md ring-2 ring-white/30 transition hover:bg-brand-600 hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-brand-700"
        >
          ✨ Daftar Sekarang
        </Link>

        <Link
          href="/login/owner"
          className="flex w-full items-center justify-center rounded-xl bg-transparent px-6 py-4 text-sm font-medium text-brand-200 transition hover:text-white focus:outline-none"
        >
          Login Owner →
        </Link>
      </div>

      <p className="mt-12 text-xs text-brand-300">
        © {new Date().getFullYear()} LinkChat. All rights reserved.
      </p>
    </main>
  );
}
