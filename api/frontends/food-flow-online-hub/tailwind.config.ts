
import type { Config } from "tailwindcss";
import tailwindcssAnimate from "tailwindcss-animate";

export default {
	darkMode: ["class"],
	content: [
		"./pages/**/*.{ts,tsx}",
		"./components/**/*.{ts,tsx}",
		"./app/**/*.{ts,tsx}",
		"./src/**/*.{ts,tsx}",
	],
	prefix: "",
	theme: {
		container: {
			center: true,
			padding: '2rem',
			screens: {
				'2xl': '1400px'
			}
		},
		extend: {
			colors: {
				border: 'hsl(var(--border))',
				input: 'hsl(var(--input))',
				ring: 'hsl(var(--ring))',
				background: 'hsl(var(--background))',
				foreground: 'hsl(var(--foreground))',
				primary: {
					DEFAULT: 'hsl(var(--primary))',
					foreground: 'hsl(var(--primary-foreground))'
				},
				secondary: {
					DEFAULT: 'hsl(var(--secondary))',
					foreground: 'hsl(var(--secondary-foreground))'
				},
				destructive: {
					DEFAULT: 'hsl(var(--destructive))',
					foreground: 'hsl(var(--destructive-foreground))'
				},
				muted: {
					DEFAULT: 'hsl(var(--muted))',
					foreground: 'hsl(var(--muted-foreground))'
				},
				accent: {
					DEFAULT: 'hsl(var(--accent))',
					foreground: 'hsl(var(--accent-foreground))'
				},
				popover: {
					DEFAULT: 'hsl(var(--popover))',
					foreground: 'hsl(var(--popover-foreground))'
				},
				card: {
					DEFAULT: 'hsl(var(--card))',
					foreground: 'hsl(var(--card-foreground))'
				},
				sidebar: {
					DEFAULT: 'hsl(var(--sidebar-background))',
					foreground: 'hsl(var(--sidebar-foreground))',
					primary: 'hsl(var(--sidebar-primary))',
					'primary-foreground': 'hsl(var(--sidebar-primary-foreground))',
					accent: 'hsl(var(--sidebar-accent))',
					'accent-foreground': 'hsl(var(--sidebar-accent-foreground))',
					border: 'hsl(var(--sidebar-border))',
					ring: 'hsl(var(--sidebar-ring))'
				},
				// Food delivery themed colors
				'food-primary': '#FF4500', // Orange-red
				'food-secondary': '#FFB72B', // Amber
				'food-accent': '#FF8C42', // Light orange
				'food-success': '#4CAF50', // Green
				'food-warning': '#FF9800', // Orange
				'food-danger': '#F44336', // Red
				'food-info': '#2196F3', // Blue
				'food-light': '#F5F5F5', // Off-white
				'food-dark': '#333333', // Dark gray
				// Dark "ink" surfaces used by the marketing landing page
				ink: {
					950: '#07070A',
					900: '#0B0B10',
					850: '#101018',
					800: '#16161F',
					700: '#1F1F2B',
					600: '#2A2A38',
				},
			},
			fontSize: {
				'display': ['clamp(2.75rem, 6vw, 5rem)', { lineHeight: '1.02', letterSpacing: '-0.03em' }],
				'display-sm': ['clamp(2rem, 4.5vw, 3.5rem)', { lineHeight: '1.05', letterSpacing: '-0.02em' }],
			},
			borderRadius: {
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)'
			},
			keyframes: {
				'accordion-down': {
					from: {
						height: '0'
					},
					to: {
						height: 'var(--radix-accordion-content-height)'
					}
				},
				'accordion-up': {
					from: {
						height: 'var(--radix-accordion-content-height)'
					},
					to: {
						height: '0'
					}
				},
				'fade-in': {
					from: { opacity: '0' },
					to: { opacity: '1' }
				},
				'slide-in': {
					from: { transform: 'translateY(20px)', opacity: '0' },
					to: { transform: 'translateY(0)', opacity: '1' }
				},
				'rise-in': {
					from: { transform: 'translateY(28px)', opacity: '0' },
					to: { transform: 'translateY(0)', opacity: '1' }
				},
				'marquee': {
					from: { transform: 'translateX(0)' },
					to: { transform: 'translateX(-50%)' }
				},
				'aurora': {
					'0%, 100%': { transform: 'translate3d(0, 0, 0) scale(1)', opacity: '0.55' },
					'50%': { transform: 'translate3d(0, -3%, 0) scale(1.12)', opacity: '0.8' }
				},
				'shimmer': {
					from: { transform: 'translateX(-100%)' },
					to: { transform: 'translateX(100%)' }
				},
				'pulse-ring': {
					'0%': { transform: 'scale(0.9)', opacity: '0.7' },
					'70%, 100%': { transform: 'scale(1.9)', opacity: '0' }
				},
				'float-y': {
					'0%, 100%': { transform: 'translateY(0)' },
					'50%': { transform: 'translateY(-10px)' }
				}
			},
			animation: {
				'accordion-down': 'accordion-down 0.2s ease-out',
				'accordion-up': 'accordion-up 0.2s ease-out',
				'fade-in': 'fade-in 0.5s ease-out',
				'slide-in': 'slide-in 0.5s ease-out',
				'rise-in': 'rise-in 0.7s cubic-bezier(0.16, 1, 0.3, 1) both',
				'marquee': 'marquee 38s linear infinite',
				'marquee-slow': 'marquee 60s linear infinite',
				'aurora': 'aurora 14s ease-in-out infinite',
				'shimmer': 'shimmer 2.4s ease-in-out infinite',
				'pulse-ring': 'pulse-ring 2.4s cubic-bezier(0.24, 0, 0.38, 1) infinite',
				'float-y': 'float-y 6s ease-in-out infinite'
			}
		}
	},
	plugins: [tailwindcssAnimate],
} satisfies Config;
