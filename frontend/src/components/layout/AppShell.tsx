"use client";
import React from "react";
import Header from "./Header";
import { useAuth } from "../features/auth/AuthContext";
import Toaster from "../ui/Toaster";

const AppShell: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return (
    <>
      {isAuthenticated && <Header />}
      <main className="container-fluid">{children}</main>
      <Toaster />
    </>
  );
};

export default AppShell;
