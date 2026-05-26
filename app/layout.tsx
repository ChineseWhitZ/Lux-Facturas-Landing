import type { Metadata } from "next";
import type { ReactNode } from "react";
import "./globals.css";

export const metadata: Metadata = {
  title: "Lux Facturas | POS Apple con facturacion SUNAT",
  description:
    "POS premium para negocios peruanos con hardware Apple, facturacion SUNAT, inventario, SIRE, instalacion y soporte.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  return (
    <html lang="es">
      <body>{children}</body>
    </html>
  );
}
