/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./{templates,generated}/**/*.html"],
    darkMode: "class",
    theme: {
        extend: {
            fontFamily: {
                "gothic": ["Gothic", "sans-serif"]
            },
            dark: {
                "bg-primary": "bg-cyan-950",
                "text-primary": "text-white",
            },
        },
    },
    plugins: [],
};
