"use client";
import React, {useState} from "react";
import {useAuth} from "./AuthContext";

const LoginForm: React.FC<{ onSwitchToRegister?: () => void; onSuccess?: () => void }> = ({
                                                                                              onSwitchToRegister,
                                                                                              onSuccess
                                                                                          }) => {
    const {login} = useAuth();
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
        <div className="login-container">
            <div className="login-card row g-0">
                <div className="col-md-6 left-panel">
                    <img src="/images/login-back.png" alt="Login background"/>
                </div>

                <div className="col-md-6 right-panel">
                    <form onSubmit={handleSubmit}>
                        <h2 className="form-title mb-4">Bienvenido</h2>
                        <div className="mb-3">
                            <label className="form-label">Email</label>
                            <input type="email" className="form-control" value={email}
                                   onChange={(e) => setEmail(e.target.value)} required/>
                        </div>
                        <div className="mb-5">
                            <label className="form-label">Clave de usuario</label>
                            <input type="password" className="form-control" value={password}
                                   onChange={(e) => setPassword(e.target.value)} required/>
                        </div>
                        <button className="btn btn-primary w-100" type="submit" disabled={loading}>
                            {loading ? "Ingresando..." : "Ingresar"}
                        </button>
                    </form>
                    <div className="mt-3 text-center">
                        <span className="text-muted">¿No tienes cuenta? </span>
                        <button type="button" className="btn btn-link p-0 align-baseline"
                                onClick={onSwitchToRegister}>Regístrate
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default LoginForm;
