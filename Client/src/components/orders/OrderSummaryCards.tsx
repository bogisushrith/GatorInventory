interface OrderSummaryCardsProps {
    totalOrders: number;
    totalRevenue: number;
}

export const OrderSummaryCards = ({ totalOrders, totalRevenue }: OrderSummaryCardsProps) => {
    return (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="rounded-2xl bg-white shadow-md p-5">
                <p className="text-sm text-gray-500">Total Orders</p>
                <p className="mt-2 text-3xl font-bold text-gray-900">{totalOrders}</p>
            </div>
            <div className="rounded-2xl bg-white shadow-md p-5">
                <p className="text-sm text-gray-500">Total Revenue</p>
                <p className="mt-2 text-3xl font-bold text-gray-900">${totalRevenue.toFixed(2)}</p>
            </div>
        </div>
    );
};
