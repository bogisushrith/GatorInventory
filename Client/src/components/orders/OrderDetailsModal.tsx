import { Order } from "../dashboard/types";
import { calculateOrderTotal, formatOrderDate, formatStatusLabel, getStatusStyles } from "./orderUtils";

interface OrderDetailsModalProps {
    order: Order | null;
    onClose: () => void;
}

export const OrderDetailsModal = ({ order, onClose }: OrderDetailsModalProps) => {
    if (!order) {
        return null;
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-8">
            <div className="w-full max-w-3xl rounded-3xl bg-white shadow-2xl">
                <div className="flex items-center justify-between border-b border-gray-100 px-6 py-4">
                    <div>
                        <h3 className="text-xl font-bold text-gray-900">Order #{order.id}</h3>
                        <p className="text-sm text-gray-500">{order.user_name || "Unknown user"} • {formatOrderDate(order.created_at)}</p>
                    </div>
                    <button
                        type="button"
                        onClick={onClose}
                        className="rounded-lg border border-gray-300 px-3 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
                    >
                        Close
                    </button>
                </div>

                <div className="px-6 py-5 space-y-4">
                    <div className="flex flex-wrap items-center gap-3">
                        <span className={["inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold", getStatusStyles(order.status)].join(" ")}>{formatStatusLabel(order.status)}</span>
                        <span className="text-sm text-gray-600">Total: ${calculateOrderTotal(order).toFixed(2)}</span>
                    </div>

                    <div className="rounded-2xl border border-gray-100 overflow-hidden">
                        <div className="grid grid-cols-12 gap-3 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-500">
                            <div className="col-span-6">Product</div>
                            <div className="col-span-2 text-center">Qty</div>
                            <div className="col-span-2 text-right">Price</div>
                            <div className="col-span-2 text-right">Line Total</div>
                        </div>
                        <div className="divide-y divide-gray-100">
                            {order.items.map((item) => (
                                <div key={item.id} className="grid grid-cols-12 gap-3 px-4 py-3 text-sm items-center">
                                    <div className="col-span-6 font-medium text-gray-900">{item.product_name || `Product #${item.product_id}`}</div>
                                    <div className="col-span-2 text-center text-gray-700">{item.quantity}</div>
                                    <div className="col-span-2 text-right text-gray-700">${item.product_price.toFixed(2)}</div>
                                    <div className="col-span-2 text-right font-semibold text-gray-900">${(item.quantity * item.product_price).toFixed(2)}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};
