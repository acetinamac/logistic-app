"use client";
import React from "react";
import { useAuth } from "../components/features/auth/AuthContext";
import Sidebar from "../components/layout/Sidebar";
import LoginForm from "../components/features/auth/LoginForm";
import RegisterForm from "../components/features/auth/RegisterForm";

export default function Home() {
  const { isAuthenticated } = useAuth();
  const [view, setView] = React.useState<'crear'|'consultar'>('crear');
  const [showRegister, setShowRegister] = React.useState(false);

  if (!isAuthenticated) {
    return (
      <div className="container d-flex justify-content-center align-items-start" style={{ minHeight: "calc(100vh - 56px)" }}>
        <div className="w-100 d-flex justify-content-center" style={{ marginTop: 40 }}>
          {!showRegister ? (
            <LoginForm onSwitchToRegister={()=>setShowRegister(true)} />
          ) : (
            <RegisterForm onSwitchToLogin={()=>setShowRegister(false)} />
          )}
        </div>
      </div>
    );
  }

  return (
    <div className="container-fluid">
      <div className="row">
        <div className="col-12 col-md-3 col-lg-2 p-0">
          <Sidebar current={view} onSelect={(v)=>setView(v as any)} />
        </div>
        <div className="col-12 col-md-9 col-lg-10 p-4">
          {view === 'crear' && (
            <>
              <h2>Crear órdenes</h2>
              <p className="text-muted">(Próximamente: formulario para crear una orden)</p>
            </>
          )}
          {view === 'consultar' && (
            <>
              <h2>Consultar mis órdenes</h2>
              <p className="text-muted">(Próximamente: listado de órdenes del usuario)</p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
