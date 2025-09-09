import React from 'react';
import Header from "../components/layout/Header";
import { AuthProvider } from "../components/features/auth/AuthContext";
import Toaster from "../components/ui/Toaster";
import 'bootstrap/dist/css/bootstrap.min.css'
import {Metadata} from "next";
import BootstrapClient from "../hooks/BootstrapClient";
import './globals.css';

export const metadata: Metadata = {
    title: 'Logistic App',
    description: "Logistic Application",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es">
      <body style={{ fontFamily: 'system-ui, -apple-system, Segoe UI, Roboto, sans-serif', margin: 0 }}>
        <AuthProvider>
          <Header />
          <main className="container-fluid p-3">
            {children}
          </main>
          <Toaster />
        </AuthProvider>
        <BootstrapClient />
      </body>
    </html>
  );
}
