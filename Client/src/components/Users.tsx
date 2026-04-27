import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { DashboardShell } from "./dashboard/DashboardShell.tsx";

interface UserItem {
    id: number;
    username: string;
    email: string;
    role: "admin" | "user";
}

export const Users = () => {
    const navigate = useNavigate();
    const role = useMemo(() => (localStorage.getItem("role") || "user").toLowerCase(), []);

    const [users, setUsers] = useState<UserItem[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>("");
    const [successMessage, setSuccessMessage] = useState<string>("");

    const fetchUsers = async () => {
        setLoading(true);
        setError("");
        try {
            const response = await fetch("/api/users", { credentials: "include" });
            if (!response.ok) {
                const errorData = await response.json();
                setError(errorData.error_message || "Failed to load users");
                return;
            }

            const data = (await response.json()) as UserItem[];
            setUsers(data.map((user) => ({ ...user, role: user.role.toLowerCase() as "admin" | "user" })));
        } catch (err) {
            setError((err as Error).message);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (role !== "admin") {
            navigate("/dashboard", { replace: true });
            return;
        }
        fetchUsers();
    }, [role]);

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", { method: "POST", credentials: "include" });
        } finally {
            localStorage.removeItem("role");
            localStorage.removeItem("username");
            navigate("/");
        }
    };

    const handleRoleChange = async (id: number, nextRole: "admin" | "user") => {
        setSuccessMessage("");
        try {
            const response = await fetch(`/api/users/${id}/role`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ role: nextRole })
            });

            if (!response.ok) {
                const errorData = await response.json();
                setError(errorData.error_message || "Failed to update role");
                return;
            }

            setUsers((previous) => previous.map((user) => user.id === id ? { ...user, role: nextRole } : user));
            setSuccessMessage("Role updated successfully");
        } catch (err) {
            setError((err as Error).message);
        }
    };

    return (
        <DashboardShell
            role={role}
            activeTab="users"
            title="User Management"
            subtitle="Manage user roles and permissions"
            onLogout={handleLogout}
            contentVariant="stack"
            showSidebar={false}
            fullWidth={true}
        >
            {error && <div className="bg-red-50 border-l-4 border-red-500 p-5 rounded-lg text-red-700">❌ {error}</div>}
            {successMessage && <div className="bg-green-50 border-l-4 border-green-500 p-5 rounded-lg text-green-700">✅ {successMessage}</div>}

            <div className="w-full min-w-0 h-full">
                <div className="bg-white rounded-2xl shadow-md overflow-hidden h-full">
                    <div className="p-5 border-b border-gray-100">
                        <div className="flex justify-between items-center gap-4">
                            <h3 className="text-xl font-semibold text-gray-900">Users</h3>
                            <p className="text-sm text-gray-500">Showing {users.length} users</p>
                        </div>
                    </div>

                    <div className="overflow-x-auto">
                        <table className="w-full min-w-[680px]">
                            <thead>
                                <tr className="bg-gradient-to-r from-indigo-600 to-pink-500 text-white">
                                    <th className="px-5 py-3 text-center text-xs font-bold uppercase tracking-wide">Username</th>
                                    <th className="px-5 py-3 text-center text-xs font-bold uppercase tracking-wide">Email</th>
                                    <th className="px-5 py-3 text-center text-xs font-bold uppercase tracking-wide">Role</th>
                                    <th className="px-5 py-3 text-center text-xs font-bold uppercase tracking-wide">Change Role</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100">
                                {loading ? (
                                    <tr>
                                        <td colSpan={4} className="px-6 py-10 text-center text-gray-500">Loading users...</td>
                                    </tr>
                                ) : users.length === 0 ? (
                                    <tr>
                                        <td colSpan={4} className="px-6 py-10 text-center text-gray-500">No users found</td>
                                    </tr>
                                ) : (
                                    users.map((user) => (
                                        <tr key={user.id} className="hover:bg-gradient-to-r hover:from-indigo-50 hover:to-pink-50 transition-colors">
                                            <td className="px-5 py-4 text-sm font-medium text-gray-900 text-center">{user.username}</td>
                                            <td className="px-5 py-4 text-sm font-medium text-gray-900 text-center">{user.email?.trim() ? user.email : "N/A"}</td>
                                            <td className="px-5 py-4 text-sm text-center">
                                                <span className={[
                                                    "inline-flex px-3 py-1 rounded-full text-xs font-semibold",
                                                    user.role === "admin" ? "bg-indigo-100 text-indigo-700" : "bg-cyan-100 text-cyan-700"
                                                ].join(" ")}>
                                                    {user.role}
                                                </span>
                                            </td>
                                            <td className="px-5 py-4 text-sm text-center">
                                                <select
                                                    value={user.role}
                                                    onChange={(event) => handleRoleChange(user.id, event.target.value as "admin" | "user")}
                                                    className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100"
                                                >
                                                    <option value="user">user</option>
                                                    <option value="admin">admin</option>
                                                </select>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </DashboardShell>
    );
};
