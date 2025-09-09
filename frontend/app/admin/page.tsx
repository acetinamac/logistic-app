"use client";
import React from 'react';
const API = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export default function AdminPage() {
  const [token, setToken] = React.useState<string>("");
  const [orders, setOrders] = React.useState<any[]>([]);

  const login = async () => {
    const res = await fetch(`${API}/api/login`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ user_id: 999, role: 'admin' }) });
    const data = await res.json();
    setToken(data.token);
    loadAll(data.token);
  };

  const loadAll = async (tok = token) => {
    const res = await fetch(`${API}/api/admin/orders`, { headers: { Authorization: `Bearer ${tok}` } });
    const data = await res.json();
    setOrders(data);
  };

  const updateStatus = async (id: number, status: string) => {
    const res = await fetch(`${API}/api/admin/orders/${id}/status`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify({ status }) });
    if (res.status === 204) {
      await loadAll();
    } else {
      alert('Error: ' + (await res.text()));
    }
  };

  return (
    <div>
      <h2>Admin</h2>
      <button onClick={login}>Login (admin)</button>
      {token && <>
        <button onClick={() => loadAll()}>Refrescar</button>
        <ul>
          {orders.map((o) => (
            <li key={o.id}>
              #{o.id} - status: {o.status} - peso: {o.weight_kg}kg - size: {o.size}
              <select onChange={(e) => updateStatus(o.id, e.target.value)} defaultValue="">
                <option value="" disabled>Actualizar…</option>
                <option value="creado">creado</option>
                <option value="recolectado">recolectado</option>
                <option value="en_estacion">en_estacion</option>
                <option value="en_ruta">en_ruta</option>
                <option value="entregado">entregado</option>
                <option value="cancelado">cancelado</option>
              </select>
            </li>
          ))}
        </ul>
      </>}
    </div>
  );
}
