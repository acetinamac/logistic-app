"use client";
import React, { useState } from "react";
import { useAuth } from "../features/auth/AuthContext";

const Header: React.FC = () => {
  const { isAuthenticated, email, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <header className="navbar navbar-dark px-3" style={{ position: "sticky", top: 0, zIndex: 1030 }}>
      <div className="container-fluid">
        <span className="navbar-brand mb-0 h1">Logistic App</span>
        <div className="d-flex align-items-center gap-2 position-relative">
          {isAuthenticated && (
            <div className="dropdown">
              <button className="btn btn-outline-light btn-sm dropdown-toggle" type="button" onClick={()=>setMenuOpen((v)=>!v)}>
                {email}
              </button>
              {menuOpen && (
                <ul className="dropdown-menu dropdown-menu-end show" style={{ position: "absolute", inset: "auto 0 0 auto", transform: "translate3d(0, 38px, 0)" }}>
                  <li>
                    <button className="dropdown-item" onClick={()=>{ setMenuOpen(false); logout(); }}>Cerrar sesión</button>
                  </li>
                </ul>
              )}
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;
