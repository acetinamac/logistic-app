"use client";
import React from "react";
import { DataGrid, type Column } from "react-data-grid";
import "react-data-grid/lib/styles.css";
import { API_BASE } from "../../../lib/constants";
import { useAuth } from "../auth/AuthContext";

// Types that match the backend DTO structure
export type OrderRow = {
  order_number: string;
  created_at: string; // formatted as DD/MM/YYYY from backend
  full_name: string;
  origin_full_address: string;
  destination_full_address: string;
  actual_weight_kg: number;
  size_code: string;
  status: string;
};

const MainContent: React.FC = () => {
  const { token } = useAuth();
  const [rows, setRows] = React.useState<OrderRow[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  // Filters
  const [qOrder, setQOrder] = React.useState("");
  const [qName, setQName] = React.useState("");
  const [qDate, setQDate] = React.useState(""); // DD/MM/YYYY or partial
  const [qStatus, setQStatus] = React.useState("");

  const fetchOrders = React.useCallback(async () => {
    if (!token) {
      setError("No autenticado. Inicia sesión para ver tus órdenes.");
      setRows([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE}/api/orders`, {
        credentials: "include",
        headers: {
          Authorization: `Bearer ${token}`,
        }
      });
      if (!res.ok) {
        const msg = await res.text();
        throw new Error(msg || `HTTP ${res.status}`);
      }
      const data: OrderRow[] = await res.json();
      setRows(data ?? []);
    } catch (e: any) {
      setError(e?.message || "Error cargando ordenes");
    } finally {
      setLoading(false);
    }
  }, [token]);

  React.useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  // Columns for react-data-grid
  const columns = React.useMemo<readonly Column<OrderRow>[]>(
    () => [
      { key: "order_number", name: "Orden", resizable: true, width: 110, frozen: true },
      { key: "created_at", name: "Fecha", resizable: true, width: 110 },
      { key: "full_name", name: "Cliente", resizable: true, minWidth: 140 },
      { key: "origin_full_address", name: "Origen", resizable: true, minWidth: 220 },
      { key: "destination_full_address", name: "Destino", resizable: true, minWidth: 220 },
      { key: "actual_weight_kg", name: "Peso (kg)", resizable: true, width: 100 },
      { key: "size_code", name: "Tamaño", resizable: true, width: 90 },
      { key: "status", name: "Estado", resizable: true, width: 120 },
      {
        key: "actions",
        name: "Acciones",
        width: 120,
        frozen: true,
        renderCell({ row }) {
          return (
            <div className="d-flex gap-2">
              <button className="btn btn-sm btn-outline-primary" title={`Visualizar ${row.order_number}`} onClick={() => { /* placeholder */ }}>
                Visualizar
              </button>
            </div>
          );
        }
      }
    ],
    []
  );

  // Client-side filtering
  const filteredRows = React.useMemo(() => {
    const order = qOrder.trim().toLowerCase();
    const name = qName.trim().toLowerCase();
    const date = qDate.trim().toLowerCase();
    const status = qStatus.trim().toLowerCase();
    return rows.filter(r =>
      (!order || r.order_number.toLowerCase().includes(order)) &&
      (!name || r.full_name.toLowerCase().includes(name)) &&
      (!date || r.created_at.toLowerCase().includes(date)) &&
      (!status || r.status.toLowerCase().includes(status))
    );
  }, [rows, qOrder, qName, qDate, qStatus]);

  return (
    <div className="container-fluid px-3 py-3" style={{ minHeight: "calc(100vh - 56px)" }}>
      <div className="d-flex align-items-center justify-content-between flex-wrap gap-2 mb-3">
        <h2 className="m-0">Ordenes disponibles</h2>
        <div className="d-flex align-items-center gap-2">
            <button className="btn btn-sm btn-secondary">
                {" Crear Orden "}
            </button>
            <button className="btn btn-sm btn-outline-secondary" onClick={fetchOrders} disabled={loading}>
                {loading ? "Cargando..." : " Refrescar "}
            </button>
        </div>
      </div>

      <div className="card shadow-sm mb-3">
        <div className="card-body">
          <div className="row g-2">
            <div className="col-12 col-sm-6 col-md-3">
              <label className="form-label">Buscar por número</label>
              <input className="form-control" placeholder="e.g. 1500" value={qOrder} onChange={e=>setQOrder(e.target.value)} />
            </div>
            <div className="col-12 col-sm-6 col-md-3">
              <label className="form-label">Buscar por cliente</label>
              <input className="form-control" placeholder="e.g. Admin" value={qName} onChange={e=>setQName(e.target.value)} />
            </div>
            <div className="col-12 col-sm-6 col-md-3">
              <label className="form-label">Buscar por fecha</label>
              <input className="form-control" placeholder="DD/MM/AAAA" value={qDate} onChange={e=>setQDate(e.target.value)} />
            </div>
            <div className="col-12 col-sm-6 col-md-3">
              <label className="form-label">Buscar por estado</label>
              <input className="form-control" placeholder="e.g. created" value={qStatus} onChange={e=>setQStatus(e.target.value)} />
            </div>
          </div>
        </div>
      </div>

      {error && (
        <div className="alert alert-danger" role="alert">
          Error al cargar órdenes: {error}
        </div>
      )}

      <div className="card shadow-sm" style={{ height: "calc(100vh - 260px)" }}>
        <div className="card-body p-0 h-100">
          <DataGrid
            className="rdg-light h-100"
            columns={columns}
            rows={filteredRows}
            rowKeyGetter={(r)=>r.order_number}
            defaultColumnOptions={{ sortable: true }}
          />
        </div>
      </div>
    </div>
  );
};

export default MainContent;