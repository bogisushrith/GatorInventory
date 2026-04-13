import { Product } from "../components/dashboard/types";

export const updateProductStock = async (productId: number, quantity: number): Promise<Product> => {
    const response = await fetch(`/api/products/${productId}/stock`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({ quantity })
    });

    if (!response.ok) {
        let errorMessage = "Failed to update stock";
        try {
            const errorData = await response.json() as { error_message?: string };
            if (errorData.error_message) {
                errorMessage = errorData.error_message;
            }
        } catch {
            // Keep default error message when body is not JSON.
        }
        throw new Error(errorMessage);
    }

    return response.json() as Promise<Product>;
};
