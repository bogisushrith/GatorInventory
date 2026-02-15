import { useEffect, useState } from "react";
import Cookies from "js-cookie";
import { Link } from "react-router-dom";

export const Root = () => {
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);

    useEffect(() => {
        const token = Cookies.get("token");
        if (token) {
            setIsLoggedIn(() => true);
        } else {
            setIsLoggedIn(() => false);
        }
    }, []);

    return (
        <div className="w-full h-[90vh] flex flex-col justify-between">
            <div className="flex flex-row justify-between items-center mb-8">
                <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-gradient-to-br from-indigo-600 to-pink-500 rounded-lg flex items-center justify-center shadow-lg">
                        <span className="text-white font-bold text-xl">📦</span>
                    </div>
                    <h1 className="text-4xl font-bold gradient-text">Inventory System</h1>
                </div>
                {isLoggedIn
                    ? <Link to="dashboard"
                            className="btn-primary">
                        📊 Dashboard
                    </Link>
                    : <div className="flex items-center gap-3">
                        <Link to="signup"
                              className="btn-outline">
                            ✨ Sign Up
                        </Link>
                        <Link to="login"
                              className="btn-primary">
                            🔐 Login
                        </Link>
                    </div>
                }
            </div>
            <div className="w-full flex-1 flex items-center justify-center">
                <div className="flex flex-col items-center justify-center max-h-full gap-8">
                    <div className="text-6xl mb-4">📈</div>
                    <h2 className="text-5xl font-bold text-gray-900">A Modern Inventory System</h2>
                    <p className="text-xl text-gray-600 max-w-2xl leading-relaxed">
                        An intelligent inventory management platform designed to optimize tracking, 
                        organization, and control of warehouse products. Manage quantities, prices, 
                        and product data effortlessly.
                    </p>
                    {!isLoggedIn && (
                        <div className="mt-8 bg-gradient-to-r from-indigo-50 to-pink-50 border-2 border-indigo-200 rounded-xl p-8 max-w-md">
                            <h3 className="text-lg font-bold text-gray-900 mb-4">✨ Demo Credentials</h3>
                            <div className="space-y-3 text-left">
                                <p className="text-gray-700"><span className="font-semibold">Username:</span> <code className="bg-gray-200 px-2 py-1 rounded text-sm">erkin</code></p>
                                <p className="text-gray-700"><span className="font-semibold">Password:</span> <code className="bg-gray-200 px-2 py-1 rounded text-sm">Test1234</code></p>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};
