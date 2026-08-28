import "./globals.css";
import type { ReactNode } from "react";

export const metadata = {
  title: "EIP Training Platform",
  description: "Local MVP console for the EIP training loop",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
