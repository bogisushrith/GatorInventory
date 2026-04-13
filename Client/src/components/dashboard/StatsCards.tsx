import { Product } from "./types";

interface StatsCardsProps {
    products: Product[];
    totalResults: number;
}

export const StatsCards = ({ products, totalResults }: StatsCardsProps) => {
    const visibleValue = products.reduce((sum, product) => sum + (product.quantity * product.price), 0);
    const visibleQuantity = products.reduce((sum, product) => sum + product.quantity, 0);

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-white rounded-2xl shadow-md p-5 transition-all duration-200 hover:-translate-y-0.5">
                <p className="text-sm text-gray-500 mb-2">📊 Total Matching Results</p> 
                <p className="text-lg font-semibold text-gray-900">{totalResults}</p>
            </div>
            <div className="bg-white rounded-2xl shadow-md p-5 transition-all duration-200 hover:-translate-y-0.5">
                <p className="text-sm text-gray-500 mb-2">📈 Visible Quantity</p>
                <p className="text-lg font-semibold text-gray-900">{visibleQuantity}</p>
            </div>
            <div className="bg-white rounded-2xl shadow-md p-5 transition-all duration-200 hover:-translate-y-0.5">
                <p className="text-sm text-gray-500 mb-2">💼 Visible Inventory Value</p>
                <p className="text-lg font-semibold text-gray-900">
                    {new Intl.NumberFormat("en-US", {
                        style: "currency",
                        currency: "USD",
                        maximumFractionDigits: 0
                    }).format(visibleValue)}
                </p>
            </div>
        </div>
    );
};
