import { ReactNode } from "react";
import { Link } from "react-router-dom";

interface DashboardShellProps {
    role: string;
    activeTab: "products" | "cart" | "checkout" | "history" | "users" | "orders";
    title: string;
    subtitle: string;
    onLogout: () => Promise<void> | void;
    children: ReactNode;
    headerActions?: ReactNode;
    contentVariant?: "stack" | "grid";
    showSidebar?: boolean;
}

export const DashboardShell = ({ role, activeTab, title, subtitle, children, headerActions, contentVariant = "stack", showSidebar = true }: DashboardShellProps) => {
    const isAdmin = role.toLowerCase() === "admin";
    const contentContainerClassName = contentVariant === "grid"
        ? "lg:col-span-3 min-w-0 grid grid-cols-1 lg:grid-cols-3 gap-6 items-start"
        : "lg:col-span-3 min-w-0 flex flex-col gap-6";
    const contentContainerWithoutSidebarClassName = contentVariant === "grid"
        ? "min-w-0 grid grid-cols-1 lg:grid-cols-4 gap-6 items-start"
        : "min-w-0 flex flex-col gap-6";
    const headerClassName = contentVariant === "grid"
        ? `${showSidebar ? "lg:col-span-3" : "lg:col-span-4"} flex justify-between items-center gap-3`
        : "flex justify-between items-center gap-3";

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-50 via-indigo-50 to-pink-50">
            <div className="max-w-7xl mx-auto px-6 py-6">
                {showSidebar ? (
                    <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 items-start">
                        <aside className="lg:col-span-1 h-fit bg-white rounded-2xl shadow-md p-5 flex flex-col gap-4 transition-all duration-200">
                            <nav className="flex flex-col gap-4">
                                <Link
                                    to="/dashboard"
                                    className={[
                                        "block rounded-xl px-4 py-2.5 text-sm font-medium transition-all duration-200",
                                        activeTab === "products"
                                            ? "bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-md"
                                            : "text-gray-700 hover:bg-indigo-50"
                                    ].join(" ")}
                                >
                                    Products
                                </Link>
                                {!isAdmin && (
                                    <>
                                        <Link
                                            to="/cart"
                                            className={[
                                                "block rounded-xl px-4 py-2.5 text-sm font-medium transition-all duration-200",
                                                activeTab === "cart"
                                                    ? "bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md"
                                                    : "text-gray-700 hover:bg-indigo-50"
                                            ].join(" ")}
                                        >
                                            Cart
                                        </Link>
                                        <Link
                                            to="/history"
                                            className={[
                                                "block rounded-xl px-4 py-2.5 text-sm font-medium transition-all duration-200",
                                                activeTab === "history"
                                                    ? "bg-gradient-to-r from-sky-500 to-cyan-500 text-white shadow-md"
                                                    : "text-gray-700 hover:bg-indigo-50"
                                            ].join(" ")}
                                        >
                                            History
                                        </Link>
                                    </>
                                )}
                                {isAdmin && (
                                    <>
                                        <Link
                                            to="/orders"
                                            className={[
                                                "block rounded-xl px-4 py-2.5 text-sm font-medium transition-all duration-200",
                                                activeTab === "orders"
                                                    ? "bg-gradient-to-r from-sky-500 to-cyan-500 text-white shadow-md"
                                                    : "text-gray-700 hover:bg-indigo-50"
                                            ].join(" ")}
                                        >
                                            Orders
                                        </Link>
                                        <Link
                                            to="/users"
                                            className={[
                                                "block rounded-xl px-4 py-2.5 text-sm font-medium transition-all duration-200",
                                                activeTab === "users"
                                                    ? "bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md"
                                                    : "text-gray-700 hover:bg-indigo-50"
                                            ].join(" ")}
                                        >
                                            Users
                                        </Link>
                                    </>
                                )}
                            </nav>
                        </aside>

                        <div className={contentContainerClassName}>
                            <div className={headerClassName}>
                                <div>
                                    <h1 className="text-3xl font-bold tracking-tight text-gray-900">{title}</h1>
                                    <p className="text-gray-500 text-sm mt-1">{subtitle}</p>
                                </div>
                                <div className="flex items-center gap-3">
                                    {headerActions}
                                    <Link
                                        to="/"
                                        className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-all duration-200"
                                    >
                                        🏠 Home
                                    </Link>
                                </div>
                            </div>

                            {children}
                        </div>
                    </div>
                ) : (
                    <div className={contentContainerWithoutSidebarClassName}>
                        <div className={headerClassName}>
                            <div>
                                <h1 className="text-3xl font-bold tracking-tight text-gray-900">{title}</h1>
                                <p className="text-gray-500 text-sm mt-1">{subtitle}</p>
                            </div>
                            <div className="flex items-center gap-3">
                                {headerActions}
                                <Link
                                    to="/"
                                    className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-all duration-200"
                                >
                                    🏠 Home
                                </Link>
                            </div>
                        </div>

                        {children}
                    </div>
                )}
            </div>
        </div>
    );
};
