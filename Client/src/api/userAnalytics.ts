import { getOrders } from "./orders";
import { Order, UserAnalyticsBundle, UserAnalyticsSummary, UserRecentOrder, UserRecommendation, UserTopProduct } from "../components/dashboard/types";

type TrendPoint = {
    date: string;
    value: number;
};

const parseErrorMessage = async (response: Response, fallback: string) => {
    try {
        const errorPayload = await response.json() as { error_message?: string; message?: string };
        return errorPayload.error_message || errorPayload.message || fallback;
    } catch {
        return fallback;
    }
};

const getJson = async <T>(path: string): Promise<T> => {
    const response = await fetch(path, { credentials: "include" });
    if (!response.ok) {
        throw new Error(await parseErrorMessage(response, "Failed to load user analytics"));
    }
    return response.json() as Promise<T>;
};

const calculateOrderTotal = (order: Order) => (
    order.items.reduce((sum, item) => sum + (item.quantity * item.product_price), 0)
);

const buildFallbackBundleFromOrders = async (): Promise<Partial<UserAnalyticsBundle>> => {
    const orders = await getOrders();

    const totalSpent = orders.reduce((sum, order) => sum + calculateOrderTotal(order), 0);
    const pendingOrders = orders.filter((order) => order.status.trim().toLowerCase() === "pending").length;

    const recentOrders: UserRecentOrder[] = orders.slice(0, 5).map((order) => ({
        orderId: order.id,
        productNames: order.items.map((item) => item.product_name).filter(Boolean).join(", ") || "No items",
        status: order.status,
        createdAt: order.created_at,
        totalPrice: calculateOrderTotal(order),
    }));

    const productCounts = new Map<number, { name: string; purchaseCount: number }>();
    orders.forEach((order) => {
        order.items.forEach((item) => {
            const existing = productCounts.get(item.product_id);
            if (existing) {
                existing.purchaseCount += item.quantity;
                return;
            }
            productCounts.set(item.product_id, {
                name: item.product_name || `Product #${item.product_id}`,
                purchaseCount: item.quantity,
            });
        });
    });

    const topProducts: UserTopProduct[] = Array.from(productCounts.entries())
        .map(([productId, item]) => ({
            productId,
            name: item.name,
            purchaseCount: item.purchaseCount,
        }))
        .sort((left, right) => right.purchaseCount - left.purchaseCount || left.name.localeCompare(right.name))
        .slice(0, 5);

    const spendingByDate = new Map<string, number>();
    orders.forEach((order) => {
        const key = new Date(order.created_at).toISOString().slice(0, 10);
        spendingByDate.set(key, (spendingByDate.get(key) || 0) + calculateOrderTotal(order));
    });

    const spendingTrend = Array.from(spendingByDate.entries())
        .map(([date, value]) => ({ date, value }))
        .sort((left, right) => left.date.localeCompare(right.date));

    return {
        summary: {
            totalOrders: orders.length,
            totalSpent,
            pendingOrders,
        },
        recentOrders,
        topProducts,
        spendingTrend,
    };
};

export const getUserAnalyticsBundle = async (): Promise<UserAnalyticsBundle> => {
    const analyticsResults = await Promise.allSettled([
        getJson<UserAnalyticsSummary>("/api/user/analytics/summary"),
        getJson<UserRecentOrder[]>("/api/user/analytics/recent-orders"),
        getJson<UserTopProduct[]>("/api/user/analytics/top-products"),
        getJson<TrendPoint[]>("/api/user/analytics/spending-trend"),
        getJson<UserRecommendation[]>("/api/user/analytics/recommendations"),
    ]);

    const fallbackBundle = analyticsResults.some((result) => result.status === "rejected")
        ? await buildFallbackBundleFromOrders()
        : {};

    const summaryResult = analyticsResults[0];
    const recentOrdersResult = analyticsResults[1];
    const topProductsResult = analyticsResults[2];
    const spendingTrendResult = analyticsResults[3];
    const recommendationsResult = analyticsResults[4];

    return {
        summary: summaryResult.status === "fulfilled"
            ? summaryResult.value
            : (fallbackBundle.summary || { totalOrders: 0, totalSpent: 0, pendingOrders: 0 }),
        recentOrders: recentOrdersResult.status === "fulfilled"
            ? recentOrdersResult.value
            : (fallbackBundle.recentOrders || []),
        topProducts: topProductsResult.status === "fulfilled"
            ? topProductsResult.value
            : (fallbackBundle.topProducts || []),
        spendingTrend: spendingTrendResult.status === "fulfilled"
            ? spendingTrendResult.value
            : (fallbackBundle.spendingTrend || []),
        recommendations: recommendationsResult.status === "fulfilled"
            ? recommendationsResult.value
            : [],
    };
};
