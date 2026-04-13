import { Outlet } from "react-router-dom";
import { ProfileDropdown } from "./ProfileDropdown";

export const Layout = () => {
    return (
        <div className="min-h-screen">
            <ProfileDropdown />
            <Outlet />
        </div>
    );
};
