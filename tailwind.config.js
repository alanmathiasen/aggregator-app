/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/views/**/*.{html,js,templ}"],
  theme: { 
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: true
  }
};
