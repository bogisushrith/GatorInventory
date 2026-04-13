import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createOrder } from "../api/orders";
import { useCart } from "../context/CartContext";
import { DashboardShell } from "./dashboard/DashboardShell";

export const Checkout = () => {
    const navigate = useNavigate();
    const role = (localStorage.getItem("role") || "user").toLowerCase();
    const { cartItems, refreshCart, clearClientCart } = useCart();

    const [placingOrder, setPlacingOrder] = useState<boolean>(false);
    const [errorMessage, setErrorMessage] = useState<string>("");

    const totalAmount = useMemo(() => cartItems.reduce((total, item) => total + item.quantity * item.product_price, 0), [cartItems]);

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

    const handlePlaceOrder = async () => {
        setErrorMessage("");
        if (cartItems.length === 0) {
            setErrorMessage("Your cart is empty.");
            return;
        }

        setPlacingOrder(true);
        try {
            await createOrder(cartItems.map((item) => ({ product_id: item.product_id, quantity: item.quantity })));
            clearClientCart();
            await refreshCart();
            navigate("/history");
        } catch (error) {
            setErrorMessage((error as Error).message || "Failed to place order.");
        } finally {
            setPlacingOrder(false);
        }
    };

    return (
        <DashboardShell
            role={role}
            activeTab="checkout"
            title="Checkout"
            subtitle="Confirm order details and place your order"
            onLogout={handleLogout}
            contentVariant="grid"
            showSidebar={true}
        >
            {errorMessage && <div className="lg:col-span-3 rounded-lg border border-red-200 bg-red-50 text-red-700 px-4 py-3 text-sm font-medium">{errorMessage}</div>}

            <div className="lg:col-span-3 bg-white rounded-2xl shadow-md p-5">
                {cartItems.length === 0 ? (
                    <p className="text-sm text-gray-500">No items to checkout.</p>
                ) : (
                    <div className="space-y-3">
                        {cartItems.map((item) => (
                            <div key={item.id} className="flex items-center justify-between rounded-xl border border-gray-100 p-3">
                                <div>
                                    <p className="text-sm font-semibold text-gray-900">{item.product_name}</p>
                                    <p className="text-xs text-gray-500">Qty: {item.quantity}</p>
                                </div>
                                <p className="text-sm font-semibold text-gray-900">${(item.quantity * item.product_price).toFixed(2)}</p>
                            </div>
                        ))}

                        <div className="pt-4 border-t border-gray-100 flex items-center justify-between">
                            <p className="text-base font-semibold text-gray-900">Total: ${totalAmount.toFixed(2)}</p>
                            <button type="button" onClick={handlePlaceOrder} disabled={placingOrder || cartItems.length === 0} className="rounded-lg bg-indigo-600 text-white px-4 py-2 font-semibold hover:bg-indigo-700 disabled:opacity-60">
                                {placingOrder ? "Placing Order..." : "Place Order"}
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </DashboardShell>
    );
};
