export interface Product {
    id: number;
    name: string;
    price: number;
    quantity: number;
    category: string;
}

export interface ProductPagination {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
}

export interface ProductListResponse {
    data: Product[];
    pagination: ProductPagination;
}

export interface CartItem {
    id: number;
    user_id: number;
    product_id: number;
    quantity: number;
    product_name: string;
    product_price: number;
    product_stock: number;
}

export interface OrderItem {
    id: number;
    order_id: number;
    product_id: number;
    product_name: string;
    product_price: number;
    quantity: number;
}

export type OrderStatus = "pending" | "completed" | "cancelled";

export interface Order {
    id: number;
    user_id: number;
    user_name: string;
    status: OrderStatus | string;
    created_at: string;
    items: OrderItem[];
}

export interface OrderCreateItem {
    product_id: number;
    quantity: number;
}

export interface OrderQueryParams {
    search?: string;
    status?: string;
    date_from?: string;
    date_to?: string;
}
