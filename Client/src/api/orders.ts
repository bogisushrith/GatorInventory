import { Order, OrderCreateItem, OrderQueryParams } from "../components/dashboard/types";

export const createOrder = async (items: OrderCreateItem[] = []): Promise<{ order_id: number }> => {
    const response = await fetch("/api/orders", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({ items })
    });

    if (!response.ok) {
        let errorMessage = "Failed to place order";
        try {
            const errorData = await response.json() as { error_message?: string };
            if (errorData.error_message) {
                errorMessage = errorData.error_message;
            }
        } catch {
            // Keep fallback message when response body is not JSON.
        }
        throw new Error(errorMessage);
    }

    return response.json() as Promise<{ order_id: number }>;
};

export const getOrders = async (query: OrderQueryParams = {}): Promise<Order[]> => {
    const queryParams = new URLSearchParams();

    Object.entries(query).forEach(([key, value]) => {
        if (value && value.trim().length > 0) {
            queryParams.set(key, value.trim());
        }
    });

    const response = await fetch(`/api/orders${queryParams.toString().length > 0 ? `?${queryParams.toString()}` : ""}`, {
        credentials: "include"
    });

    if (!response.ok) {
        throw new Error("Failed to load orders");
    }

    return response.json() as Promise<Order[]>;
};

export const getOrderById = async (id: number): Promise<Order> => {
    const response = await fetch(`/api/orders/${id}`, {
        credentials: "include"
    });

    if (!response.ok) {
        throw new Error("Failed to load order details");
    }

    return response.json() as Promise<Order>;
};
