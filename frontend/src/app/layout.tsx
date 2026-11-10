import React from 'react';

export const metadata = {
  title: 'E-commerce Grupo 3',
  description: 'Arquitectura Monolítica Go + React + PostgreSQL con Patrones Composite y Decorator',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="es">
      <body>{children}</body>
    </html>
  );
}
