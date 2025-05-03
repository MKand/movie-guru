/**
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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