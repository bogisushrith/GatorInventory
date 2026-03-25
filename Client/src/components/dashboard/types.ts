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
