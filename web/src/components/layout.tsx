import { ReactNode } from "react";
import Sidebar from "./sidebar";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="grid grid-cols-[200px_1fr]">
      <div className="h-screen bg-zinc-100"><Sidebar /></div>
      <div className="">{children}</div>
    </div>
  );
}