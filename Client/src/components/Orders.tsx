import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getOrders } from "../api/orders";
import { DashboardShell } from "./dashboard/DashboardShell";
import { Order } from "./dashboard/types";
import { OrderDetailsModal } from "./orders/OrderDetailsModal";
import { OrderFilters } from "./orders/OrderFilters";
import { OrderSummaryCards } from "./orders/OrderSummaryCards";
import { calculateOrderTotal, formatOrderDate, formatStatusLabel, getStatusStyles } from "./orders/orderUtils";

export const Orders = () => {
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

    const totalRevenue = useMemo(() => orders.reduce((sum, order) => sum + calculateOrderTotal(order), 0), [orders]);

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
            setErrorMessage((error as Error).message || "Failed to load orders.");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (!isAdmin) {
            navigate("/history", { replace: true });
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
            activeTab="orders"
            title="All Orders"
            subtitle="Manage customer orders across the system"
            onLogout={handleLogout}
            contentVariant="grid"
            showSidebar={true}
        >
            <div className="lg:col-span-3 space-y-6">
                <OrderSummaryCards totalOrders={orders.length} totalRevenue={totalRevenue} />

                <OrderFilters
                    search={search}
                    status={status}
                    dateFrom={dateFrom}
                    dateTo={dateTo}
                    searchPlaceholder="Search order ID or username"
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

                <div className="bg-white rounded-2xl shadow-md overflow-hidden">
                    <div className="border-b border-gray-100 px-5 py-4 flex flex-wrap items-center justify-between gap-3">
                        <div>
                            <h3 className="text-xl font-semibold text-gray-900">Order Management</h3>
                            <p className="text-sm text-gray-500">Click View Details to inspect any order</p>
                        </div>
                        <p className="text-sm text-gray-500">{orders.length} order(s)</p>
                    </div>

                    {loading ? (
                        <div className="p-5 space-y-3">
                            <div className="h-16 animate-pulse rounded-xl bg-gray-100" />
                            <div className="h-16 animate-pulse rounded-xl bg-gray-100" />
                            <div className="h-16 animate-pulse rounded-xl bg-gray-100" />
                        </div>
                    ) : orders.length === 0 ? (
                        <div className="px-6 py-12 text-center">
                            <p className="text-lg font-semibold text-gray-700">No orders found</p>
                            <p className="mt-2 text-sm text-gray-500">Try adjusting search, status, or date filters.</p>
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full min-w-[820px]">
                                <thead>
                                    <tr className="bg-gradient-to-r from-indigo-600 to-cyan-500 text-white">
                                        <th className="px-5 py-3 text-left text-xs font-bold uppercase tracking-wide">Order ID</th>
                                        <th className="px-5 py-3 text-left text-xs font-bold uppercase tracking-wide">User Name</th>
                                        <th className="px-5 py-3 text-left text-xs font-bold uppercase tracking-wide">Date</th>
                                        <th className="px-5 py-3 text-left text-xs font-bold uppercase tracking-wide">Total Amount</th>
                                        <th className="px-5 py-3 text-left text-xs font-bold uppercase tracking-wide">Status</th>
                                        <th className="px-5 py-3 text-left text-xs font-bold uppercase tracking-wide">Actions</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-100">
                                    {orders.map((order) => (
                                        <tr key={order.id} className="hover:bg-gray-50 transition-colors">
                                            <td className="px-5 py-4 text-sm font-semibold text-gray-900">#{order.id}</td>
                                            <td className="px-5 py-4 text-sm text-gray-700">{order.user_name || "Unknown user"}</td>
                                            <td className="px-5 py-4 text-sm text-gray-700">{formatOrderDate(order.created_at)}</td>
                                            <td className="px-5 py-4 text-sm font-semibold text-gray-900">${calculateOrderTotal(order).toFixed(2)}</td>
                                            <td className="px-5 py-4 text-sm">
                                                <span className={["inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold", getStatusStyles(order.status)].join(" ")}>{formatStatusLabel(order.status)}</span>
                                            </td>
                                            <td className="px-5 py-4 text-sm">
                                                <button
                                                    type="button"
                                                    onClick={() => setSelectedOrder(order)}
                                                    className="rounded-lg bg-indigo-100 px-3 py-2 text-xs font-semibold text-indigo-700 hover:bg-indigo-200 transition-all duration-200"
                                                >
                                                    View Details
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            </div>

            <OrderDetailsModal order={selectedOrder} onClose={() => setSelectedOrder(null)} />
        </DashboardShell>
    );
};
