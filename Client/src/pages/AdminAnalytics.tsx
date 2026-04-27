import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { AnalyticsBundle, AnalyticsDays, fallbackAnalyticsBundle, getAnalyticsBundle } from "../api/analytics";
import { DateRangeFilter } from "../components/analytics/DateRangeFilter";
import { LowStockTable } from "../components/analytics/LowStockTable";
import { OrdersChartCard } from "../components/analytics/OrdersChartCard";
import { RevenueChartCard } from "../components/analytics/RevenueChartCard";
import { SummaryCards } from "../components/analytics/SummaryCards";
import { TopProductsTable } from "../components/analytics/TopProductsTable";
import { DashboardShell } from "../components/dashboard/DashboardShell";

const emptyBundle: AnalyticsBundle = {
    summary: {
        totalRevenue: 0,
        totalOrders: 0,
        totalProducts: 0,
        lowStockCount: 0,
    },
    revenueTrend: [],
    orderTrend: [],
    topProducts: [],
    lowStockProducts: [],
};

export const AdminAnalytics = () => {
    const navigate = useNavigate();
    const role = (localStorage.getItem("role") || "").toLowerCase();
    const isAdmin = role === "admin";

    const [days, setDays] = useState<AnalyticsDays>(30);
    const [loading, setLoading] = useState<boolean>(false);
    const [errorMessage, setErrorMessage] = useState<string>("");
    const [bundle, setBundle] = useState<AnalyticsBundle>(emptyBundle);

    const loadAnalytics = useCallback(async (selectedDays: AnalyticsDays) => {
        setLoading(true);
        setErrorMessage("");
        try {
            const data = await getAnalyticsBundle(selectedDays);
            setBundle({
                summary: data.summary,
                revenueTrend: data.revenueTrend.length > 0 ? data.revenueTrend : fallbackAnalyticsBundle.revenueTrend,
                orderTrend: data.orderTrend.length > 0 ? data.orderTrend : fallbackAnalyticsBundle.orderTrend,
                topProducts: data.topProducts.length > 0 ? data.topProducts : fallbackAnalyticsBundle.topProducts,
                lowStockProducts: data.lowStockProducts,
            });
        } catch (error) {
            setErrorMessage((error as Error).message || "Unable to load live analytics. Showing sample data.");
            setBundle(fallbackAnalyticsBundle);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        if (!isAdmin) {
            navigate("/dashboard", { replace: true });
            return;
        }
        void loadAnalytics(days);
    }, [days, isAdmin, loadAnalytics, navigate]);

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", {
                method: "POST",
                credentials: "include",
            });
        } finally {
            localStorage.removeItem("token");
            localStorage.removeItem("role");
            localStorage.removeItem("username");
            navigate("/");
        }
    };

    return (
        <DashboardShell
            role={role}
            activeTab="analytics"
            title="Admin Analytics"
            subtitle="Revenue, order trends, top products, and low-stock insights"
            onLogout={handleLogout}
            contentVariant="stack"
            headerActions={<DateRangeFilter value={days} loading={loading} onChange={setDays} />}
        >
            <div className="space-y-6">
                {errorMessage && (
                    <div className="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm font-medium text-amber-800">
                        {errorMessage}
                    </div>
                )}

                <SummaryCards summary={bundle.summary} loading={loading} />

                <div className="grid grid-cols-1 gap-6 xl:grid-cols-2">
                    <RevenueChartCard data={bundle.revenueTrend} loading={loading} />
                    <OrdersChartCard data={bundle.orderTrend} loading={loading} />
                </div>

                <div className="grid grid-cols-1 gap-6 xl:grid-cols-2">
                    <TopProductsTable data={bundle.topProducts} loading={loading} />
                    <LowStockTable data={bundle.lowStockProducts} loading={loading} />
                </div>

                {!loading && bundle.summary.totalOrders === 0 && bundle.topProducts.length === 0 && (
                    <div className="rounded-2xl border border-dashed border-slate-300 bg-white px-6 py-10 text-center shadow-md">
                        <p className="text-lg font-semibold text-slate-800">No data available yet</p>
                        <p className="mt-2 text-sm text-slate-500">Analytics will populate here as soon as products and orders are available.</p>
                    </div>
                )}
            </div>
        </DashboardShell>
    );
};
