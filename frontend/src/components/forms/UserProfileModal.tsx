"use client";
import React from "react";
import { API_BASE } from "../../lib/constants";
import { useAuth } from "../features/auth/AuthContext";
import { DataGrid, type Column } from "react-data-grid";
import "react-data-grid/lib/styles.css";

export type UserProfile = {
  id: number;
  full_name: string;
  email: string;
  phone: string;
  role: "client" | "admin";
  is_active: boolean;
  created_at: string;
};

export type Address = {
  id: number;
  customer_id?: number;
  street: string;
  exterior_number: string;
  interior_number?: string;
  neighborhood: string;
  postal_code: string;
  state: string;
  city: string;
  country: string;
  is_active: boolean;
  latitude?: number;
  longitude?: number;
};

export type UserProfileModalProps = {
  open: boolean;
  onClose: () => void;
};

const modalStyle: React.CSSProperties = {
  position: "fixed",
  inset: 0,
  background: "rgba(0,0,0,.45)",
  zIndex: 1050,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  padding: 16
};

const dialogStyle: React.CSSProperties = {
  width: "min(100%, 1000px)",
  maxHeight: "92vh",
  overflow: "auto",
  background: "#fff",
  borderRadius: 8,
  boxShadow: "0 10px 30px rgba(0,0,0,.2)"
};

const mapBoxStyle: React.CSSProperties = {
    width: "100%",
    height: 320,
    borderRadius: 8,
    overflow: "hidden",
    background: "#f6f6f6",
    border: "1px solid #e5e5e5"
};

