import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getOrders } from "../api/orders";
import { DashboardShell } from "./dashboard/DashboardShell";
import { Order } from "./dashboard/types";
import { OrderDetailsModal } from "./orders/OrderDetailsModal";
import { OrderFilters } from "./orders/OrderFilters";
import { calculateOrderTotal, formatOrderDate, formatStatusLabel, getStatusStyles } from "./orders/orderUtils";

export const OrderHistory = () => {
    const navigate = useNavigate();
    const role = (localStorage.getItem("role") || "user").toLowerCase();
    const isAdmin = role === "admin";

    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [errorMessage, setErrorMessage] = useState<string>("");
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [search, setSearch] = useState<string>("");
    const [status, setStatus] = useState<string>("");
    const [dateFrom, setDateFrom] = useState<string>("");
    const [dateTo, setDateTo] = useState<string>("");

    const totalAmount = useMemo(() => orders.reduce((sum, order) => sum + calculateOrderTotal(order), 0), [orders]);

    const loadOrders = async (filters?: { search?: string; status?: string; date_from?: string; date_to?: string }) => {
        setLoading(true);
        setErrorMessage("");
        try {
            const result = await getOrders(filters ?? {});
            setOrders(result);
            if (selectedOrder) {
                const updatedSelection = result.find((order) => order.id === selectedOrder.id) ?? null;
                setSelectedOrder(updatedSelection);
            }
        } catch (error) {
            setErrorMessage((error as Error).message || "Failed to load order history.");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (isAdmin) {
            navigate("/orders", { replace: true });
            return;
        }
        loadOrders();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isAdmin]);

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", { method: "POST", credentials: "include" });
        } finally {
            localStorage.removeItem("token");
            localStorage.removeItem("role");
            localStorage.removeItem("username");
            navigate("/");
        }
    };

    const handleApplyFilters = () => {
        void loadOrders({
            search,
            status,
            date_from: dateFrom,
            date_to: dateTo,
        });
    };

    const handleResetFilters = () => {
        setSearch("");
        setStatus("");
        setDateFrom("");
        setDateTo("");
        void loadOrders();
    };

    return (
        <DashboardShell
            role={role}
            activeTab="history"
            title="My Orders"
            subtitle="Review your personal order history"
            onLogout={handleLogout}
            contentVariant="grid"
            showSidebar={true}
        >
            <div className="lg:col-span-3 space-y-6">
                <OrderFilters
                    search={search}
                    status={status}
                    dateFrom={dateFrom}
                    dateTo={dateTo}
                    searchPlaceholder="Search order ID"
                    onSearchChange={setSearch}
                    onStatusChange={setStatus}
                    onDateFromChange={setDateFrom}
                    onDateToChange={setDateTo}
                    onApply={handleApplyFilters}
                    onReset={handleResetFilters}
                    loading={loading}
                />

                {errorMessage && (
                    <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-medium text-red-700">
                        {errorMessage}
                    </div>
                )}

                {orders.length > 0 && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div className="rounded-2xl bg-white shadow-md p-5">
                            <p className="text-sm text-gray-500">Total Orders</p>
                            <p className="mt-2 text-3xl font-bold text-gray-900">{orders.length}</p>
                        </div>
                        <div className="rounded-2xl bg-white shadow-md p-5">
                            <p className="text-sm text-gray-500">Total Amount</p>
                            <p className="mt-2 text-3xl font-bold text-gray-900">${totalAmount.toFixed(2)}</p>
                        </div>
                        <div className="rounded-2xl bg-white shadow-md p-5">
                            <p className="text-sm text-gray-500">Latest Order</p>
                            <p className="mt-2 text-xl font-bold text-gray-900">#{orders[0]?.id ?? "-"}</p>
                        </div>
                    </div>
                )}

                <div className="space-y-4">
                    {loading ? (
                        <div className="space-y-3">
                            <div className="h-24 animate-pulse rounded-2xl bg-white" />
                            <div className="h-24 animate-pulse rounded-2xl bg-white" />
                            <div className="h-24 animate-pulse rounded-2xl bg-white" />
                        </div>
                    ) : orders.length === 0 ? (
                        <div className="rounded-2xl border border-dashed border-gray-200 bg-white px-6 py-12 text-center shadow-sm">
                            <p className="text-lg font-semibold text-gray-700">No orders yet</p>
                            <p className="mt-2 text-sm text-gray-500">Your placed orders will appear here with item-level details and totals.</p>
                        </div>
                    ) : (
                        orders.map((order) => (
                            <button
                                key={order.id}
                                type="button"
                                onClick={() => setSelectedOrder(order)}
                                className="w-full rounded-2xl border border-gray-100 bg-white p-5 text-left shadow-md transition-all duration-200 hover:-translate-y-0.5 hover:shadow-lg"
                            >
                                <div className="flex flex-wrap items-center justify-between gap-3">
                                    <div>
                                        <p className="text-sm font-semibold uppercase tracking-wide text-gray-500">Order #{order.id}</p>
                                        <p className="text-xs text-gray-500">{formatOrderDate(order.created_at)}</p>
                                    </div>
                                    <div className="flex flex-wrap items-center gap-3">
                                        <span className={["inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold", getStatusStyles(order.status)].join(" ")}>{formatStatusLabel(order.status)}</span>
                                        <span className="text-sm font-semibold text-gray-900">${calculateOrderTotal(order).toFixed(2)}</span>
                                    </div>
                                </div>

                                <div className="mt-4 grid grid-cols-1 gap-2 md:grid-cols-2">
                                    {order.items.slice(0, 3).map((item) => (
                                        <div key={item.id} className="flex items-center justify-between rounded-xl bg-gray-50 px-3 py-2 text-sm text-gray-700">
                                            <span className="truncate pr-3">{item.product_name || `Product #${item.product_id}`}</span>
                                            <span className="font-semibold">Qty: {item.quantity}</span>
                                        </div>
                                    ))}
                                </div>
                                {order.items.length > 3 && (
                                    <p className="mt-3 text-xs text-gray-500">+ {order.items.length - 3} more item(s)</p>
                                )}
                            </button>
                        ))
                    )}
                </div>
            </div>

            <OrderDetailsModal order={selectedOrder} onClose={() => setSelectedOrder(null)} />
        </DashboardShell>
    );
};
