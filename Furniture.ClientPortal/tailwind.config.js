/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./Pages/**/*.{razor,html,cshtml}",
    "./Shared/**/*.{razor,html,cshtml}",
    "./Components/**/*.{razor,html,cshtml}",
    "./wwwroot/**/*.{html,js}"
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#fef2f2',
          100: '#fee2e2',
          200: '#fecaca',
          300: '#fca5a5',
          400: '#f87171',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
          800: '#991b1b',
          900: '#7f1d1d',
        },
        accent: {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#c0c0c0',
          600: '#808080',
          700: '#6b7280',
          800: '#4b5563',
          900: '#374151',
        }
      },
      fontFamily: {
        'playwrite': ['Playwrite US Trad', 'cursive'],
      },
      animation: {
        'gradient-shift': 'gradientShift 20s ease infinite',
        'shimmer': 'shimmer 3s ease-in-out infinite',
        'slide-up': 'slideUp 0.6s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'slide-down': 'slideDown 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'shake': 'shake 0.4s ease',
        'spin': 'spin 0.8s linear infinite',
      },
      keyframes: {
        gradientShift: {
          '0%, 100%': { transform: 'translate(0, 0)' },
          '50%': { transform: 'translate(30px, -30px)' },
        },
        shimmer: {
          '0%, 100%': { opacity: '0.3' },
          '50%': { opacity: '0.8' },
        },
        slideUp: {
          from: {
            opacity: '0',
            transform: 'translateY(30px)',
          },
          to: {
            opacity: '1',
            transform: 'translateY(0)',
          },
        },
        slideDown: {
          from: {
            opacity: '0',
            transform: 'translateY(-10px)',
          },
          to: {
            opacity: '1',
            transform: 'translateY(0)',
          },
        },
        shake: {
          '0%, 100%': { transform: 'translateX(0)' },
          '25%': { transform: 'translateX(-5px)' },
          '75%': { transform: 'translateX(5px)' },
        },
        spin: {
          to: { transform: 'rotate(360deg)' },
        },
      },
    },
  },
  plugins: [],
}
