import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCart } from "../context/CartContext";
import { DashboardShell } from "./dashboard/DashboardShell";

export const Cart = () => {
    const navigate = useNavigate();
    const role = (localStorage.getItem("role") || "user").toLowerCase();
    const { cartItems, loading, updateItem, removeItem } = useCart();
    const [errorMessage, setErrorMessage] = useState<string>("");

    const totalAmount = useMemo(() => {
        return cartItems.reduce((total, item) => total + item.quantity * item.product_price, 0);
    }, [cartItems]);

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

    const handleIncrease = async (productId: number, quantity: number, stock: number) => {
        setErrorMessage("");
        if (quantity + 1 > stock) {
            setErrorMessage("Insufficient stock for this product.");
            return;
        }

        try {
            await updateItem(productId, quantity + 1);
        } catch (error) {
            setErrorMessage((error as Error).message || "Failed to update cart item.");
        }
    };

    const handleDecrease = async (productId: number, quantity: number) => {
        setErrorMessage("");
        try {
            if (quantity - 1 <= 0) {
                await removeItem(productId);
                return;
            }
            await updateItem(productId, quantity - 1);
        } catch (error) {
            setErrorMessage((error as Error).message || "Failed to update cart item.");
        }
    };

    return (
        <DashboardShell
            role={role}
            activeTab="cart"
            title="Cart"
            subtitle="Review your items before checkout"
            onLogout={handleLogout}
            contentVariant="stack"
            showSidebar={false}
            fullWidth={true}
        >
            {errorMessage && <div className="rounded-lg border border-red-200 bg-red-50 text-red-700 px-4 py-3 text-sm font-medium">{errorMessage}</div>}

            <div className="w-full bg-white rounded-2xl shadow-md p-5">
                {loading ? (
                    <p className="text-sm text-gray-500">Loading cart...</p>
                ) : cartItems.length === 0 ? (
                    <p className="text-sm text-gray-500">Your cart is empty.</p>
                ) : (
                    <div className="space-y-3">
                        {cartItems.map((item) => (
                            <div key={item.id} className="grid grid-cols-1 md:grid-cols-5 gap-3 items-center rounded-xl border border-gray-100 p-3">
                                <p className="text-sm font-semibold text-gray-900 md:col-span-2">{item.product_name}</p>
                                <p className="text-sm text-gray-700">${item.product_price.toFixed(2)}</p>
                                <div className="flex items-center gap-2">
                                    <button 
                                        type="button"
                                        onClick={() => handleDecrease(item.product_id, item.quantity)} 
                                        className="w-8 h-8 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200"
                                    >
                                        −
                                    </button>
                                    <span className="min-w-6 text-center text-sm font-semibold">{item.quantity}</span>
                                    <button 
                                        type="button"
                                        onClick={() => handleIncrease(item.product_id, item.quantity, item.product_stock)} 
                                        className="w-8 h-8 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200"
                                    >
                                        +
                                    </button>
                                </div>
                                <div className="flex items-center justify-between md:justify-end gap-3">
                                    <p className="text-sm font-semibold text-gray-900">${(item.quantity * item.product_price).toFixed(2)}</p>
                                    <button 
                                        type="button"
                                        onClick={() => removeItem(item.product_id)} 
                                        className="px-2 py-1 rounded-md bg-red-100 text-red-700 text-xs font-medium hover:bg-red-200"
                                    >
                                        Remove
                                    </button>
                                </div>
                            </div>
                        ))}

                        <div className="pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">
                            <p className="text-base font-semibold text-gray-900">Total: ${totalAmount.toFixed(2)}</p>
                            <button 
                                type="button"
                                onClick={() => navigate("/checkout")} 
                                disabled={cartItems.length === 0} 
                                className="rounded-lg bg-emerald-600 text-white px-4 py-2 font-semibold hover:bg-emerald-700 disabled:opacity-60"
                            >
                                Checkout
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </DashboardShell>
    );
};
