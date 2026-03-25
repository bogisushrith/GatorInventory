import { useNavigate } from "react-router-dom";
import { FormEvent, useState } from "react";

export const Login = () => {
    const navigate = useNavigate();
    const [errorMessage, setErrorMessage] = useState("");

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const data = {
            username: formData.get("username"),
            password: formData.get("password")
        };

        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data),
                credentials: 'include'
            });

            if (response.ok) {
                const loginResult = await response.json();
                const role = String(loginResult.role || "user").toLowerCase();
                localStorage.setItem("role", role);
                localStorage.setItem("username", String(data.username || ""));
                navigate("/dashboard");
            } else {
                const errorData = await response.json();
                setErrorMessage(errorData.error_message || "Login failed");
            }
        } catch (error) {
            setErrorMessage("Error: " + (error as Error).message);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8">
            <div className="w-full max-w-md">
                {/* Header */}
                <div className="text-center mb-10">
                    <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-indigo-600 to-pink-500 rounded-xl shadow-lg mb-4">
                        <span className="text-white text-3xl">🔐</span>
                    </div>
                    <h2 className="text-4xl font-bold text-gray-900 gradient-text">Welcome Back</h2>
                    <p className="text-gray-600 mt-2">Sign in to your inventory account</p>
                </div>

                {/* Form Card */}
                <div className="bg-white rounded-2xl shadow-soft-lg p-8 space-y-6">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {/* Username Field */}
                        <div>
                            <label htmlFor="username" className="block text-sm font-semibold text-gray-900 mb-2">
                                👤 Username
                            </label>
                            <input
                                id="username"
                                name="username"
                                type="text"
                                required
                                placeholder="Enter your username"
                                className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"
                            />
                        </div>

                        {/* Password Field */}
                        <div>
                            <label htmlFor="password" className="block text-sm font-semibold text-gray-900 mb-2">
                                🔑 Password
                            </label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                required
                                placeholder="Enter your password"
                                className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"
                            />
                        </div>

                        {/* Submit Button */}
                        <button
                            type="submit"
                            className="w-full btn-primary justify-center py-3 text-base font-semibold"
                        >
                            ✨ Login
                        </button>
                    </form>

                    {/* Error Message */}
                    {errorMessage && (
                        <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
                            <p className="text-red-700 font-medium">❌ {errorMessage}</p>
                        </div>
                    )}

                    {/* Demo Credentials */}
                    <div className="bg-gradient-to-r from-indigo-50 to-pink-50 border-2 border-indigo-200 rounded-xl p-6 space-y-3">
                        <h3 className="text-sm font-bold text-gray-900">✨ Demo Credentials</h3>
                        <div className="space-y-2 text-sm">
                            <p className="text-gray-700">
                                <span className="font-semibold text-indigo-600">Username:</span>
                                <code className="bg-white px-2 py-1 rounded ml-2 text-gray-900 font-mono">erkin</code>
                            </p>
                            <p className="text-gray-700">
                                <span className="font-semibold text-indigo-600">Password:</span>
                                <code className="bg-white px-2 py-1 rounded ml-2 text-gray-900 font-mono">Test1234</code>
                            </p>
                        </div>
                    </div>

                    {/* Sign Up Link */}
                    <div className="text-center pt-4 border-t border-gray-200">
                        <p className="text-gray-600 text-sm">
                            Don't have an account?{' '}
                            <a href="/signup" className="font-semibold text-indigo-600 hover:text-indigo-700 transition-colors">
                                Sign up here
                            </a>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};
