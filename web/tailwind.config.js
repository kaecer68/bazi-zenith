/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 五行配色
        'wood': '#2E7D32',
        'fire': '#D32F2F',
        'earth': '#F9A825',
        'metal': '#C0C0C0',
        'water': '#1976D2',
        // 陰陽配色
        'yang': '#FF5722',
        'yin': '#607D8B',
        // 傳統配色
        'traditional': {
          'red': '#8B0000',
          'gold': '#DAA520',
          'ink': '#1a1a1a',
          'paper': '#F5F5DC',
          'jade': '#00A86B',
        }
      },
      fontFamily: {
        'chinese': ['"Noto Serif TC"', '"Microsoft YaHei"', 'serif'],
      }
    },
  },
  plugins: [],
}
