import { Product } from "./types";

interface ProductTableProps {
    products: Product[];
    loading: boolean;
    page: number;
    limit: number;
    totalResults: number;
    onLimitChange: (nextLimit: number) => void;
    canManageProducts: boolean;
    onUpdate: (product: Product) => void;
    onDelete: (id: number) => void;
}

const LoadingSkeleton = () => {
    return (
        <div className="animate-pulse space-y-3 p-4">
            {Array.from({ length: 5 }).map((_, index) => (
                <div key={index} className="h-12 bg-gray-100 rounded-lg" />
            ))}
        </div>
    );
};

export const ProductTable = ({ products, loading, page, limit, totalResults, onLimitChange, canManageProducts, onUpdate, onDelete }: ProductTableProps) => {
    const start = totalResults > 0 ? ((page - 1) * limit) + 1 : 0;
    const end = totalResults > 0 ? ((page - 1) * limit) + products.length : 0;

    return (
        <div className="bg-white rounded-2xl shadow-md overflow-hidden transition-all duration-200 h-full">
            <div className="p-5 border-b border-gray-100">
                <div className="flex justify-between items-center gap-4">
                    <h3 className="text-xl font-semibold text-gray-900">Products</h3>
                    <div className="flex items-center gap-4 text-sm text-gray-500">
                        <p className="whitespace-nowrap">Showing {start}–{end} of {totalResults} results</p>
                        <div className="flex items-center gap-2">
                            <label className="whitespace-nowrap">Items per page</label>
                            <select
                                value={limit}
                                onChange={(event) => onLimitChange(Number(event.target.value))}
                                className="rounded-lg border border-gray-200 px-2.5 py-1.5 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-purple-400 transition-all duration-200 bg-white"
                                disabled={loading}
                            >
                                <option value={5}>5</option>
                                <option value={10}>10</option>
                                <option value={25}>25</option>
                            </select>
                        </div>
                    </div>
                </div>
            </div>

            <div className="overflow-x-auto">
                {loading ? (
                    <LoadingSkeleton />
                ) : (
                    <table className="w-full min-w-[760px]">
                        <thead>
                            <tr className="bg-gradient-to-r from-indigo-600 to-pink-500 text-white">
                                <th className="px-6 py-4 text-left text-xs font-bold uppercase tracking-wide">Name</th>
                                <th className="px-6 py-4 text-left text-xs font-bold uppercase tracking-wide">Price</th>
                                <th className="px-6 py-4 text-left text-xs font-bold uppercase tracking-wide">Quantity</th>
                                <th className="px-6 py-4 text-left text-xs font-bold uppercase tracking-wide">Total</th>
                                <th className="px-6 py-4 text-left text-xs font-bold uppercase tracking-wide">Category</th>
                                {canManageProducts && <th className="px-6 py-4 text-left text-xs font-bold uppercase tracking-wide">Actions</th>}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {products.length === 0 ? (
                                <tr>
                                    <td colSpan={canManageProducts ? 6 : 5} className="px-6 py-14 text-center">
                                        <p className="text-gray-500 text-lg font-medium">No products found. Try adjusting filters.</p>
                                    </td>
                                </tr>
                            ) : (
                                products.map((product) => (
                                    <tr key={product.id} className="hover:bg-gray-50 transition-all duration-200">
                                        <td className="px-6 py-4 text-sm font-semibold text-gray-900">{product.name}</td>
                                        <td className="px-6 py-4 text-sm text-gray-700">
                                            {new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(product.price)}
                                        </td>
                                        <td className="px-6 py-4 text-sm text-gray-700">{product.quantity}</td>
                                        <td className="px-6 py-4 text-sm text-gray-700">
                                            {new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(product.price * product.quantity)}
                                        </td>
                                        <td className="px-6 py-4 text-sm">
                                            <span className="inline-flex rounded-full px-3 py-1 text-sm font-semibold bg-cyan-100 text-cyan-800">
                                                {product.category}
                                            </span>
                                        </td>
                                        {canManageProducts && (
                                            <td className="px-6 py-4 text-sm">
                                                <div className="flex items-center gap-2">
                                                    <button
                                                        onClick={() => onUpdate(product)}
                                                        className="px-3 py-1.5 rounded-lg bg-indigo-100 text-indigo-700 font-medium hover:bg-indigo-200 transition-all duration-200"
                                                    >
                                                        Edit
                                                    </button>
                                                    <button
                                                        onClick={() => onDelete(product.id)}
                                                        className="px-3 py-1.5 rounded-lg bg-red-100 text-red-700 font-medium hover:bg-red-200 transition-all duration-200"
                                                    >
                                                        Delete
                                                    </button>
                                                </div>
                                            </td>
                                        )}
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
};
