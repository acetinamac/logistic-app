"use client";
import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { API_BASE } from "../../../lib/constants";

export type AuthState = {
  token: string | null;
  email: string | null;
  isAuthenticated: boolean;
};

export type Toast = { id: number; type: "success" | "danger" | "info" | "warning"; message: string; timeout?: number };

type AuthContextType = AuthState & {
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
  notify: (toast: Omit<Toast, "id">) => void;
  toasts: Toast[];
  removeToast: (id: number) => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);
  const [email, setEmail] = useState<string | null>(null);
  const [toasts, setToasts] = useState<Toast[]>([]);

  useEffect(() => {
    const t = localStorage.getItem("auth_token");
    const e = localStorage.getItem("auth_email");
    if (t) setToken(t);
    if (e) setEmail(e);
  }, []);

  const isAuthenticated = !!token;

  const notify = useCallback((toast: Omit<Toast, "id">) => {
    const id = Date.now() + Math.floor(Math.random() * 1000);
    setToasts((prev) => [...prev, { id, ...toast }]);
    if (toast.timeout !== 0) {
      const timeoutMs = toast.timeout ?? 4000;
      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
      }, timeoutMs);
    }
  }, []);

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await fetch(`${API_BASE}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      const msg = await res.text();
      notify({ type: "danger", message: msg || "Error de autenticación" });
      throw new Error(msg || "login failed");
    }
    const data = (await res.json()) as { token: string };
    setToken(data.token);
    setEmail(email);
    localStorage.setItem("auth_token", data.token);
    localStorage.setItem("auth_email", email);
    notify({ type: "success", message: "Has iniciado sesión" });
  }, [notify]);

  const register = useCallback(async (email: string, password: string) => {
    const res = await fetch(`${API_BASE}/api/users`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password, role: "client" }),
    });
    if (!res.ok) {
      const msg = await res.text();
      notify({ type: "danger", message: msg || "No se pudo registrar" });
      throw new Error(msg || "register failed");
    }
    notify({ type: "success", message: "Usuario registrado. Ahora puedes iniciar sesión." });
  }, [notify]);

  const logout = useCallback(() => {
    setToken(null);
    setEmail(null);
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_email");
    notify({ type: "info", message: "Sesión cerrada" });
  }, [notify]);

  const value = useMemo<AuthContextType>(
    () => ({ token, email, isAuthenticated, login, register, logout, notify, toasts, removeToast }),
    [token, email, isAuthenticated, login, register, logout, notify, toasts, removeToast]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
};
