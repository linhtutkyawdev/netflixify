/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./cmd/web/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        poppins: 'Poppins',
        geist: 'Geist'
      }
    },
  },
  plugins: [],
}
