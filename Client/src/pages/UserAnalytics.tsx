import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis, CartesianGrid } from "recharts";
import { getUserAnalyticsBundle } from "../api/userAnalytics";
import { useCart } from "../context/CartContext";
import { DashboardShell } from "../components/dashboard/DashboardShell";
import { UserAnalyticsBundle } from "../components/dashboard/types";

const emptyBundle: UserAnalyticsBundle = {
    summary: {
        totalOrders: 0,
        totalSpent: 0,
        pendingOrders: 0,
    },
    recentOrders: [],
    topProducts: [],
    spendingTrend: [],
    recommendations: [],
};

const currencyFormatter = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 2,
});

const formatStatus = (status: string) => {
    const normalized = status.trim().toLowerCase();
    return normalized.charAt(0).toUpperCase() + normalized.slice(1);
};

const statusClasses = (status: string) => {
    const normalized = status.trim().toLowerCase();
    if (normalized === "completed") {
        return "border-emerald-200 bg-emerald-50 text-emerald-700";
    }
    if (normalized === "cancelled") {
        return "border-rose-200 bg-rose-50 text-rose-700";
    }
    return "border-amber-200 bg-amber-50 text-amber-700";
};

export const UserAnalytics = () => {
    const navigate = useNavigate();
    const role = (localStorage.getItem("role") || "user").toLowerCase();
    const { addItem, refreshCart } = useCart();

    const [bundle, setBundle] = useState<UserAnalyticsBundle>(emptyBundle);
    const [loading, setLoading] = useState<boolean>(false);
    const [buyingProductId, setBuyingProductId] = useState<number | null>(null);

    const hasActivity = useMemo(() => (
        bundle.summary.totalOrders > 0 ||
        bundle.recentOrders.length > 0 ||
        bundle.topProducts.length > 0 ||
        bundle.spendingTrend.length > 0
    ), [bundle]);

    useEffect(() => {
        if (role !== "user") {
            navigate("/dashboard", { replace: true });
            return;
        }

        const loadAnalytics = async () => {
            setLoading(true);
            try {
                const data = await getUserAnalyticsBundle();
                setBundle(data);
            } catch {
                setBundle(emptyBundle);
            } finally {
                setLoading(false);
            }
        };

        void loadAnalytics();
    }, [navigate, role]);

    const handleBuyAgain = async (productId: number) => {
        setBuyingProductId(productId);
        try {
            await addItem(productId, 1);
            await refreshCart();
            navigate("/cart");
        } finally {
            setBuyingProductId(null);
        }
    };

    return (
        <DashboardShell
            role={role}
            activeTab="history"
            title="Your Insights"
            subtitle="Personal spending, order activity, and recommendations tailored to you"
            contentVariant="stack"
        >
            <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
                <div className="rounded-2xl bg-white p-5 shadow-md">
                    <p className="text-sm font-semibold uppercase tracking-wide text-slate-500">Total Orders</p>
                    <p className="mt-3 text-3xl font-bold text-slate-900">{bundle.summary.totalOrders}</p>
                </div>
                <div className="rounded-2xl bg-white p-5 shadow-md">
                    <p className="text-sm font-semibold uppercase tracking-wide text-slate-500">Total Spent</p>
                    <p className="mt-3 text-3xl font-bold text-slate-900">{currencyFormatter.format(bundle.summary.totalSpent)}</p>
                </div>
                <div className="rounded-2xl bg-white p-5 shadow-md">
                    <p className="text-sm font-semibold uppercase tracking-wide text-slate-500">Pending Orders</p>
                    <p className="mt-3 text-3xl font-bold text-slate-900">{bundle.summary.pendingOrders}</p>
                </div>
            </div>

            {!loading && !hasActivity ? (
                <div className="rounded-2xl bg-white px-6 py-12 text-center shadow-md">
                    <p className="text-lg font-semibold text-slate-800">No activity yet</p>
                    <p className="mt-2 text-sm text-slate-500">Your personal analytics will appear here once you place orders.</p>
                </div>
            ) : (
                <>
                    <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                        <section className="rounded-2xl bg-white p-5 shadow-md">
                            <div className="mb-4">
                                <h3 className="text-lg font-semibold text-slate-900">Spending Trend</h3>
                                <p className="text-sm text-slate-500">Track how your spending changes over time.</p>
                            </div>
                            {loading ? (
                                <div className="h-72 animate-pulse rounded-xl bg-slate-100" />
                            ) : bundle.spendingTrend.length === 0 ? (
                                <div className="flex h-72 items-center justify-center rounded-xl border border-dashed border-slate-200 text-sm font-medium text-slate-500">
                                    No activity yet
                                </div>
                            ) : (
                                <div className="h-72">
                                    <ResponsiveContainer width="100%" height="100%">
                                        <LineChart data={bundle.spendingTrend} margin={{ top: 8, right: 18, left: 0, bottom: 6 }}>
                                            <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                                            <XAxis dataKey="date" tick={{ fontSize: 12 }} stroke="#64748b" />
                                            <YAxis tick={{ fontSize: 12 }} stroke="#64748b" tickFormatter={(value) => `$${value}`} />
                                            <Tooltip formatter={(value) => currencyFormatter.format(Number(value ?? 0))} />
                                            <Line type="monotone" dataKey="value" stroke="#2563eb" strokeWidth={3} dot={false} />
                                        </LineChart>
                                    </ResponsiveContainer>
                                </div>
                            )}
                        </section>

                        <section className="rounded-2xl bg-white p-5 shadow-md">
                            <div className="mb-4">
                                <h3 className="text-lg font-semibold text-slate-900">Order Status Snapshot</h3>
                                <p className="text-sm text-slate-500">A quick view of your current order pipeline.</p>
                            </div>
                            <div className="space-y-4">
                                <div className="rounded-2xl border border-slate-100 bg-slate-50 p-4">
                                    <p className="text-sm text-slate-500">Completed or Historical Orders</p>
                                    <p className="mt-2 text-2xl font-bold text-slate-900">
                                        {Math.max(bundle.summary.totalOrders - bundle.summary.pendingOrders, 0)}
                                    </p>
                                </div>
                                <div className="rounded-2xl border border-amber-100 bg-amber-50 p-4">
                                    <p className="text-sm text-amber-700">Pending Orders</p>
                                    <p className="mt-2 text-2xl font-bold text-amber-800">{bundle.summary.pendingOrders}</p>
                                </div>
                                <div className="rounded-2xl border border-blue-100 bg-blue-50 p-4">
                                    <p className="text-sm text-blue-700">Average Spend per Order</p>
                                    <p className="mt-2 text-2xl font-bold text-blue-800">
                                        {bundle.summary.totalOrders > 0
                                            ? currencyFormatter.format(bundle.summary.totalSpent / bundle.summary.totalOrders)
                                            : currencyFormatter.format(0)}
                                    </p>
                                </div>
                            </div>
                        </section>
                    </div>

                    <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                        <section className="rounded-2xl bg-white p-5 shadow-md">
                            <div className="mb-4">
                                <h3 className="text-lg font-semibold text-slate-900">Recent Orders</h3>
                                <p className="text-sm text-slate-500">Your last five purchases at a glance.</p>
                            </div>
                            {loading ? (
                                <div className="space-y-3">
                                    <div className="h-16 animate-pulse rounded-xl bg-slate-100" />
                                    <div className="h-16 animate-pulse rounded-xl bg-slate-100" />
                                    <div className="h-16 animate-pulse rounded-xl bg-slate-100" />
                                </div>
                            ) : bundle.recentOrders.length === 0 ? (
                                <div className="flex h-64 items-center justify-center rounded-xl border border-dashed border-slate-200 text-sm font-medium text-slate-500">
                                    No activity yet
                                </div>
                            ) : (
                                <div className="space-y-3">
                                    {bundle.recentOrders.map((order) => (
                                        <div key={order.orderId} className="rounded-2xl border border-slate-100 bg-slate-50 p-4">
                                            <div className="flex flex-wrap items-start justify-between gap-3">
                                                <div>
                                                    <p className="text-sm font-semibold text-slate-900">Order #{order.orderId}</p>
                                                    <p className="mt-1 text-sm text-slate-600">{order.productNames}</p>
                                                    <p className="mt-1 text-xs text-slate-500">{new Date(order.createdAt).toLocaleDateString()}</p>
                                                </div>
                                                <div className="text-right">
                                                    <span className={["inline-flex rounded-full border px-3 py-1 text-xs font-semibold", statusClasses(order.status)].join(" ")}>
                                                        {formatStatus(order.status)}
                                                    </span>
                                                    <p className="mt-2 text-sm font-semibold text-slate-900">{currencyFormatter.format(order.totalPrice)}</p>
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </section>

                        <section className="rounded-2xl bg-white p-5 shadow-md">
                            <div className="mb-4">
                                <h3 className="text-lg font-semibold text-slate-900">Top Products</h3>
                                <p className="text-sm text-slate-500">The products you come back to most often.</p>
                            </div>
                            {loading ? (
                                <div className="space-y-3">
                                    <div className="h-12 animate-pulse rounded-lg bg-slate-100" />
                                    <div className="h-12 animate-pulse rounded-lg bg-slate-100" />
                                    <div className="h-12 animate-pulse rounded-lg bg-slate-100" />
                                </div>
                            ) : bundle.topProducts.length === 0 ? (
                                <div className="flex h-64 items-center justify-center rounded-xl border border-dashed border-slate-200 text-sm font-medium text-slate-500">
                                    No activity yet
                                </div>
                            ) : (
                                <div className="overflow-x-auto">
                                    <table className="w-full">
                                        <thead>
                                            <tr className="border-b border-slate-200 text-left text-xs uppercase tracking-wide text-slate-500">
                                                <th className="py-3 pr-3">Product</th>
                                                <th className="py-3">Purchase Count</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {bundle.topProducts.map((product) => (
                                                <tr key={product.productId} className="border-b border-slate-100 text-sm text-slate-700">
                                                    <td className="py-3 pr-3 font-medium text-slate-900">{product.name}</td>
                                                    <td className="py-3">{product.purchaseCount}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </section>
                    </div>

                    <section className="rounded-2xl bg-white p-5 shadow-md">
                        <div className="mb-4">
                            <h3 className="text-lg font-semibold text-slate-900">Recommendations</h3>
                            <p className="text-sm text-slate-500">Suggested picks based on your past purchases and shopping patterns.</p>
                        </div>
                        {loading ? (
                            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
                                <div className="h-48 animate-pulse rounded-2xl bg-slate-100" />
                                <div className="h-48 animate-pulse rounded-2xl bg-slate-100" />
                                <div className="h-48 animate-pulse rounded-2xl bg-slate-100" />
                            </div>
                        ) : bundle.recommendations.length === 0 ? (
                            <div className="flex h-48 items-center justify-center rounded-xl border border-dashed border-slate-200 text-sm font-medium text-slate-500">
                                No activity yet
                            </div>
                        ) : (
                            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
                                {bundle.recommendations.map((product) => (
                                    <article key={product.productId} className="rounded-2xl border border-slate-100 bg-slate-50 p-5">
                                        <span className="inline-flex rounded-full bg-blue-100 px-3 py-1 text-xs font-semibold text-blue-700">
                                            {product.category || "General"}
                                        </span>
                                        <h4 className="mt-4 text-lg font-semibold text-slate-900">{product.name}</h4>
                                        <p className="mt-2 text-sm text-slate-500">{product.reason}</p>
                                        <p className="mt-4 text-2xl font-bold text-slate-900">{currencyFormatter.format(product.price)}</p>
                                        <button
                                            type="button"
                                            onClick={() => handleBuyAgain(product.productId)}
                                            disabled={buyingProductId === product.productId}
                                            className="mt-5 rounded-full bg-blue-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-800 disabled:opacity-60"
                                        >
                                            {buyingProductId === product.productId ? "Adding..." : "Buy Again"}
                                        </button>
                                    </article>
                                ))}
                            </div>
                        )}
                    </section>
                </>
            )}
        </DashboardShell>
    );
};