const UserProfileModal: React.FC<UserProfileModalProps> = ({ open, onClose }) => {
  const { token, userId, notify } = useAuth();

  const [loading, setLoading] = React.useState(false);
  const [saving, setSaving] = React.useState(false);
  const [user, setUser] = React.useState<UserProfile | null>(null);
  const [addresses, setAddresses] = React.useState<Address[]>([]);

  const [addrForm, setAddrForm] = React.useState({
    exterior_number: "",
    interior_number: "",
    neighborhood: "",
    postal_code: "",
    state: "",
    street: "",
    city: "",
    country: "México",
    latitude: 0,
    longitude: 0,
  });

  const onAddrField = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setAddrForm((f) => ({ ...f, [name]: value }));
  };

  const clearAddr = () => {
    setAddrForm({
      exterior_number: "",
      interior_number: "",
      neighborhood: "",
      postal_code: "",
      state: "",
      street: "",
      city: "",
      country: "México",
      latitude: 0,
      longitude: 0,
    });
  };

  const fetchUser = React.useCallback(async () => {
    if (!token || !userId) return;
    setLoading(true);
    try {
      const [userRes, addrRes] = await Promise.all([
        fetch(`${API_BASE}/api/users/${userId}`, { headers: { Authorization: `Bearer ${token}` } }),
        fetch(`${API_BASE}/api/addresses`, { headers: { Authorization: `Bearer ${token}` } })
      ]);

      if (!userRes.ok) throw new Error(await userRes.text());
      if (!addrRes.ok) throw new Error(await addrRes.text());

      const u = (await userRes.json()) as UserProfile;
      const addrs = (await addrRes.json()) as Address[];
      setUser(u);

      const filtered = addrs.filter(a => (a as any).customer_id ? (a as any).customer_id === userId : true);
      setAddresses(filtered);
    } catch (e: any) {
      notify({ type: "danger", message: e?.message || "Error cargando datos de perfil" });
    } finally {
      setLoading(false);
    }
  }, [token, userId, notify]);

  React.useEffect(() => {
    if (open) {
      fetchUser();
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition((pos) => {
          setAddrForm(f => ({ ...f, latitude: pos.coords.latitude, longitude: pos.coords.longitude }));
        });
      }
    }
  }, [open, fetchUser]);

  // Google Maps embedding: load script when open
  const mapRef = React.useRef<HTMLDivElement | null>(null);
  const mapObjRef = React.useRef<any>(null);
  const markerRef = React.useRef<any>(null);

  const goToCurrentLocation = React.useCallback(() => {
    if (!navigator.geolocation) {
      notify({ type: "warning", message: "Geolocalización no disponible" });
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const lat = pos.coords.latitude;
        const lng = pos.coords.longitude;
        setAddrForm((f) => ({ ...f, latitude: lat, longitude: lng }));
        try {
          if (markerRef.current) markerRef.current.setPosition({ lat, lng });
          if (mapObjRef.current) mapObjRef.current.setCenter({ lat, lng });
        } catch {}
      },
      () => notify({ type: "warning", message: "No se pudo obtener la ubicación actual" }),
      { enableHighAccuracy: true, timeout: 10000 }
    );
  }, [notify]);

  const initMap = React.useCallback(() => {
    if (!mapRef.current || (window as any).google?.maps == null) return;
    const { latitude, longitude } = addrForm.latitude && addrForm.longitude ? addrForm : { latitude: 19.432608, longitude: -99.133209 };
    const center = { lat: Number(latitude) || 19.432608, lng: Number(longitude) || -99.133209 };
    const map = new (window as any).google.maps.Map(mapRef.current, { zoom: 14, center, mapTypeControl: false, streetViewControl: false });
    const marker = new (window as any).google.maps.Marker({ position: center, map, draggable: true });
    map.addListener("click", (e: any) => {
      const lat = e.latLng.lat();
      const lng = e.latLng.lng();
      marker.setPosition({ lat, lng });
      setAddrForm((f) => ({ ...f, latitude: lat, longitude: lng }));
    });
    marker.addListener("dragend", () => {
      const p = marker.getPosition();
      const lat = p.lat();
      const lng = p.lng();
      setAddrForm((f) => ({ ...f, latitude: lat, longitude: lng }));
    });
    mapObjRef.current = map;
    markerRef.current = marker;
  }, [addrForm.latitude, addrForm.longitude]);

  React.useEffect(() => {
    if (!open) return;
    if ((window as any).google?.maps) {
      initMap();
      return;
    }

    const scriptId = "google-maps-sdk";
    if (document.getElementById(scriptId)) return;

    const script = document.createElement("script");
    script.id = scriptId;
    // NOTE: expects an environment-configured API key available via ?key=... in .env or index.html injection. If not set, map won't load but form still works.
    const key = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || "";
    const url = `https://maps.googleapis.com/maps/api/js?libraries=places${key ? `&key=${key}` : ""}`;
    script.src = url;
    script.async = true;
    script.onload = () => initMap();
    script.onerror = () => console.warn("Google Maps SDK no pudo cargarse");
    document.body.appendChild(script);
  }, [open, initMap]);

  React.useEffect(() => {
    // update marker when lat/lng typed manually
    if (markerRef.current && mapObjRef.current) {
      const pos = { lat: Number(addrForm.latitude) || 0, lng: Number(addrForm.longitude) || 0 };
      if (!isNaN(pos.lat) && !isNaN(pos.lng) && (pos.lat !== 0 || pos.lng !== 0)) {
        markerRef.current.setPosition(pos);
        mapObjRef.current.setCenter(pos);
      }
    }
  }, [addrForm.latitude, addrForm.longitude]);

  const canSave = () => {
    const f = addrForm;
    return !!(f.street && f.exterior_number && f.neighborhood && f.postal_code && f.city && f.state && f.country);
  };

  const saveAddress = async () => {
    if (!token || !userId) {
      notify({ type: "danger", message: "No autenticado" });
      return;
    }
    if (!canSave()) {
      notify({ type: "warning", message: "Completa los campos obligatorios" });
      return;
    }
    setSaving(true);
    try {
      const body = {
        exterior_number: addrForm.exterior_number,
        interior_number: addrForm.interior_number || "",
        neighborhood: addrForm.neighborhood,
        postal_code: addrForm.postal_code,
        state: addrForm.state,
        street: addrForm.street,
        city: addrForm.city,
        country: addrForm.country,
        coordinates: {
          latitude: Number(addrForm.latitude) || 0,
          longitude: Number(addrForm.longitude) || 0,
        },
        customer_id: userId,
      } as any;
      const res = await fetch(`${API_BASE}/api/addresses`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify(body),
      });
      if (!res.ok) throw new Error(await res.text());
      notify({ type: "success", message: "Dirección guardada" });
      clearAddr();
      goToCurrentLocation();

      // refresh
      const listRes = await fetch(`${API_BASE}/api/addresses`, { headers: { Authorization: `Bearer ${token}` } });
      if (listRes.ok) {
        const addrs = (await listRes.json()) as Address[];
        const filtered = addrs.filter(a => (a as any).customer_id ? (a as any).customer_id === userId : true);
        setAddresses(filtered);
      }
    } catch (e: any) {
      notify({ type: "danger", message: e?.message || "No se pudo guardar la dirección" });
    } finally {
      setSaving(false);
    }
  };

  if (!open) return null;

  const addrColumns: readonly Column<any>[] = [
    { key: "id", name: "Id", width: 80 },
    { key: "direccion", name: "Dirección", resizable: true },
    {
      key: "is_active",
      name: "Estado",
      width: 110,
      renderCell({ row }) {
        const active = !!row.is_active;
        return <span className={active ? "text-success" : "text-danger"}>{active ? "Activo" : "Inactivo"}</span>;
      }
    },
  ];

  const addrRows = addresses.map(a => ({
    ...a,
    direccion: `${a.street ?? ""} ${a.exterior_number ?? ""} ${a.neighborhood ?? ""} ${a.city ?? ""} ${a.country ?? ""} ${a.postal_code ?? ""}`.replace(/\s+/g, " ").trim()
  }));

  return (
    <div style={modalStyle} role="dialog" aria-modal>
      <div style={dialogStyle} className="p-3">
        <div className="d-flex align-items-center justify-content-between mb-3">
          <h5 className="m-0">Perfil de usuario</h5>
          <button className="btn btn-sm btn-outline-secondary" onClick={onClose}>Cerrar</button>
        </div>

        {/* User info */}
        <div className="card mb-3">
          <div className="card-body">
            <div className="row g-3">
              <div className="col-12 col-md-4">
                <label className="form-label">Nombre completo</label>
                <input className="form-control-plaintext bg-transparent" style={{border:0}} value={user?.full_name ?? ""} readOnly tabIndex={-1} />
              </div>
              <div className="col-12 col-md-4">
                <label className="form-label">Email</label>
                <input className="form-control-plaintext bg-transparent" style={{border:0}} value={user?.email ?? ""} readOnly tabIndex={-1} />
              </div>
              <div className="col-12 col-md-4">
                <label className="form-label">Teléfono</label>
                <input className="form-control-plaintext bg-transparent" style={{border:0}} value={user?.phone ?? ""} readOnly tabIndex={-1} />
              </div>
              <div className="col-12 col-md-4">
                <label className="form-label">Rol</label>
                <input className="form-control-plaintext bg-transparent" style={{border:0}} value={user?.role ?? ""} readOnly tabIndex={-1} />
              </div>
              <div className="col-12 col-md-4">
                <label className="form-label">Activo</label>
                <input className="form-control-plaintext bg-transparent" style={{border:0}} value={user?.is_active ? "Sí" : "No"} readOnly tabIndex={-1} />
              </div>
              <div className="col-12 col-md-4">
                <label className="form-label">Creado</label>
                <input className="form-control-plaintext bg-transparent" style={{border:0}} value={user?.created_at ?? ""} readOnly tabIndex={-1} />
              </div>
            </div>
          </div>
        </div>

        {/* Address form */}
        <div className="card mb-3">
          <div className="card-header bg-light">Agregar dirección</div>
          <div className="card-body">
            <div className="row g-3">
              <div className="col-12">
                <div className="mb-2">Selecciona ubicación en el mapa (clic o arrastra el marcador)</div>
                <div ref={mapRef} style={mapBoxStyle} />
                <div className="row g-2 mt-2 align-items-end">
                  <div className="col-6 col-md-3">
                    <label className="form-label">Latitud</label>
                    <input type="number" step="any" className="form-control" name="latitude" value={addrForm.latitude} onChange={(e:any)=>setAddrForm(f=>({...f, latitude: parseFloat(e.target.value)||0}))} />
                  </div>
                  <div className="col-6 col-md-3">
                    <label className="form-label">Longitud</label>
                    <input type="number" step="any" className="form-control" name="longitude" value={addrForm.longitude} onChange={(e:any)=>setAddrForm(f=>({...f, longitude: parseFloat(e.target.value)||0}))} />
                  </div>
                  <div className="col-12 col-md-3">
                    <button type="button" className="btn btn-outline-primary mt-3 mt-md-0" onClick={goToCurrentLocation}>
                      Mandar a mi ubicación
                    </button>
                  </div>
                </div>
              </div>

              <div className="col-12 col-md-4">
                <label className="form-label">Calle</label>
                <input className="form-control" name="street" value={addrForm.street} onChange={onAddrField} />
              </div>
              <div className="col-6 col-md-2">
                <label className="form-label">No. exterior</label>
                <input className="form-control" name="exterior_number" value={addrForm.exterior_number} onChange={onAddrField} />
              </div>
              <div className="col-6 col-md-2">
                <label className="form-label">No. interior</label>
                <input className="form-control" name="interior_number" value={addrForm.interior_number} onChange={onAddrField} />
              </div>
              <div className="col-12 col-md-4">
                <label className="form-label">Colonia</label>
                <input className="form-control" name="neighborhood" value={addrForm.neighborhood} onChange={onAddrField} />
              </div>

              <div className="col-6 col-md-3">
                <label className="form-label">Código postal</label>
                <input className="form-control" name="postal_code" value={addrForm.postal_code} onChange={onAddrField} />
              </div>
              <div className="col-6 col-md-3">
                <label className="form-label">Ciudad</label>
                <input className="form-control" name="city" value={addrForm.city} onChange={onAddrField} />
              </div>
              <div className="col-6 col-md-3">
                <label className="form-label">Estado</label>
                <input className="form-control" name="state" value={addrForm.state} onChange={onAddrField} />
              </div>
              <div className="col-6 col-md-3">
                <label className="form-label">País</label>
                <input className="form-control" name="country" value={addrForm.country} onChange={onAddrField} />
              </div>

              <div className="col-12 d-flex gap-2">
                <button className="btn btn-sm btn-secondary" onClick={saveAddress} disabled={saving || !canSave()}>
                  {saving ? "Guardando..." : "Guardar dirección"}
                </button>
                <button className="btn btn-sm btn-outline-secondary" onClick={() => { clearAddr(); goToCurrentLocation(); }} disabled={saving}>Limpiar</button>
              </div>
            </div>
          </div>
        </div>

        {/* Address list */}
        <div className="card">
          <div className="card-header bg-light">Mis direcciones</div>
          <div className="card-body p-0" style={{height: 260}}>
            <DataGrid className="rdg-light h-100" columns={addrColumns} rows={addrRows} rowKeyGetter={(r)=>String(r.id)} />
          </div>
        </div>

        <div className="d-flex justify-content-end mt-3">
          <button className="btn btn-sm btn-outline-secondary" onClick={onClose}>Cerrar</button>
        </div>
      </div>
    </div>
  );
};

export default UserProfileModal;
