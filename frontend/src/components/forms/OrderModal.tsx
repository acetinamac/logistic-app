"use client";
import React from "react";
import {API_BASE} from "../../lib/constants";
import {useAuth} from "../features/auth/AuthContext";

export type OrderStatusOption = { label: string; value: string };
export type PackageType = {
    id: number;
    size_code: string;
    max_weight_kg: number;
    description: string;
    is_active: boolean
};
export type Address = {
    id: number;
    street: string;
    exterior_number: string;
    interior_number?: string;
    neighborhood: string;
    postal_code: string;
    city: string;
    state?: string;
    country?: string;
    is_active: boolean
};

export type OrderDetail = {
    id: number;
    order_number: string;
    created_at: string;
    user_id: number;
    full_name: string;
    origin_address_id: number;
    ao_street: string;
    ao_exterior: string;
    ao_neighborhood: string;
    ao_city: string;
    ao_postal: string;
    destination_address_id: number;
    ad_street: string;
    ad_exterior: string;
    ad_neighborhood: string;
    ad_city: string;
    ad_postal: string;
    quantity: number;
    actual_weight_kg: number;
    package_type_id: number;
    size_code: string;
    observations: string;
    internal_notes: string;
    updated_at: string;
    status: string;
};

export type OrderModalProps = {
    open: boolean;
    mode: "create" | "view";
    orderId?: number;
    onClose: () => void;
    onSaved?: () => void; // called after successful create or status update
};

const modalStyle: React.CSSProperties = {
    position: "fixed",
    inset: 0,
    background: "rgba(0,0,0,.45)",
    zIndex: 1050,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    padding: "16px"
};

const dialogStyle: React.CSSProperties = {
    width: "min(100%, 900px)",
    maxHeight: "90vh",
    overflow: "auto",
    background: "#fff",
    borderRadius: 8,
    boxShadow: "0 10px 30px rgba(0,0,0,.2)"
};

