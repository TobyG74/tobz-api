import { useState } from "react";
import { AnimatePresence, motion } from "motion/react";
import { ChevronRight, Download, Globe, KeyRound, LayoutDashboard, LogOut, Menu, Search, UserCog, X } from "lucide-react";
import type { ReactNode } from "react";
import clsx from "clsx";
import { useAuth } from "../../context/AuthContext";
import { GROUPS, type GroupId } from "../../lib/endpoints";
import { Logo } from "../primitives";
import { ThemeSwitcher } from "../ThemeSwitcher";
import { KeysManager } from "../KeysManager";
import { WhitelistManager } from "../WhitelistManager";
import { AccountPanel } from "./AccountPanel";
import { Overview } from "./Overview";
import { EndpointDoc } from "./EndpointDoc";
import { PanelHeader } from "./shared";

type View =
  | { kind: "overview" }
  | { kind: "keys" }
  | { kind: "whitelist" }
  | { kind: "account" }
  | { kind: "endpoint"; id: string };

const GROUP_ICON: Record<GroupId, ReactNode> = {
  downloader: <Download size={18} />,
  search: <Search size={18} />,
};

export function Dashboard() {
  const { user, logout } = useAuth();
  const [view, setView] = useState<View>({ kind: "overview" });
  const [open, setOpen] = useState<Set<GroupId>>(new Set(["downloader"]));
  const [navOpen, setNavOpen] = useState(false); // mobile drawer

  const go = (v: View) => {
    setView(v);
    setNavOpen(false);
  };

  const toggle = (g: GroupId) =>
    setOpen((prev) => {
      const next = new Set(prev);
      next.has(g) ? next.delete(g) : next.add(g);
      return next;
    });

  const openGroup = (g: GroupId) => {
    setOpen((prev) => new Set(prev).add(g));
    const first = GROUPS.find((x) => x.id === g)?.endpoints[0];
    if (first) go({ kind: "endpoint", id: first.id });
  };

  const currentEndpoint =
    view.kind === "endpoint" ? GROUPS.flatMap((g) => g.endpoints).find((e) => e.id === view.id) : undefined;

  const itemCls = (active: boolean) =>
    clsx(
      "flex w-full items-center gap-3 rounded-xl px-3 py-2 text-sm font-500 transition",
      active
        ? "bg-gradient-to-r from-blue/20 to-violet/20 text-paper ring-1 ring-violet/40"
        : "text-fog hover:bg-surface-2/60 hover:text-paper"
    );

  function renderNav() {
    return (
      <nav className="space-y-1">
        <button className={itemCls(view.kind === "overview")} onClick={() => go({ kind: "overview" })}>
          <LayoutDashboard size={18} /> Overview
        </button>
        <button className={itemCls(view.kind === "keys")} onClick={() => go({ kind: "keys" })}>
          <KeyRound size={18} /> API Keys
        </button>
        <button className={itemCls(view.kind === "whitelist")} onClick={() => go({ kind: "whitelist" })}>
          <Globe size={18} /> IP Whitelist
        </button>
        <button className={itemCls(view.kind === "account")} onClick={() => go({ kind: "account" })}>
          <UserCog size={18} /> Akun
        </button>

        {GROUPS.map((g) => {
          const expanded = open.has(g.id);
          return (
            <div key={g.id}>
              <button className={itemCls(false)} onClick={() => toggle(g.id)} aria-expanded={expanded}>
                {GROUP_ICON[g.id]}
                <span className="flex-1 text-left">{g.label}</span>
                <ChevronRight size={15} className={clsx("text-mist transition-transform", expanded && "rotate-90")} />
              </button>
              <AnimatePresence initial={false}>
                {expanded && (
                  <motion.div
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: "auto", opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    transition={{ duration: 0.2 }}
                    className="overflow-hidden"
                  >
                    <div className="ml-4 mt-1 space-y-0.5 border-l border-line pl-3">
                      {g.endpoints.map((e) => {
                        const active = view.kind === "endpoint" && view.id === e.id;
                        return (
                          <button
                            key={e.id}
                            onClick={() => go({ kind: "endpoint", id: e.id })}
                            className={clsx(
                              "block w-full rounded-lg px-3 py-1.5 text-left text-[13px] transition",
                              active ? "bg-surface-2/70 text-cyan" : "text-mist hover:text-paper"
                            )}
                          >
                            {e.label}
                          </button>
                        );
                      })}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          );
        })}
      </nav>
    );
  }

  return (
    <div className="min-h-screen">
      {/* Topbar */}
      <header className="sticky top-0 z-40 border-b border-line bg-canvas/80 backdrop-blur-md">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-3 px-4 py-3">
          <div className="flex items-center gap-2">
            <button
              onClick={() => setNavOpen(true)}
              className="grid h-9 w-9 place-items-center rounded-xl border border-line text-fog md:hidden"
              aria-label="Menu"
            >
              <Menu size={18} />
            </button>
            <Logo />
          </div>
          <div className="flex items-center gap-2">
            <ThemeSwitcher />
            <div className="flex items-center gap-2 rounded-xl border border-line bg-surface/60 py-1.5 pl-1.5 sm:pr-3">
              <span className="grid h-7 w-7 place-items-center rounded-lg bg-gradient-to-br from-blue to-violet text-xs font-700 text-white">
                {(user?.display_name || "U").slice(0, 1).toUpperCase()}
              </span>
              <span className="hidden max-w-[9rem] truncate pr-1 text-sm text-paper sm:block">{user?.display_name}</span>
            </div>
            <button
              onClick={logout}
              className="grid h-9 w-9 place-items-center rounded-xl border border-line text-fog transition hover:border-red-500/50 hover:text-red-300"
              aria-label="Keluar"
            >
              <LogOut size={16} />
            </button>
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-7xl gap-6 px-4 py-6 md:grid md:grid-cols-[250px_1fr]">
        {/* Desktop sidebar */}
        <aside className="hidden md:block">
          <div className="sticky top-20">{renderNav()}</div>
        </aside>

        {/* Content */}
        <motion.main
          key={view.kind === "endpoint" ? view.id : view.kind}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
          className="min-w-0"
        >
          {view.kind === "overview" && <Overview openGroup={openGroup} openKeys={() => go({ kind: "keys" })} />}
          {view.kind === "keys" && (
            <div>
              <PanelHeader icon={<KeyRound size={20} />} title="API Keys" desc="Buat, aktifkan, dan cabut API key. Pantau kuota harian." />
              <KeysManager />
            </div>
          )}
          {view.kind === "whitelist" && (
            <div>
              <PanelHeader icon={<Globe size={20} />} title="IP Whitelist" desc="Batasi login & pemakaian API key ke IP tertentu (maks 5)." />
              <WhitelistManager />
            </div>
          )}
          {view.kind === "account" && (
            <div>
              <PanelHeader icon={<UserCog size={20} />} title="Akun" desc="Kelola profil & ganti password." />
              <AccountPanel />
            </div>
          )}
          {currentEndpoint && (
            <EndpointDoc key={currentEndpoint.id} endpoint={currentEndpoint} goKeys={() => go({ kind: "keys" })} />
          )}
        </motion.main>
      </div>

      {/* Mobile drawer */}
      <AnimatePresence>
        {navOpen && (
          <motion.div
            className="fixed inset-0 z-50 md:hidden"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <div className="absolute inset-0 bg-void/70 backdrop-blur-sm" onClick={() => setNavOpen(false)} />
            <motion.aside
              className="glass absolute left-0 top-0 flex h-full w-72 max-w-[80%] flex-col p-4"
              initial={{ x: "-100%" }}
              animate={{ x: 0 }}
              exit={{ x: "-100%" }}
              transition={{ type: "spring", stiffness: 300, damping: 32 }}
            >
              <div className="mb-4 flex items-center justify-between">
                <Logo />
                <button
                  onClick={() => setNavOpen(false)}
                  className="rounded-lg p-2 text-mist hover:bg-surface-2 hover:text-paper"
                  aria-label="Tutup"
                >
                  <X size={18} />
                </button>
              </div>
              <div className="overflow-y-auto">{renderNav()}</div>
            </motion.aside>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
