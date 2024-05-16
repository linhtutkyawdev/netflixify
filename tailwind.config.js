/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./cmd/web/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        geist: 'Geist'
      }
    },
  },
  plugins: [],
}
