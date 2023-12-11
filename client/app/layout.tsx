import type { Metadata } from "next";
import { ThemeProvider } from "@/components/theme-provider";
import "./globals.css";
import { ModeToggle } from "@/components/mode-toggle";
import { Github } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Task Together",
  description: "Task Together",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <html lang="en" suppressHydrationWarning={true}>
        <head />
        <body>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            <header className="w-full flex justify-center border-b">
              <div className="w-full flex justify-between py-3 px-4 max-w-7xl">
                <h1 className="scroll-m-20 text-2xl font-semibold tracking-tight">
                  Task Together
                </h1>
                <div className="flex justify-between items-center gap-2">
                  <Button variant="outline" size="icon" asChild>
                    <Link href="https://github.com/michaelmagen/task-together">
                      <Github />
                    </Link>
                  </Button>
                  <ModeToggle />
                </div>
              </div>
            </header>
            {children}
          </ThemeProvider>
        </body>
      </html>
    </>
  );
}
