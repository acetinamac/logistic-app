"use client";
import React, { useState } from "react";
import { useAuth } from "./AuthContext";

const LoginForm: React.FC<{ onSwitchToRegister?: () => void; onSuccess?: () => void }>=({ onSwitchToRegister, onSuccess })=>{
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login(email, password);
      onSuccess?.();
    } catch {
      // notification handled by context
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card shadow-sm" style={{ maxWidth: 420, width: "100%" }}>
      <div className="card-body">
        <h5 className="card-title mb-3">Iniciar sesión</h5>
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label className="form-label">Email</label>
            <input type="email" className="form-control" value={email} onChange={(e)=>setEmail(e.target.value)} required />
          </div>
          <div className="mb-3">
            <label className="form-label">Password</label>
            <input type="password" className="form-control" value={password} onChange={(e)=>setPassword(e.target.value)} required />
          </div>
          <button className="btn btn-primary w-100" type="submit" disabled={loading}>
            {loading? "Ingresando..." : "Ingresar"}
          </button>
        </form>
        <div className="mt-3 text-center">
          <span className="text-muted">¿No tienes cuenta? </span>
          <button type="button" className="btn btn-link p-0 align-baseline" onClick={onSwitchToRegister}>Regístrate</button>
        </div>
      </div>
    </div>
  );
};

export default LoginForm;
