/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: {
                    50: '#eef2ff',
                    100: '#e0e7ff',
                    500: '#6366f1',
                    600: '#4f46e5',
                    700: '#4338ca',
                    800: '#3730a3',
                    900: '#312e81',
                },
                secondary: {
                    400: '#f472b6',
                    500: '#ec4899',
                    600: '#db2777',
                },
                accent: {
                    400: '#22d3ee',
                    500: '#06b6d4',
                    600: '#0891b2',
                },
            },
            fontFamily: {
                sans: ['Roboto', 'sans-serif'],
            },
            boxShadow: {
                'soft': '0 4px 15px rgba(99, 102, 241, 0.1)',
                'soft-lg': '0 8px 25px rgba(99, 102, 241, 0.15)',
            },
        },
    },
    plugins: [],
};

