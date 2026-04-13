import { useEffect, useMemo, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";

export const ProfileDropdown = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const containerRef = useRef<HTMLDivElement | null>(null);

    const username = useMemo(() => (localStorage.getItem("username") || "User").trim() || "User", []);
    const role = useMemo(() => (localStorage.getItem("role") || "guest").toLowerCase(), []);
    const initial = username.charAt(0).toUpperCase();

    const isLoggedIn = !!localStorage.getItem("token") || !!Cookies.get("token") || !!localStorage.getItem("role");
    const shouldHideOnRoute = location.pathname === "/login" || location.pathname === "/signup";

    useEffect(() => {
        const onClickOutside = (event: MouseEvent) => {
            const target = event.target as Node;
            if (containerRef.current && !containerRef.current.contains(target)) {
                setIsOpen(false);
            }
        };

        document.addEventListener("mousedown", onClickOutside);
        return () => document.removeEventListener("mousedown", onClickOutside);
    }, []);

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", { method: "POST", credentials: "include" });
        } finally {
            localStorage.removeItem("token");
            localStorage.removeItem("role");
            localStorage.removeItem("username");
            setIsOpen(false);
            navigate("/login");
        }
    };

    if (!isLoggedIn || shouldHideOnRoute) {
        return null;
    }

    return (
        <div className="fixed top-4 right-5 z-50" ref={containerRef}>
            <button
                type="button"
                onClick={() => setIsOpen((previous) => !previous)}
                className="w-11 h-11 rounded-full bg-gray-900 text-white text-sm font-semibold flex items-center justify-center shadow-lg hover:bg-gray-800 transition-all"
                title={`${username} (${role})`}
                aria-label="Open profile menu"
                aria-expanded={isOpen}
            >
                {initial}
            </button>

            {isOpen && (
                <div className="absolute right-0 mt-2 w-52 rounded-xl border border-gray-100 bg-white p-2 shadow-xl">
                    <div className="px-3 py-2">
                        <p className="truncate text-sm font-semibold text-gray-900">{username}</p>
                        <p className="mt-1 text-xs capitalize text-gray-500">{role}</p>
                    </div>

                    <div className="my-1 h-px bg-gray-100" />

                    <button
                        type="button"
                        onClick={handleLogout}
                        className="w-full rounded-lg px-3 py-2 text-left text-sm font-medium text-gray-700 transition-all duration-200 hover:bg-gray-50"
                    >
                        Logout
                    </button>
                </div>
            )}
        </div>
    );
};
