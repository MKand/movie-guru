/** @type {import('tailwindcss').Config} */
import typographyPlugin from "@tailwindcss/typography";

export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  safelist: [
    {
      pattern:
        /bg-(slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose|seaweed)-(50|100|200|300|400|500|600|700|800|900|950)/,
    },
  ],
  theme: {
    screens: {
      xs: "320px",
      sm: "640px",
      md: "768px",
      lg: "1024px",
      xl: "1280px",
      xxl: "1536px",
    },
    extend: {
        boxShadow: {
          '3xl': '10px 10px 10px 5px rgba(0, 0, 0, 0.75)',
        },
        backgroundImage: {
          'stars1': "url('/src/assets/stars1.jpeg')",
          'reel': "url('/src/assets/reel-2.jpeg')",

                  }
        ,
      colors: {
        start: "#050a0d",
        primary: "#244855", 
        pop: "#E64833",
        accent: "#2e5c6b",
        text: "#FBE9D0",
        secondary: "#90AEAD",
        negative: "#660000",
        gurusilver: "#C6D4D2"
      },
      fontFamily: {
        'guru-title': ['Montserrat', 'sans-serif'],

        sans: [
          "Roboto",
        ],
        serif: [
          //font-serif
          "Roboto",
        ],
        mono: [
          //font-mono
          "ui-monospace",
          "SFMono-Regular",
          "Menlo",
          "Monaco",
          "Consolas",
          '"Liberation Mono"',
          '"Courier New"',
          "monospace",
        ],
        display: [
          // font-display
          "Lora",
        ],
      },
    },
  },
  plugins: [
    require('tailwind-scrollbar'),
  ],
}