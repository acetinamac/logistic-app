import React from 'react';
import { AuthProvider } from "../components/features/auth/AuthContext";
import 'bootstrap/dist/css/bootstrap.min.css'
import {Metadata} from "next";
import BootstrapClient from "../hooks/BootstrapClient";
import './globals.css';
import AppShell from "../components/layout/AppShell";

export const metadata: Metadata = {
    title: 'Logistic App',
    description: "Logistic Application",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es">
      <body style={{ fontFamily: 'system-ui, -apple-system, Segoe UI, Roboto, sans-serif', margin: 0 }}>
        <AuthProvider>
          <AppShell>{children}</AppShell>
        </AuthProvider>
        <BootstrapClient />
      </body>
    </html>
  );
}
