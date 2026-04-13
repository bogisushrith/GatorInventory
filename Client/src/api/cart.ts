import { CartItem } from "../components/dashboard/types";

type ApiErrorPayload = {
    error_message?: string;
    message?: string;
};

const readApiErrorMessage = async (response: Response, fallbackMessage: string): Promise<string> => {
    try {
        const errorData = await response.json() as ApiErrorPayload;
        return errorData.error_message || errorData.message || fallbackMessage;
    } catch {
        try {
            const rawText = await response.text();
            if (rawText.trim().length > 0) {
                return rawText;
            }
        } catch {
            // Ignore and keep fallback.
        }
        return fallbackMessage;
    }
};

const fetchCartEndpoint = async (path: string, init?: RequestInit): Promise<Response> => {
    const apiPath = `/api${path}`;
    const directPath = path;

    const firstResponse = await fetch(apiPath, init);
    if (firstResponse.status !== 404) {
        return firstResponse;
    }

    const fallbackMessage = await readApiErrorMessage(firstResponse, "Not Found");
    if (!fallbackMessage.toLowerCase().includes("not found")) {
        return firstResponse;
    }

    return fetch(directPath, init);
};

export const getCart = async (): Promise<CartItem[]> => {
    const response = await fetchCartEndpoint("/cart", { credentials: "include" });
    if (!response.ok) {
        const errorMessage = await readApiErrorMessage(response, "Failed to load cart");
        throw new Error(errorMessage);
    }
    return response.json() as Promise<CartItem[]>;
};

export const addToCart = async (productId: number, quantity: number): Promise<void> => {
    const response = await fetchCartEndpoint("/cart/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ product_id: productId, quantity })
    });

    if (!response.ok) {
        const errorMessage = await readApiErrorMessage(response, "Failed to add item to cart");
        throw new Error(errorMessage);
    }
};

export const updateCartItem = async (productId: number, quantity: number): Promise<void> => {
    const response = await fetchCartEndpoint("/cart/update", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ product_id: productId, quantity })
    });

    if (!response.ok) {
        const errorMessage = await readApiErrorMessage(response, "Failed to update cart item");
        throw new Error(errorMessage);
    }
};

export const removeCartItem = async (productId: number): Promise<void> => {
    const response = await fetchCartEndpoint("/cart/remove", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ product_id: productId })
    });

    if (!response.ok) {
        const errorMessage = await readApiErrorMessage(response, "Failed to remove cart item");
        throw new Error(errorMessage);
    }
};
