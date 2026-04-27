import { FormEvent, useMemo, useState } from "react";
import Cookies from "js-cookie";
import { Link, useLocation, useNavigate } from "react-router-dom";

interface NavigationItem {
    label: string;
    to: string;
    match: (pathname: string) => boolean;
}

const getNavigationItems = (isAuthenticated: boolean, role: string): NavigationItem[] => {
    if (!isAuthenticated) {
        return [];
    }

    if (role === "admin") {
        return [
            {
                label: "Home",
                to: "/",
                match: (pathname) => pathname === "/",
            },
            {
                label: "Products",
                to: "/products",
                match: (pathname) => pathname === "/products" || pathname === "/dashboard" || pathname === "/add-product" || pathname === "/update-product",
            },
            {
                label: "Sales",
                to: "/orders",
                match: (pathname) => pathname === "/orders" || pathname === "/sales",
            },
            {
                label: "Analytics",
                to: "/analytics",
                match: (pathname) => pathname === "/analytics" || pathname === "/admin/analytics",
            },
            {
                label: "Users",
                to: "/users",
                match: (pathname) => pathname === "/users",
            },
        ];
    }

    return [
        {
            label: "Home",
            to: "/",
            match: (pathname) => pathname === "/",
        },
        {
            label: "Products",
            to: "/products",
            match: (pathname) => pathname === "/products" || pathname === "/dashboard",
        },
        {
            label: "Cart",
            to: "/cart",
            match: (pathname) => pathname === "/cart" || pathname === "/checkout",
        },
        {
            label: "Order History",
            to: "/history",
            match: (pathname) => pathname === "/history" || pathname === "/orders/history" || pathname === "/sales",
        },
        {
            label: "Insights",
            to: "/user/analytics",
            match: (pathname) => pathname === "/user/analytics",
        },
    ];
};

export const TopNavbar = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [searchValue, setSearchValue] = useState("");

    const username = (localStorage.getItem("username") || "User").trim() || "User";
    const role = (localStorage.getItem("role") || "user").toLowerCase();
    const isAuthenticated = Boolean(localStorage.getItem("token") || localStorage.getItem("role") || Cookies.get("token"));

    const visibleItems = useMemo(() => (
        getNavigationItems(isAuthenticated, role)
    ), [isAuthenticated, role]);

    const roleLabel = role === "admin" ? "Admin" : "User";

    const handleSearch = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const query = searchValue.trim();
        if (!query) {
            navigate("/products");
            return;
        }

        navigate(`/products?search=${encodeURIComponent(query)}`);
    };

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", { method: "POST", credentials: "include" });
        } finally {
            localStorage.removeItem("token");
            localStorage.removeItem("role");
            localStorage.removeItem("username");
            navigate("/");
        }
    };

    return (
        <header className="sticky top-0 z-50 w-full border-b border-blue-800 bg-blue-700 text-white shadow-md">
            <div className="flex w-full flex-wrap items-center gap-4 px-6 py-3 xl:flex-nowrap">
                <Link to={isAuthenticated ? "/" : "/"} className="flex shrink-0 items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-white text-lg font-bold text-blue-700">
                        IC
                    </div>
                    <h1 className="text-lg font-bold tracking-tight">Gator Inventory</h1>
                </Link>

                <form onSubmit={handleSearch} className="flex w-full flex-1 items-center rounded-full bg-white px-2 py-1 shadow-sm">
                    <input
                        value={searchValue}
                        onChange={(event) => setSearchValue(event.target.value)}
                        className="w-full rounded-full px-4 py-2 text-black outline-none placeholder:text-slate-400"
                        placeholder="Search..."
                    />
                    <button type="submit" className="rounded-full bg-blue-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-800">
                        Search
                    </button>
                </form>

                <div className="ml-auto flex flex-wrap items-center gap-5 text-sm font-medium">
                    {isAuthenticated ? (
                        <>
                            {visibleItems.map((item) => {
                                const isActive = item.match(location.pathname);

                                return (
                                    <Link
                                        key={item.label}
                                        to={item.to}
                                        className={isActive ? "text-white underline underline-offset-8" : "text-blue-100 transition hover:text-white"}
                                    >
                                        {item.label}
                                    </Link>
                                );
                            })}

                            <button type="button" onClick={handleLogout} className="text-blue-100 transition hover:text-white">
                                Logout
                            </button>

                            <div className="flex items-center gap-3 rounded-full bg-white/10 px-2 py-1">
                                <span className="rounded-full bg-white/20 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-white">
                                    {roleLabel}
                                </span>
                                <div
                                    className="flex h-8 w-8 items-center justify-center rounded-full bg-white text-sm font-semibold text-black"
                                    title={`${username} (${roleLabel})`}
                                >
                                    {username.charAt(0).toUpperCase()}
                                </div>
                            </div>
                        </>
                    ) : (
                        <>
                            <Link to="/login" className="text-blue-100 transition hover:text-white">Login</Link>
                            <Link to="/signup" className="rounded-full bg-white px-4 py-2 text-sm font-semibold text-blue-700 transition hover:bg-blue-50">
                                Sign Up
                            </Link>
                        </>
                    )}
                </div>
            </div>
        </header>
    );
};
