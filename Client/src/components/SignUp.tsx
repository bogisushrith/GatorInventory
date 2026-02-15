import { useNavigate } from "react-router-dom";
import { FormEvent, useState } from "react";

export const SignUp = () => {
    const navigate = useNavigate();
    const [errorMessage, setErrorMessage] = useState("");
    const [successMessage, setSuccessMessage] = useState("");

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const password = formData.get("password");
        const confirmPassword = formData.get("confirmPassword");

        if (password !== confirmPassword) {
            setErrorMessage("Passwords do not match");
            return;
        }

        const data = {
            username: formData.get("username"),
            email: formData.get("email"),
            password: password
        };

        try {
            const response = await fetch("/api/signup", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data),
                credentials: 'include'
            });

            if (response.ok) {
                setSuccessMessage("Account created successfully! Redirecting to login...");
                setErrorMessage("");
                setTimeout(() => {
                    navigate("/login");
                }, 2000);
            } else {
                const errorData = await response.json();
                setErrorMessage(errorData.error_message || "Sign up failed");
                setSuccessMessage("");
            }
        } catch (error) {
            setErrorMessage("Error: " + (error as Error).message);
            setSuccessMessage("");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8">
            <div className="w-full max-w-md">
                {/* Header */}
                <div className="text-center mb-10">
                    <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-indigo-600 to-pink-500 rounded-xl shadow-lg mb-4">
                        <span className="text-white text-3xl">✨</span>
                    </div>
                    <h2 className="text-4xl font-bold text-gray-900 gradient-text">Create Account</h2>
                    <p className="text-gray-600 mt-2">Join our inventory management system</p>
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
                                placeholder="Choose a username"
                                className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"
                            />
                        </div>

                        {/* Email Field */}
                        <div>
                            <label htmlFor="email" className="block text-sm font-semibold text-gray-900 mb-2">
                                📧 Email Address
                            </label>
                            <input
                                id="email"
                                name="email"
                                type="email"
                                required
                                placeholder="Enter your email"
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
                                placeholder="Create a password"
                                className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"
                            />
                        </div>

                        {/* Confirm Password Field */}
                        <div>
                            <label htmlFor="confirmPassword" className="block text-sm font-semibold text-gray-900 mb-2">
                                🔒 Confirm Password
                            </label>
                            <input
                                id="confirmPassword"
                                name="confirmPassword"
                                type="password"
                                required
                                placeholder="Confirm your password"
                                className="w-full rounded-lg border-2 border-gray-200 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-indigo-600 focus:outline-none focus:ring-2 focus:ring-indigo-100 transition-all"
                            />
                        </div>

                        {/* Submit Button */}
                        <button
                            type="submit"
                            className="w-full btn-primary justify-center py-3 text-base font-semibold"
                        >
                            🚀 Create Account
                        </button>
                    </form>

                    {/* Error Message */}
                    {errorMessage && (
                        <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
                            <p className="text-red-700 font-medium">❌ {errorMessage}</p>
                        </div>
                    )}

                    {/* Success Message */}
                    {successMessage && (
                        <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
                            <p className="text-green-700 font-medium">✅ {successMessage}</p>
                        </div>
                    )}

                    {/* Login Link */}
                    <div className="text-center pt-4 border-t border-gray-200">
                        <p className="text-gray-600 text-sm">
                            Already have an account?{' '}
                            <a href="/login" className="font-semibold text-indigo-600 hover:text-indigo-700 transition-colors">
                                Login here
                            </a>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};
