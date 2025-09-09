"use client";
import React from "react";

const Sidebar: React.FC<{ current: string; onSelect: (k: string) => void }>=({ current, onSelect })=>{
  return (
    <aside className="bg-light border-end" style={{ minHeight: "calc(100vh - 56px)" }}>
      <div className="list-group list-group-flush">
        <button className={`list-group-item list-group-item-action ${current==='crear'?'active':''}`} onClick={()=>onSelect('crear')}>Crear órdenes</button>
        <button className={`list-group-item list-group-item-action ${current==='consultar'?'active':''}`} onClick={()=>onSelect('consultar')}>Consultar mis órdenes</button>
      </div>
    </aside>
  );
};

export default Sidebar;
