import { ChangeEvent, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

interface Product {
    id: number;
    name: string;
    price: number;
    quantity: number;
    category: string;
}

export const Dashboard = () => {
    const navigate = useNavigate();
    const [products, setProducts] = useState<Product[]>([]);
    const [allProducts, setAllProducts] = useState<Product[]>([]);
    const [categories, setCategories] = useState<string[]>([]);
    const [selectedCategory, setSelectedCategory] = useState<string>("");

    const handleLogout = async () => {
        try {
            const response = await fetch("/api/logout", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                credentials: 'include'
            });

            if (response.ok) {
                navigate("/");
            } else {
                console.error("Logout failed");
            }
        } catch (error) {
            console.error("Error:", error);
        }
    };

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                const response = await fetch("/api/products", {
                    credentials: 'include'
                });
                if (response.ok) {
                    const data = await response.json() || [];
                    setProducts(data);
                    setAllProducts(data);
                } else {
                    console.error("Failed to fetch products:", response.status);
                }
            } catch (error) {
                console.error("Error fetching products:", error);
            }
        };

        fetchProducts();
    }, []);

    useEffect(() => {
        setCategories(() => Array.from(new Set(allProducts.map(p => p.category))));
    }, [allProducts]);

    const handleCategoryChange = (event: ChangeEvent<HTMLSelectElement>) => {
        const category = event.target.value;
        setSelectedCategory(category);
        if (category) {
            setProducts(allProducts.filter(product => product.category === category));
        } else {
            setProducts(allProducts);
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-50 via-indigo-50 to-pink-50">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* Header */}
                <div className="mb-10">
                    <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                        <div className="flex items-center gap-3">
                            <div className="w-14 h-14 bg-gradient-to-br from-indigo-600 to-pink-500 rounded-xl flex items-center justify-center shadow-lg">
                                <span className="text-white text-2xl">📊</span>
                            </div>
                            <div>
                                <h1 className="text-4xl font-bold gradient-text">Inventory Dashboard</h1>
                                <p className="text-gray-600 text-sm mt-1">View all products</p>
                            </div>
                        </div>
                        <button
                            onClick={handleLogout}
                            className="btn-primary flex items-center gap-2 whitespace-nowrap"
                        >
                            🚪 Logout
                        </button>
                    </div>
                </div>

                {/* Filter Section */}
                <div className="mb-8">
                    <label className="block text-sm font-semibold text-gray-900 mb-3">
                        🏷️ Filter by Category
                    </label>
                    <select
                        value={selectedCategory}
                        onChange={handleCategoryChange}
                        className="w-full md:w-64 rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 font-medium focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all bg-white"
                    >
                        <option value="">All Categories ({allProducts.length})</option>
                        {categories.map((category, index) => (
                            <option key={index} value={category}>
                                {category} ({allProducts.filter(p => p.category === category).length})
                            </option>
                        ))}
                    </select>
                </div>

                {/* Products Table */}
                <div className="bg-white rounded-2xl shadow-soft-lg overflow-hidden card-shadow">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead>
                                <tr className="bg-gradient-to-r from-indigo-600 to-pink-500 text-white">
                                    <th className="px-6 py-4 text-left text-sm font-bold uppercase tracking-wider">📦 Product Name</th>
                                    <th className="px-6 py-4 text-left text-sm font-bold uppercase tracking-wider">💵 Price</th>
                                    <th className="px-6 py-4 text-left text-sm font-bold uppercase tracking-wider">📈 Quantity</th>
                                    <th className="px-6 py-4 text-left text-sm font-bold uppercase tracking-wider">💰 Total Value</th>
                                    <th className="px-6 py-4 text-left text-sm font-bold uppercase tracking-wider">🏷️ Category</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {products.length > 0 ? (
                                    products.map((product, index) => (
                                        <tr key={index} className="hover:bg-gradient-to-r hover:from-indigo-50 hover:to-pink-50 transition-colors duration-200">
                                            <td className="px-6 py-4 text-sm font-semibold text-gray-900">{product.name}</td>
                                            <td className="px-6 py-4 text-sm text-gray-700">
                                                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
                                                    {new Intl.NumberFormat('en-US', {
                                                        style: 'currency',
                                                        currency: 'USD'
                                                    }).format(product.price)}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 text-sm font-semibold text-gray-900">
                                                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-pink-100 text-pink-800">
                                                    {product.quantity} units
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 text-sm font-bold">
                                                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-bold bg-gradient-to-r from-indigo-100 to-pink-100 text-indigo-900">
                                                    {new Intl.NumberFormat('en-US', {
                                                        style: 'currency',
                                                        currency: 'USD'
                                                    }).format(parseFloat((product.quantity * product.price).toFixed(2)))}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 text-sm text-gray-700">
                                                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-cyan-100 text-cyan-800">
                                                    {product.category}
                                                </span>
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={5} className="px-6 py-12 text-center">
                                            <p className="text-gray-500 text-lg font-medium">No products found</p>
                                            <p className="text-gray-400 text-sm mt-1">Try selecting a different category</p>
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Summary Stats */}
                {products.length > 0 && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-10">
                        <div className="bg-white rounded-xl shadow-soft p-6 border-l-4 border-indigo-600">
                            <p className="text-gray-600 text-sm font-medium mb-2">📦 Total Products</p>
                            <p className="text-3xl font-bold text-gray-900">{products.length}</p>
                        </div>
                        <div className="bg-white rounded-xl shadow-soft p-6 border-l-4 border-pink-500">
                            <p className="text-gray-600 text-sm font-medium mb-2">📊 Total Quantity</p>
                            <p className="text-3xl font-bold text-gray-900">{products.reduce((sum, p) => sum + p.quantity, 0)}</p>
                        </div>
                        <div className="bg-white rounded-xl shadow-soft p-6 border-l-4 border-cyan-500">
                            <p className="text-gray-600 text-sm font-medium mb-2">💼 Total Inventory Value</p>
                            <p className="text-3xl font-bold text-gray-900">
                                {new Intl.NumberFormat('en-US', {
                                    style: 'currency',
                                    currency: 'USD',
                                    maximumFractionDigits: 0
                                }).format(products.reduce((sum, p) => sum + (p.quantity * p.price), 0))}
                            </p>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
