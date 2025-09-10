"use client";
import React, {useState, useRef, useEffect} from "react";
import {useAuth} from "../features/auth/AuthContext";

const Header: React.FC = () => {
    const {isAuthenticated, email, logout} = useAuth();
    const [menuOpen, setMenuOpen] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setMenuOpen(false);
            }
        };

        if (menuOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [menuOpen]);

    const handleLogout = () => {
        setMenuOpen(false);
        logout();
    };

    return (
        <header className="navbar navbar-dark px-3 shadow-sm"
                style={{position: "sticky", top: 0, zIndex: 1030}}>
            <div className="container-fluid">
                <span className="navbar-brand mb-0 h1 fw-bold">Logistic App</span>

                <div className="d-flex align-items-center">
                    {isAuthenticated && (
                        <div className="dropdown" ref={dropdownRef}>
                            <div
                                className="d-flex align-items-center text-white cursor-pointer user-select-none"
                                onClick={() => setMenuOpen(!menuOpen)}
                                style={{cursor: 'pointer'}}
                                role="button"
                                aria-expanded={menuOpen}
                                aria-haspopup="true"
                            >
                                <div
                                    className="rounded-circle bg-light d-flex align-items-center justify-content-center me-2 fw-bold text-primary"
                                    style={{width: '32px', height: '32px', fontSize: '14px'}}
                                >
                                    {email?.charAt(0).toUpperCase() || 'U'}
                                </div>

                                {/* Email */}
                                <span className="me-2 d-none d-md-inline text-truncate" style={{maxWidth: '150px'}}>
                  {email}
                </span>

                                {/* Icono dropdown */}
                                <svg
                                    width="12"
                                    height="12"
                                    fill="currentColor"
                                    className={`transition-transform ${menuOpen ? 'rotate-180' : ''}`}
                                    style={{transition: 'transform 0.2s ease'}}
                                    viewBox="0 0 16 16"
                                >
                                    <path fillRule="evenodd"
                                          d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
                                </svg>
                            </div>

                            {/* Menú dropdown */}
                            {menuOpen && (
                                <div
                                    className="dropdown-menu dropdown-menu-end show position-absolute"
                                    style={{
                                        top: '100%',
                                        right: '0',
                                        marginTop: '8px',
                                        minWidth: '200px',
                                        border: '1px solid rgba(0,0,0,0.1)',
                                        boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
                                        borderRadius: '8px'
                                    }}
                                >
                                    <div className="dropdown-header border-bottom pb-2 mb-2">
                                        <div className="d-flex align-items-center">
                                            <div
                                                className="rounded-circle bg-primary d-flex align-items-center justify-content-center me-2 text-white fw-bold"
                                                style={{width: '24px', height: '24px', fontSize: '12px'}}
                                            >
                                                {email?.charAt(0).toUpperCase() || 'U'}
                                            </div>
                                            <div>
                                                <div className="fw-semibold text-dark" style={{fontSize: '14px'}}>
                                                    {email}
                                                </div>
                                                <small className="text-muted">Usuario activo</small>
                                            </div>
                                        </div>
                                    </div>

                                    <li>
                                        <button
                                            className="dropdown-item d-flex align-items-center py-2"
                                            onClick={handleLogout}
                                        >
                                            <svg
                                                width="16"
                                                height="16"
                                                fill="currentColor"
                                                className="me-2 text-muted"
                                                viewBox="0 0 16 16"
                                            >
                                                <path fillRule="evenodd"
                                                      d="M6 12.5a.5.5 0 0 0 .5.5h8a.5.5 0 0 0 .5-.5v-9a.5.5 0 0 0-.5-.5h-8a.5.5 0 0 0-.5.5v2a.5.5 0 0 1-1 0v-2A1.5 1.5 0 0 1 6.5 2h8A1.5 1.5 0 0 1 16 3.5v9a1.5 1.5 0 0 1-1.5 1.5h-8A1.5 1.5 0 0 1 5 12.5v-2a.5.5 0 0 1 1 0v2z"/>
                                                <path fillRule="evenodd"
                                                      d="M.146 8.354a.5.5 0 0 1 0-.708l3-3a.5.5 0 1 1 .708.708L1.707 7.5H10.5a.5.5 0 0 1 0 1H1.707l2.147 2.146a.5.5 0 0 1-.708.708l-3-3z"/>
                                            </svg>
                                            Cerrar sesión
                                        </button>
                                    </li>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </header>
    );
};

export default Header;