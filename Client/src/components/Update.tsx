import { FormEvent, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

interface Product {
    id: number;
    name: string;
    price: number;
    quantity: number;
    category: string;
}

export const Update = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const product = location.state?.product as Product | undefined;

    const [name, setName] = useState("");
    const [price, setPrice] = useState<number>(0);
    const [quantity, setQuantity] = useState<number>(0);
    const [category, setCategory] = useState("");
    const [errorMessage, setErrorMessage] = useState("");
    const [successMessage, setSuccessMessage] = useState("");

    useEffect(() => {
        if (!product) {
            navigate("/dashboard");
            return;
        }

        setName(product.name);
        setPrice(product.price);
        setQuantity(product.quantity);
        setCategory(product.category);
    }, [product, navigate]);

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (!product) {
            setErrorMessage("No product selected for update");
            return;
        }

        const data = {
            name,
            price: Number(price),
            quantity: Number(quantity),
            category
        };

        try {
            const response = await fetch(`/api/products/${product.id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data),
                credentials: "include"
            });

            if (response.ok) {
                setSuccessMessage("Product updated successfully! Redirecting to dashboard...");
                setErrorMessage("");
                setTimeout(() => navigate("/dashboard"), 1500);
            } else {
                const errorData = await response.json();
                setErrorMessage(errorData.error_message || "Failed to update product");
                setSuccessMessage("");
            }
        } catch (error) {
            setErrorMessage("Error: " + (error as Error).message);
            setSuccessMessage("");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-gray-50 via-indigo-50 to-pink-50">
            <div className="w-full max-w-xl">
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-indigo-600 to-pink-500 rounded-xl shadow-lg mb-4">
                        <span className="text-white text-2xl">✏️</span>
                    </div>
                    <h2 className="text-4xl font-bold gradient-text">Update Product</h2>
                    <p className="text-gray-600 mt-2">Edit product details and save changes</p>
                </div>

                <div className="bg-white rounded-2xl shadow-soft-lg p-8">
                    <form onSubmit={handleSubmit} className="space-y-5">
                        <div>
                            <label htmlFor="name" className="block text-sm font-semibold text-gray-900 mb-2">📦 Product Name</label>
                            <input id="name" name="name" type="text" required value={name} onChange={(event) => setName(event.target.value)}
                                   className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"/>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label htmlFor="price" className="block text-sm font-semibold text-gray-900 mb-2">💵 Price</label>
                                <input id="price" name="price" type="number" step="0.01" min="0" required value={price} onChange={(event) => setPrice(Number(event.target.value))}
                                       className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"/>
                            </div>
                            <div>
                                <label htmlFor="quantity" className="block text-sm font-semibold text-gray-900 mb-2">📈 Quantity</label>
                                <input id="quantity" name="quantity" type="number" min="0" required value={quantity} onChange={(event) => setQuantity(Number(event.target.value))}
                                       className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"/>
                            </div>
                        </div>

                        <div>
                            <label htmlFor="category" className="block text-sm font-semibold text-gray-900 mb-2">🏷️ Category</label>
                            <input id="category" name="category" type="text" required value={category} onChange={(event) => setCategory(event.target.value)}
                                   className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"/>
                        </div>

                        <div className="flex gap-3 pt-2">
                            <button type="button" onClick={() => navigate("/dashboard")} className="btn-outline w-full">↩ Back</button>
                            <button type="submit" className="btn-primary w-full">💾 Save Changes</button>
                        </div>
                    </form>

                    {errorMessage && <div className="mt-5 bg-red-50 border-l-4 border-red-500 p-3 rounded"><p className="text-red-700 font-medium">❌ {errorMessage}</p></div>}
                    {successMessage && <div className="mt-5 bg-green-50 border-l-4 border-green-500 p-3 rounded"><p className="text-green-700 font-medium">✅ {successMessage}</p></div>}
                </div>
            </div>
        </div>
    );
};
