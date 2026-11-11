import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "LinkChat",
  description: "Platform komunikasi pelanggan berbasis SaaS",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id">
      <body className="min-h-screen">{children}</body>
    </html>
  );
}
