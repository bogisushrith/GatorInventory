import { Outlet } from "react-router-dom";
import { TopNavbar } from "./TopNavbar";

export const Layout = () => {
    return (
        <div className="w-full min-h-screen bg-gray-50">
            <TopNavbar />
            <main className="w-full px-6 py-6">
                <Outlet />
            </main>
        </div>
    );
};
