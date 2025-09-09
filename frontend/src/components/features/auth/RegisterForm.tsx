"use client";
import React, { useState } from "react";
import { useAuth } from "./AuthContext";

const RegisterForm: React.FC<{ onSwitchToLogin?: () => void }>=({ onSwitchToLogin })=>{
  const { register } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirm) {
      alert("Las contraseñas no coinciden");
      return;
    }
    setLoading(true);
    try {
      await register(email, password);
      // Volver al login
      onSwitchToLogin?.();
    } catch {
      // notifications handled by context
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card shadow-sm" style={{ maxWidth: 420, width: "100%" }}>
      <div className="card-body">
        <h5 className="card-title mb-3">Registro</h5>
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label className="form-label">Email</label>
            <input type="email" className="form-control" value={email} onChange={(e)=>setEmail(e.target.value)} required />
          </div>
          <div className="mb-3">
            <label className="form-label">Password</label>
            <input type="password" className="form-control" value={password} onChange={(e)=>setPassword(e.target.value)} required />
          </div>
          <div className="mb-3">
            <label className="form-label">Confirmar Password</label>
            <input type="password" className="form-control" value={confirm} onChange={(e)=>setConfirm(e.target.value)} required />
          </div>
          <button className="btn btn-primary w-100" type="submit" disabled={loading}>
            {loading? "Registrando..." : "Registrar"}
          </button>
        </form>
        <div className="mt-3 text-center">
          <button type="button" className="btn btn-link p-0 align-baseline" onClick={onSwitchToLogin}>¿Ya tienes cuenta? Inicia sesión</button>
        </div>
      </div>
    </div>
  );
};

export default RegisterForm;
