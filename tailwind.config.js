/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.{html,js,templ}"],
  theme: { 
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: true
  }
};
