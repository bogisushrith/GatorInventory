import { ReactNode, useEffect, useMemo, useRef, useState } from "react";
import { Link } from "react-router-dom";

interface DashboardShellProps {
    role: string;
    activeTab: "products" | "users";
    title: string;
    subtitle: string;
    onLogout: () => Promise<void> | void;
    children: ReactNode;
    contentVariant?: "stack" | "grid";
    showSidebar?: boolean;
}

export const DashboardShell = ({ role, activeTab, title, subtitle, onLogout, children, contentVariant = "stack", showSidebar = true }: DashboardShellProps) => {
    const normalizedRole = role.toLowerCase();
    const contentContainerClassName = contentVariant === "grid"
        ? "lg:col-span-3 min-w-0 grid grid-cols-1 lg:grid-cols-3 gap-6 items-start"
        : "lg:col-span-3 min-w-0 flex flex-col gap-6";
    const contentContainerWithoutSidebarClassName = contentVariant === "grid"
        ? "min-w-0 grid grid-cols-1 lg:grid-cols-4 gap-6 items-start"
        : "min-w-0 flex flex-col gap-6";
    const headerClassName = contentVariant === "grid"
        ? `${showSidebar ? "lg:col-span-3" : "lg:col-span-4"} flex justify-between items-center gap-3`
        : "flex justify-between items-center gap-3";
    const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false);
    const profileMenuRef = useRef<HTMLDivElement | null>(null);

    const username = useMemo(() => {
        return (localStorage.getItem("username") || "User").trim() || "User";
    }, []);

    const initials = useMemo(() => {
        return username.charAt(0).toUpperCase();
    }, [username]);

    const formattedRole = useMemo(() => {
        return normalizedRole === "admin" ? "Admin" : "User";
    }, [normalizedRole]);

    useEffect(() => {
        const handleOutsideClick = (event: MouseEvent) => {
            const target = event.target as Node;
            if (profileMenuRef.current && !profileMenuRef.current.contains(target)) {
                setIsMenuOpen(false);
            }
        };

        document.addEventListener("mousedown", handleOutsideClick);
        return () => {
            document.removeEventListener("mousedown", handleOutsideClick);
        };
    }, []);

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
                            </nav>
                        </aside>

                        <div className={contentContainerClassName}>
                            <div className={headerClassName}>
                                <div>
                                    <h1 className="text-3xl font-bold tracking-tight text-gray-900">{title}</h1>
                                    <p className="text-gray-500 text-sm mt-1">{subtitle}</p>
                                </div>
                                <div className="flex items-center gap-3">
                                    <Link
                                        to="/"
                                        className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-all duration-200"
                                    >
                                        🏠 Home
                                    </Link>

                                    <div className="relative" ref={profileMenuRef}>
                                        <button
                                            type="button"
                                            onClick={() => setIsMenuOpen((previous) => !previous)}
                                            className="w-10 h-10 rounded-full bg-gray-200 text-gray-700 font-semibold flex items-center justify-center cursor-pointer hover:bg-gray-300 transition-all duration-200"
                                            aria-haspopup="menu"
                                            aria-expanded={isMenuOpen}
                                            aria-label="Open profile menu"
                                        >
                                            {initials}
                                        </button>

                                        <div
                                            className={[
                                                "absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 p-2 transition-all duration-200 z-20",
                                                isMenuOpen ? "opacity-100 scale-100 pointer-events-auto" : "opacity-0 scale-95 pointer-events-none"
                                            ].join(" ")}
                                            role="menu"
                                        >
                                            <div className="px-3 py-2">
                                                <p className="text-sm font-semibold text-gray-900 truncate">{username}</p>
                                                <p className="text-xs text-gray-500 mt-1">{formattedRole}</p>
                                            </div>

                                            <div className="h-px bg-gray-100 my-1" />

                                            <button
                                                type="button"
                                                onClick={onLogout}
                                                className="w-full text-left rounded-lg px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-all duration-200"
                                                role="menuitem"
                                            >
                                                🚪 Logout
                                            </button>
                                        </div>
                                    </div>
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
                                <Link
                                    to="/"
                                    className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-all duration-200"
                                >
                                    🏠 Home
                                </Link>

                                <div className="relative" ref={profileMenuRef}>
                                    <button
                                        type="button"
                                        onClick={() => setIsMenuOpen((previous) => !previous)}
                                        className="w-10 h-10 rounded-full bg-gray-200 text-gray-700 font-semibold flex items-center justify-center cursor-pointer hover:bg-gray-300 transition-all duration-200"
                                        aria-haspopup="menu"
                                        aria-expanded={isMenuOpen}
                                        aria-label="Open profile menu"
                                    >
                                        {initials}
                                    </button>

                                    <div
                                        className={[
                                            "absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 p-2 transition-all duration-200 z-20",
                                            isMenuOpen ? "opacity-100 scale-100 pointer-events-auto" : "opacity-0 scale-95 pointer-events-none"
                                        ].join(" ")}
                                        role="menu"
                                    >
                                        <div className="px-3 py-2">
                                            <p className="text-sm font-semibold text-gray-900 truncate">{username}</p>
                                            <p className="text-xs text-gray-500 mt-1">{formattedRole}</p>
                                        </div>

                                        <div className="h-px bg-gray-100 my-1" />

                                        <button
                                            type="button"
                                            onClick={onLogout}
                                            className="w-full text-left rounded-lg px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-all duration-200"
                                            role="menuitem"
                                        >
                                            🚪 Logout
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        {children}
                    </div>
                )}
            </div>
        </div>
    );
};
