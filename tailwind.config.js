/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./cmd/web/**/*.{templ,js}'],
  theme: {
    extend: {
      fontFamily: {
        geist: 'Geist',
        bebas: 'Bebas Neue',
      },
    },
  },
  plugins: [],
};
