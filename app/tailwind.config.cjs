/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,svelte}"],
  theme: {
    extend: {
      colors: {
        primary: "#074884",
        background: "#96eac2",
        button: "#ffc40c",
        shadow: "#161511",
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