const OrderModal: React.FC<OrderModalProps> = ({open, mode, orderId, onClose, onSaved}) => {
    const {token, userId, role, notify} = useAuth();

    const [addresses, setAddresses] = React.useState<Address[]>([]);
    const [pkgTypes, setPkgTypes] = React.useState<PackageType[]>([]);
    const [statusOptions, setStatusOptions] = React.useState<OrderStatusOption[]>([]);
    const [loading, setLoading] = React.useState(false);
    const [saving, setSaving] = React.useState(false);
    const [detail, setDetail] = React.useState<OrderDetail | null>(null);

    const isAdmin = role === "admin";
    const isView = mode === "view";

    const [form, setForm] = React.useState({
        quantity: 1,
        actual_weight_kg: 0,
        origin_address_id: 0,
        destination_address_id: 0,
        observations: "",
        internal_notes: "",
        status: "created" as string,
    });

    const computePackageTypeId = React.useCallback((weight: number) => {
        if (!pkgTypes.length) return 0;

        const sorted = [...pkgTypes].sort((a, b) => a.max_weight_kg - b.max_weight_kg);
        const lastPkg = sorted[sorted.length - 1];

        // If weight exceeds the maximum weight of the last package type, return -1
        if (weight > lastPkg.max_weight_kg) {
            return -1;
        }

        const hit = sorted.find(p => weight <= p.max_weight_kg);
        return hit?.id ?? 0;
    }, [pkgTypes]);

    const fetchBaseData = React.useCallback(async (customerIdFromOrder?: number
    ) => {
        if (!token) return;

        setLoading(true);
        try {
            let addressesUrl = `${API_BASE}/api/addresses`;
            if (customerIdFromOrder && customerIdFromOrder !== userId) {
                addressesUrl += `?customer_id=${customerIdFromOrder}`;
            }

            const [addrRes, pkgRes, statusRes] = await Promise.all([
                fetch(addressesUrl, {headers: {Authorization: `Bearer ${token}`}}),
                fetch(`${API_BASE}/api/package-types`, {headers: {Authorization: `Bearer ${token}`}}),
                fetch(`${API_BASE}/api/orders/status`, {headers: {Authorization: `Bearer ${token}`}}),
            ]);

            if (!addrRes.ok) throw new Error(await addrRes.text());

            if (!pkgRes.ok) throw new Error(await pkgRes.text());

            if (!statusRes.ok) throw new Error(await statusRes.text());

            const addrs = (await addrRes.json()) as Address[];
            const pkgs = (await pkgRes.json()) as PackageType[];
            const stats = (await statusRes.json()) as OrderStatusOption[];

            // Add special element for packages that exceed the limit
            const specialPackage: PackageType = {
                id: 0,
                size_code: "",
                max_weight_kg: Infinity,
                description: "El peso del paquete excede el límite estándar de 25kg. Para envíos de este tipo, debe contactar a la empresa para generar un convenio especial",
                is_active: true
            };

            setAddresses(addrs);
            setPkgTypes([...pkgs.filter(p => p.is_active), specialPackage]);
            setStatusOptions(stats);
        } catch (e: any) {
            notify({type: "danger", message: e?.message || "Error cargando datos base"});
        } finally {
            setLoading(false);
        }
    }, [token, notify]);

    const fetchDetail = React.useCallback(async () => {
        if (!token || !orderId) return;

        setLoading(true);
        try {
            const res = await fetch(`${API_BASE}/api/orders/${orderId}`, {headers: {Authorization: `Bearer ${token}`}});
            if (!res.ok) throw new Error(await res.text());

            const d = (await res.json()) as OrderDetail;
            setDetail(d);
            if (d.user_id !== userId) {
                await fetchBaseData(d.user_id);
            }
            // preset form with detail (for potential status update)
            setForm(f => ({
                ...f,
                quantity: d.quantity,
                actual_weight_kg: d.actual_weight_kg,
                origin_address_id: d.origin_address_id,
                destination_address_id: d.destination_address_id,
                observations: d.observations || "",
                internal_notes: d.internal_notes || "",
                status: d.status,
            }));
        } catch (e: any) {
            notify({type: "danger", message: e?.message || "Error obteniendo orden"});
        } finally {
            setLoading(false);
        }
    }, [token, orderId, notify]);

    React.useEffect(() => {
        if (!open) return;
        fetchBaseData();

        if (!isView) {
            setDetail(null);
            setForm({
                quantity: 1,
                actual_weight_kg: 0,
                origin_address_id: 0,
                destination_address_id: 0,
                observations: "",
                internal_notes: "",
                status: "created",
            });
        }
    }, [open, fetchBaseData, isView]);

    React.useEffect(() => {
        if (open && isView && orderId) {
            fetchDetail();
        }
    }, [open, isView, orderId, fetchDetail]);

    // handlers
    const onField = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const {name, value} = e.target;
        setForm((f) => ({...f, [name]: name === "actual_weight_kg" ? Number(value) : value}));
    };

    React.useEffect(() => {
        // recompute package type when weight changes
        const id = computePackageTypeId(form.actual_weight_kg);
        // no need to store in form; will compute on submit
    }, [form.actual_weight_kg, computePackageTypeId]);

    const handleCreate = async () => {
        if (!token || !userId) {
            notify({type: "danger", message: "No autenticado"});
            return;
        }
        if (!form.origin_address_id || !form.destination_address_id) {
            notify({type: "warning", message: "Selecciona dirección de origen y destino"});
            return;
        }
        if (form.actual_weight_kg <= 0) {
            notify({type: "warning", message: "El peso debe ser mayor a 0"});
            return;
        }
        setSaving(true);
        try {
            const nowIso = new Date().toISOString();
            const body = {
                quantity: Number(form.quantity),
                actual_weight_kg: form.actual_weight_kg,
                created_at: nowIso,
                created_by: userId,
                customer_id: userId,
                destination_address_id: Number(form.destination_address_id),
                id: 0,
                internal_notes: isAdmin ? form.internal_notes : "",
                observations: form.observations,
                order_number: "",
                origin_address_id: Number(form.origin_address_id),
                package_type_id: computePackageTypeId(form.actual_weight_kg),
                status: "created",
                updated_at: nowIso,
                updated_by: 0
            };
            const res = await fetch(`${API_BASE}/api/orders`, {
                method: "POST",
                headers: {"Content-Type": "application/json", Authorization: `Bearer ${token}`},
                body: JSON.stringify(body),
            });

            if (!res.ok) throw new Error(await res.text());
            notify({type: "success", message: "Orden creada correctamente"});
            onClose();
            onSaved?.();
        } catch (e: any) {
            notify({type: "danger", message: e?.message || "No se pudo crear la orden"});
        } finally {
            setSaving(false);
        }
    };

    const handlePatchStatus = async () => {
        if (!token || !isAdmin || !orderId) return;

        setSaving(true);
        try {
            const res = await fetch(`${API_BASE}/api/orders/${orderId}/status`, {
                method: "PATCH",
                headers: {"Content-Type": "application/json", Authorization: `Bearer ${token}`},
                body: JSON.stringify({internal_notes: isAdmin ? form.internal_notes : "", status: form.status}),
            });
            if (!res.ok) throw new Error(await res.text());
            notify({type: "success", message: "Orden actualizada"});
            onClose();
            onSaved?.();
        } catch (e: any) {
            notify({type: "danger", message: e?.message || "No se pudo actualizar el estatus"});
        } finally {
            setSaving(false);
        }
    };

    if (!open) return null;

    const renderAddressOption = (a: Address) => (
        <option key={a.id} value={a.id}>
            {`${a.street} ${a.exterior_number} ${a.neighborhood} ${a.city} ${a.postal_code}`}
        </option>
    );

    const selectedPkg = (() => {
        const id = computePackageTypeId(form.actual_weight_kg);
        return pkgTypes.find(p => p.id === id);
    })();

    return (
        <div style={modalStyle} role="dialog" aria-modal>
            <div style={dialogStyle} className="p-3">
                <div className="d-flex align-items-center justify-content-between mb-3">
                    <h5 className="m-0">{isView ? `Orden #${detail?.order_number ?? orderId}` : "Crear Orden"}</h5>
                </div>

                {/* Content */}
                <div className="row g-3">
                    <div className="col-12 col-md-6">
                        <label className="form-label">Dirección Origen</label>
                        <select className="form-select" name="origin_address_id" value={Number(form.origin_address_id)}
                                onChange={onField} disabled={isView}>
                            <option value={0}>Selecciona...</option>
                            {addresses.map(renderAddressOption)}
                        </select>
                    </div>
                    <div className="col-12 col-md-6">
                        <label className="form-label">Dirección Destino</label>
                        <select className="form-select" name="destination_address_id"
                                value={Number(form.destination_address_id)}
                                onChange={onField} disabled={isView}>
                            <option value={0}>Selecciona...</option>
                            {addresses.map(renderAddressOption)}
                        </select>
                    </div>

                    <div className="col-12 col-md-3">
                        <label className="form-label">Cantidad (Qty)</label>
                        <input type="number" step="1" min={1} className="form-control" name="quantity"
                               value={form.quantity} onChange={onField} disabled={isView}/>
                    </div>

                    <div className="col-12 col-md-3">
                        <label className="form-label">Peso Actual (kg)</label>
                        <input type="number" step="0.01" min={0} className="form-control" name="actual_weight_kg"
                               value={form.actual_weight_kg} onChange={onField} disabled={isView}/>
                    </div>

                    <div className="col-12 col-md-6">
                        <label className="form-label">Tipo de Paquete</label>
                        <input className="form-control form-control-plaintext bg-transparent" style={{border: 0}}
                               value={selectedPkg ? `${selectedPkg.size_code} - ${selectedPkg.description}` : "—"}
                               readOnly tabIndex={-1}/>
                        <div className="form-text">Se calcula automáticamente según el peso.</div>
                    </div>

                    <div className="col-12">
                        <label className="form-label">Observaciones</label>
                        <textarea className="form-control" name="observations" rows={3} value={form.observations}
                                  onChange={onField} disabled={isView}/>
                    </div>

                    {isAdmin && (
                        <div className="col-12">
                            <label className="form-label">Notas internas</label>
                            <textarea className="form-control" name="internal_notes" rows={3}
                                      value={form.internal_notes}
                                      onChange={onField} disabled={isView && !isAdmin}/>
                        </div>
                    )}

                    <div className="col-12 col-md-6">
                        <label className="form-label">Estatus</label>
                        {!isView && (
                            <input className="form-control form-control-plaintext bg-transparent" style={{border: 0}}
                                   readOnly tabIndex={-1} value="Creado"/>
                        )}
                        {isView && isAdmin && (
                            <select className="form-select" name="status" value={form.status} onChange={onField}>
                                {statusOptions.map(s => (
                                    <option key={s.value} value={s.value}>{s.label}</option>
                                ))}
                            </select>
                        )}
                        {isView && !isAdmin && (
                            <input className="form-control form-control-plaintext bg-transparent" style={{border: 0}}
                                   readOnly tabIndex={-1}
                                   value={statusOptions.find(s => s.value === form.status)?.label || form.status}/>
                        )}
                        {!isAdmin &&
                            <div className="form-text">Solo un admin puede modificarlo.</div>}
                    </div>

                    {isView && detail && (
                        <div className="col-12 col-md-6">
                            <label className="form-label">Fecha de creación</label>
                            <input className="form-control form-control-plaintext bg-transparent" style={{border: 0}}
                                   value={detail.created_at} readOnly tabIndex={-1}/>
                        </div>
                    )}
                </div>

                <div className="d-flex justify-content-end gap-2 mt-4">
                    {!isView && (
                        <button className="btn btn-sm btn-secondary" onClick={handleCreate}
                                disabled={saving || loading}>Guardar</button>
                    )}
                    {isView && isAdmin && (
                        <button className="btn btn-sm btn-secondary" onClick={handlePatchStatus}
                                disabled={saving || loading}>Guardar estatus</button>
                    )}
                    <button className="btn btn-sm btn-outline-secondary" onClick={onClose} disabled={saving}>Cerrar
                    </button>
                </div>
            </div>
        </div>
    );
};

export default OrderModal;
