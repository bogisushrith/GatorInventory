import { Order, OrderStatus } from "../dashboard/types";

export const calculateOrderTotal = (order: Order): number => {
    return order.items.reduce((sum, item) => sum + (item.quantity * item.product_price), 0);
};

export const formatOrderDate = (dateValue: string): string => {
    return new Date(dateValue).toLocaleString();
};

export const getStatusStyles = (status: string) => {
    const normalizedStatus = status.toLowerCase() as OrderStatus;

    switch (normalizedStatus) {
        case "completed":
            return "bg-emerald-100 text-emerald-700 border-emerald-200";
        case "cancelled":
            return "bg-red-100 text-red-700 border-red-200";
        default:
            return "bg-amber-100 text-amber-700 border-amber-200";
    }
};

export const formatStatusLabel = (status: string): string => {
    return status.charAt(0).toUpperCase() + status.slice(1).toLowerCase();
};
