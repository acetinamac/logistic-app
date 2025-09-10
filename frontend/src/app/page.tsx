"use client";
import React from "react";
import { useAuth } from "../components/features/auth/AuthContext";
import LoginForm from "../components/features/auth/LoginForm";
import RegisterForm from "../components/features/auth/RegisterForm";
import MainContent from "../components/layout/MainContent";

export default function Home() {
  const { isAuthenticated } = useAuth();
  const [showRegister, setShowRegister] = React.useState(false);

  if (!isAuthenticated) {
    return (
      <div className="container-fluid min-vh-100 d-flex align-items-center justify-content-center">
        <div className="w-100" style={{ maxWidth: 480 }}>
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
    <div className="container-fluid m-0 p-0">
      <div className="row">
          <MainContent />
      </div>
    </div>
  );
}
