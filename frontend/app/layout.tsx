import React from 'react';
export const metadata = { title: 'Logistics' };
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es">
      <body style={{ fontFamily: 'sans-serif', margin: 0, padding: 20 }}>
        <nav style={{ display: 'flex', gap: 12 }}>
          <a href="/">Inicio</a>
          <a href="/client">Cliente</a>
          <a href="/admin">Admin</a>
        </nav>
        <hr />
        {children}
      </body>
    </html>
  );
}
