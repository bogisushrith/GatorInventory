import { useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { FilterPanel } from "./dashboard/FilterPanel.tsx";
import { Pagination } from "./dashboard/Pagination.tsx";
import { ProductTable } from "./dashboard/ProductTable.tsx";
import { StatsCards } from "./dashboard/StatsCards.tsx";
import { Product, ProductListResponse, ProductPagination } from "./dashboard/types.ts";
import { DashboardShell } from "./dashboard/DashboardShell.tsx";
 
export const Dashboard = () => {
    const DEFAULT_LIMIT = 5;
    const role = (localStorage.getItem("role") || "user").toLowerCase();
    const isAdmin = role === "admin";

    const navigate = useNavigate();
    const [searchParams, setSearchParams] = useSearchParams();

    const initialPage = useMemo(() => {
        const rawPage = Number(searchParams.get("page") || "1");
        if (Number.isNaN(rawPage) || rawPage < 1) {
            return 1;
        }
        return rawPage;
    }, [searchParams]);

    const initialLimit = useMemo(() => {
        const rawLimit = Number(searchParams.get("limit") || String(DEFAULT_LIMIT));
        const allowedLimits = [5, 10, 25];
        if (allowedLimits.includes(rawLimit)) {
            return rawLimit;
        }
        return DEFAULT_LIMIT;
    }, [searchParams]);

    const initialSearch = searchParams.get("search") || "";
    const initialCategory = searchParams.get("category") || "";
    const initialMinPrice = searchParams.get("min_price") || "";
    const initialMaxPrice = searchParams.get("max_price") || "";

    const [page, setPage] = useState<number>(initialPage);
    const [limit, setLimit] = useState<number>(initialLimit);
    const [searchInput, setSearchInput] = useState<string>(initialSearch);
    const [categoryInput, setCategoryInput] = useState<string>(initialCategory);
    const [minPriceInput, setMinPriceInput] = useState<string>(initialMinPrice);
    const [maxPriceInput, setMaxPriceInput] = useState<string>(initialMaxPrice);

    const [search, setSearch] = useState<string>(initialSearch);
    const [category, setCategory] = useState<string>(initialCategory);
    const [minPrice, setMinPrice] = useState<string>(initialMinPrice);
    const [maxPrice, setMaxPrice] = useState<string>(initialMaxPrice);

    const [products, setProducts] = useState<Product[]>([]);
    const [categories, setCategories] = useState<string[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [refreshKey, setRefreshKey] = useState<number>(0);
    const [pagination, setPagination] = useState<ProductPagination>({
        page: 1,
        limit: DEFAULT_LIMIT,
        total: 0,
        total_pages: 0
    });

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                credentials: 'include'
            });
        } catch (error) {
            console.error("Error:", error);
        } finally {
            localStorage.removeItem("role");
            localStorage.removeItem("username");
            navigate("/");
        }
    };

    useEffect(() => {
        const params = new URLSearchParams();
        params.set("page", String(page));
        params.set("limit", String(limit));

        if (search.trim()) {
            params.set("search", search.trim());
        }
        if (category.trim()) {
            params.set("category", category.trim());
        }
        if (minPrice.trim()) {
            params.set("min_price", minPrice.trim());
        }
        if (maxPrice.trim()) {
            params.set("max_price", maxPrice.trim());
        }

        setSearchParams(params, { replace: true });
    }, [page, limit, search, category, minPrice, maxPrice, setSearchParams]);

    useEffect(() => {
        const fetchProducts = async () => {
            setLoading(true);
            try {
                const params = new URLSearchParams({
                    page: String(page),
                    limit: String(limit)
                });

                if (search.trim()) {
                    params.set("search", search.trim());
                }
                if (category.trim()) {
                    params.set("category", category.trim());
                }
                if (minPrice.trim()) {
                    params.set("min_price", minPrice.trim());
                }
                if (maxPrice.trim()) {
                    params.set("max_price", maxPrice.trim());
                }

                const response = await fetch(`/api/products?${params.toString()}`, {
                    credentials: 'include'
                });

                if (response.ok) {
                    const result = await response.json() as ProductListResponse | Product[];

                    if (Array.isArray(result)) {
                        setProducts(result);
                        setPagination({
                            page,
                            limit,
                            total: result.length,
                            total_pages: result.length > 0 ? 1 : 0
                        });
                        setCategories((previous) => Array.from(new Set([...previous, ...result.map((product) => product.category)])));
                        return;
                    }

                    setProducts(result.data || []);
                    setPagination(result.pagination || { page, limit, total: 0, total_pages: 0 });
                    setCategories((previous) => Array.from(new Set([...previous, ...(result.data || []).map((product) => product.category)])));

                    if (result.pagination && result.pagination.total_pages > 0 && page > result.pagination.total_pages) {
                        setPage(result.pagination.total_pages);
                    }
                } else {
                    console.error("Failed to fetch products:", response.status);
                }
            } catch (error) {
                console.error("Error fetching products:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, [page, limit, search, category, minPrice, maxPrice, refreshKey]);

    const handleSearchChange = (nextSearch: string) => {
        setSearchInput(nextSearch);
    };

    const handleCategoryChange = (nextCategory: string) => {
        setCategoryInput(nextCategory);
    };

    const handleMinPriceChange = (nextMinPrice: string) => {
        setMinPriceInput(nextMinPrice);
    };

    const handleMaxPriceChange = (nextMaxPrice: string) => {
        setMaxPriceInput(nextMaxPrice);
    };

    const handleResetFilters = () => {
        setSearchInput("");
        setCategoryInput("");
        setMinPriceInput("");
        setMaxPriceInput("");

        setSearch("");
        setCategory("");
        setMinPrice("");
        setMaxPrice("");
        setLimit(DEFAULT_LIMIT);
        setPage(1);
    };

    const handleApplyFilters = () => {
        setSearch(searchInput);
        setCategory(categoryInput);
        setMinPrice(minPriceInput);
        setMaxPrice(maxPriceInput);
        setPage(1);
    };

    const handleLimitChange = (nextLimit: number) => {
        setLimit(nextLimit);
        setPage(1);
    };

    const hasPendingChanges = useMemo(() => {
        return (
            searchInput.trim() !== search.trim() ||
            categoryInput.trim() !== category.trim() ||
            minPriceInput.trim() !== minPrice.trim() ||
            maxPriceInput.trim() !== maxPrice.trim()
        );
    }, [searchInput, categoryInput, minPriceInput, maxPriceInput, search, category, minPrice, maxPrice]);

    const handleUpdate = (product: Product) => {
        navigate("/update-product", { state: { product } });
    };

    const handleDelete = async (id: number) => {
        try {
            const response = await fetch(`/api/products/${id}`, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json"
                },
                credentials: 'include'
            });

            if (response.ok) {
                setRefreshKey((previous) => previous + 1);
            } else {
                const result = await response.json();
                console.error("Delete failed:", result.error_message);
            }
        } catch (error) {
            console.error("Error:", error);
        }
    };

    return (
        <DashboardShell
            role={role}
            activeTab="products"
            title="Inventory Dashboard"
            subtitle="Search, filter, and manage products"
            onLogout={handleLogout}
            contentVariant="grid"
            showSidebar={false}
        >
            {isAdmin && (
                <div className="lg:col-span-4">
                    <button
                        onClick={() => navigate("/add-product")}
                        className="btn-secondary flex items-center gap-2 transition-all duration-200 hover:-translate-y-0.5 disabled:opacity-60 disabled:cursor-not-allowed"
                        disabled={loading}
                    >
                        ➕ Add Product
                    </button>
                </div>
            )}

            <div className="lg:col-span-1 h-full">
                <FilterPanel
                    searchInput={searchInput}
                    onSearchChange={handleSearchChange}
                    categoryInput={categoryInput}
                    onCategoryChange={handleCategoryChange}
                    minPriceInput={minPriceInput}
                    onMinPriceChange={handleMinPriceChange}
                    maxPriceInput={maxPriceInput}
                    onMaxPriceChange={handleMaxPriceChange}
                    categories={categories}
                    onReset={handleResetFilters}
                    onApply={handleApplyFilters}
                    hasPendingChanges={hasPendingChanges}
                    loading={loading}
                />
            </div>

            <div className="lg:col-span-3 min-w-0 flex flex-col gap-6 h-full">
                <ProductTable
                    products={products}
                    loading={loading}
                    page={pagination.page || page}
                    limit={pagination.limit || limit}
                    totalResults={pagination.total || 0}
                    onLimitChange={handleLimitChange}
                    canManageProducts={isAdmin}
                    onUpdate={handleUpdate}
                    onDelete={handleDelete}
                />

                <Pagination
                    page={pagination.page || page}
                    totalPages={pagination.total_pages || 0}
                    loading={loading}
                    onPageChange={setPage}
                />
            </div>

            <div className="lg:col-span-4">
                <StatsCards
                    products={products}
                    totalResults={pagination.total || 0}
                />
            </div>
        </DashboardShell>
    );
};
