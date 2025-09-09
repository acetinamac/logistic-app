"use client";
import React from 'react';

const API = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export default function ClientPage() {
  const [token, setToken] = React.useState<string>("");
  const [orders, setOrders] = React.useState<any[]>([]);
  const [form, setForm] = React.useState<any>({
    origin_coord: { lat: 19.4326, lng: -99.1332 },
    destination_coord: { lat: 20.6597, lng: -103.3496 },
    origin_address: { country: "MX", state: "CDMX", city: "CDMX", zipcode: "01000", street: "Av. Reforma", ext_num: "1", int_num: "" },
    destination_address: { country: "MX", state: "JAL", city: "GDL", zipcode: "44100", street: "Av. Juarez", ext_num: "10", int_num: "2" },
    items_count: 1,
    weight_kg: 3
  });

  const login = async () => {
    const res = await fetch(`${API}/api/login`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ user_id: 1, role: 'client' }) });
    const data = await res.json();
    setToken(data.token);
    loadOrders(data.token);
  };

  const loadOrders = async (tok = token) => {
    const res = await fetch(`${API}/api/orders`, { headers: { Authorization: `Bearer ${tok}` } });
    const data = await res.json();
    setOrders(data);
  };

  const createOrder = async () => {
    const res = await fetch(`${API}/api/orders`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(form) });
    if (res.ok) {
      await loadOrders();
      alert('Orden creada');
    } else {
      alert('Error al crear orden: ' + (await res.text()));
    }
  };

  return (
    <div>
      <h2>Cliente</h2>
      <button onClick={login}>Login (cliente)</button>
      {token && <>
        <h3>Crear Orden</h3>
        <label>Peso (kg): <input type="number" value={form.weight_kg} onChange={e => setForm({ ...form, weight_kg: parseFloat(e.target.value) })} /></label>
        <button onClick={createOrder}>Crear</button>
        <h3>Mis Órdenes</h3>
        <pre>{JSON.stringify(orders, null, 2)}</pre>
      </>}
    </div>
  );
}
