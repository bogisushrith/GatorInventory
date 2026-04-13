interface FilterPanelProps {
    searchInput: string;
    onSearchChange: (value: string) => void;
    categoryInput: string;
    onCategoryChange: (value: string) => void;
    minPriceInput: string;
    onMinPriceChange: (value: string) => void;
    maxPriceInput: string;
    onMaxPriceChange: (value: string) => void;
    categories: string[];
    onReset: () => void; 
    onApply: () => void;
    hasPendingChanges: boolean;
    loading: boolean;
}

export const FilterPanel = ({
    searchInput,
    onSearchChange,
    categoryInput,
    onCategoryChange,
    minPriceInput,
    onMinPriceChange,
    maxPriceInput,
    onMaxPriceChange,
    categories,
    onReset,
    onApply,
    hasPendingChanges,
    loading
}: FilterPanelProps) => {
    return (
        <aside className="w-full bg-white rounded-2xl shadow-md p-5 h-full transition-all duration-200">
            <h3 className="text-xl font-semibold text-gray-900 mb-4">Product Filters</h3>

            <div className="space-y-4">
                <div>
                    <label className="block text-sm text-gray-500 mb-2">Search</label>
                    <input
                        value={searchInput}
                        onChange={(event) => onSearchChange(event.target.value)}
                        placeholder="Search products..."
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 transition-all duration-200"
                        disabled={loading}
                    />
                </div>

                <div>
                    <label className="block text-sm text-gray-500 mb-2">Category</label>
                    <select
                        value={categoryInput}
                        onChange={(event) => onCategoryChange(event.target.value)}
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 transition-all duration-200 bg-white"
                        disabled={loading}
                    >
                        <option value="">All categories</option>
                        {categories.map((categoryOption) => (
                            <option key={categoryOption} value={categoryOption}>
                                {categoryOption}
                            </option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm text-gray-500 mb-2">Min price</label>
                    <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={minPriceInput}
                        onChange={(event) => onMinPriceChange(event.target.value)}
                        placeholder="0.00"
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 transition-all duration-200"
                        disabled={loading}
                    />
                </div>

                <div>
                    <label className="block text-sm text-gray-500 mb-2">Max price</label>
                    <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={maxPriceInput}
                        onChange={(event) => onMaxPriceChange(event.target.value)}
                        placeholder="9999.99"
                        className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 transition-all duration-200"
                        disabled={loading}
                    />
                </div>

                <button
                    onClick={onApply}
                    disabled={loading || !hasPendingChanges}
                    className="w-full rounded-lg bg-gradient-to-r from-purple-500 to-pink-500 text-white px-3 py-2 font-medium hover:opacity-90 transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    Apply Filters
                </button>

                <button
                    onClick={onReset}
                    disabled={loading}
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-gray-700 font-medium hover:bg-gray-100 transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    Reset filters
                </button>
            </div>
        </aside>
    );
};
