interface PaginationProps {
    page: number;
    totalPages: number;
    loading: boolean;
    onPageChange: (nextPage: number) => void; 
}

export const Pagination = ({ page, totalPages, loading, onPageChange }: PaginationProps) => {
    if (totalPages <= 1) {
        return null;
    }

    const visiblePages: number[] = [];
    const start = Math.max(1, page - 2);
    const end = Math.min(totalPages, page + 2);

    for (let pageNumber = start; pageNumber <= end; pageNumber++) {
        visiblePages.push(pageNumber);
    }

    return (
        <div className="flex justify-center items-center gap-3 mt-4 flex-wrap">
            <button
                onClick={() => onPageChange(page - 1)}
                disabled={loading || page <= 1}
                className="rounded-lg px-3 py-1 border border-gray-300 text-gray-700 hover:bg-gray-100 transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
            >
                Previous
            </button>

            <div className="flex items-center gap-2">
                {visiblePages.map((pageNumber) => (
                    <button
                        key={pageNumber}
                        onClick={() => onPageChange(pageNumber)}
                        disabled={loading}
                        className={[
                            "rounded-lg px-3 py-1 text-sm font-semibold transition-all duration-200",
                            pageNumber === page
                                ? "bg-purple-600 text-white"
                                : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-100",
                            loading ? "opacity-60 cursor-not-allowed" : ""
                        ].join(" ")}
                    >
                        {pageNumber}
                    </button>
                ))}
            </div>

            <button
                onClick={() => onPageChange(page + 1)}
                disabled={loading || page >= totalPages}
                className="rounded-lg px-3 py-1 border border-gray-300 text-gray-700 hover:bg-gray-100 transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
            >
                Next
            </button>
        </div>
    );
};
