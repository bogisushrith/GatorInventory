interface OrderFiltersProps {
    search: string;
    status: string;
    dateFrom: string;
    dateTo: string;
    searchPlaceholder: string;
    onSearchChange: (value: string) => void;
    onStatusChange: (value: string) => void;
    onDateFromChange: (value: string) => void;
    onDateToChange: (value: string) => void;
    onApply: () => void;
    onReset: () => void;
    loading: boolean;
}

export const OrderFilters = ({
    search,
    status,
    dateFrom,
    dateTo,
    searchPlaceholder,
    onSearchChange,
    onStatusChange,
    onDateFromChange,
    onDateToChange,
    onApply,
    onReset,
    loading,
}: OrderFiltersProps) => {
    return (
        <div className="bg-white rounded-2xl shadow-md p-5">
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-600 mb-2">Search</label>
                    <input
                        value={search}
                        onChange={(event) => onSearchChange(event.target.value)}
                        placeholder={searchPlaceholder}
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-400"
                        disabled={loading}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-600 mb-2">Status</label>
                    <select
                        value={status}
                        onChange={(event) => onStatusChange(event.target.value)}
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-400 bg-white"
                        disabled={loading}
                    >
                        <option value="">All statuses</option>
                        <option value="pending">Pending</option>
                        <option value="completed">Completed</option>
                        <option value="cancelled">Cancelled</option>
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-600 mb-2">Date from</label>
                    <input
                        type="date"
                        value={dateFrom}
                        onChange={(event) => onDateFromChange(event.target.value)}
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-400"
                        disabled={loading}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-600 mb-2">Date to</label>
                    <input
                        type="date"
                        value={dateTo}
                        onChange={(event) => onDateToChange(event.target.value)}
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-400"
                        disabled={loading}
                    />
                </div>
            </div>

            <div className="mt-4 flex flex-wrap items-center gap-3">
                <button
                    type="button"
                    onClick={onApply}
                    disabled={loading}
                    className="inline-flex items-center justify-center rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-700 transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    Apply Filters
                </button>
                <button
                    type="button"
                    onClick={onReset}
                    disabled={loading}
                    className="inline-flex items-center justify-center rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    Reset Filters
                </button>
            </div>
        </div>
    );
};
