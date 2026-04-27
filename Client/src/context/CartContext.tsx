import { createContext, ReactNode, useContext, useEffect, useMemo, useState } from "react";
import { addToCart, getCart, removeCartItem, updateCartItem } from "../api/cart";
import { CartItem } from "../components/dashboard/types";

const LEGACY_CART_STORAGE_KEY = "inventory_cart";
const CART_STORAGE_PREFIX = "inventory_cart_";

type CartContextValue = {
    cartItems: CartItem[];
    loading: boolean;
    refreshCart: () => Promise<void>;
    addItem: (productId: number, quantity: number) => Promise<void>;
    updateItem: (productId: number, quantity: number) => Promise<void>;
    removeItem: (productId: number) => Promise<void>;
    clearClientCart: () => void;
    totalQuantity: number;
};

const CartContext = createContext<CartContextValue | null>(null);

export const CartProvider = ({ children }: { children: ReactNode }) => {
    const [cartItems, setCartItems] = useState<CartItem[]>([]);
    const [loading, setLoading] = useState<boolean>(false);

    const isAuthenticated = () => {
        return Boolean(localStorage.getItem("role"));
    };

    const getCurrentUsername = () => {
        return (localStorage.getItem("username") || "").trim().toLowerCase();
    };

    const getCartStorageKey = () => {
        const username = getCurrentUsername();
        if (!username) {
            return null;
        }
        return `${CART_STORAGE_PREFIX}${username}`;
    };

    const saveCartToStorage = (items: CartItem[]) => {
        try {
            const storageKey = getCartStorageKey();
            if (!storageKey) {
                return;
            }
            localStorage.setItem(storageKey, JSON.stringify(items));
            localStorage.removeItem(LEGACY_CART_STORAGE_KEY);
        } catch (error) {
            console.error("Failed to save cart to localStorage:", error);
        }
    };

    const loadCartFromStorage = (): CartItem[] => {
        try {
            const storageKey = getCartStorageKey();
            if (!storageKey) {
                return [];
            }
            const stored = localStorage.getItem(storageKey);
            return stored ? JSON.parse(stored) : [];
        } catch (error) {
            console.error("Failed to load cart from localStorage:", error);
            return [];
        }
    };

    const clearCartStorage = () => {
        const storageKey = getCartStorageKey();
        if (storageKey) {
            localStorage.removeItem(storageKey);
        }
        localStorage.removeItem(LEGACY_CART_STORAGE_KEY);
    };

    const refreshCart = async () => {
        if (!isAuthenticated()) {
            setCartItems([]);
            clearCartStorage();
            return;
        }

        // Switch to current user's local snapshot immediately before API sync.
        setCartItems(loadCartFromStorage());

        setLoading(true);
        try {
            const items = await getCart();
            setCartItems(items);
            saveCartToStorage(items);
        } catch {
            // Keep cart state stable when fetch fails; use localStorage fallback
            const stored = loadCartFromStorage();
            setCartItems(stored);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        // On mount, load from localStorage first, then sync with backend
        const stored = loadCartFromStorage();
        if (stored.length > 0) {
            setCartItems(stored);
        }
        if (isAuthenticated()) {
            refreshCart();
        }
    }, []);

    const addItem = async (productId: number, quantity: number) => {
        await addToCart(productId, quantity);
        await refreshCart();
    };

    const updateItem = async (productId: number, quantity: number) => {
        await updateCartItem(productId, quantity);
        await refreshCart();
    };

    const removeItem = async (productId: number) => {
        await removeCartItem(productId);
        await refreshCart();
    };

    const clearClientCart = () => {
        setCartItems([]);
        clearCartStorage();
    };

    const totalQuantity = useMemo(() => {
        return cartItems.reduce((total, item) => total + item.quantity, 0);
    }, [cartItems]);

    return (
        <CartContext.Provider value={{ cartItems, loading, refreshCart, addItem, updateItem, removeItem, clearClientCart, totalQuantity }}>
            {children}
        </CartContext.Provider>
    );
};

export const useCart = () => {
    const context = useContext(CartContext);
    if (!context) {
        throw new Error("useCart must be used within CartProvider");
    }
    return context;
};
 
