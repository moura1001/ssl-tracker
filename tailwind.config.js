const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./src/static/views/**/*.{html,js}'],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter', 'system-ui', 'sans-serif'],
            },
            colors: {
                orange: colors.orange,
                green: colors.green,
            },
        },
    },
    plugins: [],
}