/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./src/static/views/**/*.{html,js}'],
    theme: {},
    plugins: [
        require("@tailwindcss/typography"),
        require("daisyui"),
    ],
    daisyui: {
        themes: ["light", "winter"],
    },
}