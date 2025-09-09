"use client";
import React from "react";
import { useAuth } from "../features/auth/AuthContext";

const Toaster: React.FC = () => {
  const { toasts, removeToast } = useAuth();
  return (
    <div className="toast-container position-fixed top-0 end-0 p-3" style={{ zIndex: 2000 }}>
      {toasts.map((t) => (
        <div key={t.id} className={`toast align-items-center text-bg-${t.type} show mb-2`} role="alert" aria-live="assertive" aria-atomic="true">
          <div className="d-flex">
            <div className="toast-body">{t.message}</div>
            <button type="button" className="btn-close btn-close-white me-2 m-auto" aria-label="Close" onClick={() => removeToast(t.id)}></button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default Toaster;
