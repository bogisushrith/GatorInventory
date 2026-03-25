import { Navigate } from "react-router-dom";
import { ReactNode } from "react";

interface RequireAdminProps {
    children: ReactNode;
}

export const RequireAdmin = ({ children }: RequireAdminProps) => {
    const role = (localStorage.getItem("role") || "").toLowerCase();

    if (role !== "admin") {
        return <Navigate to="/dashboard" replace />;
    }

    return <>{children}</>;
};
