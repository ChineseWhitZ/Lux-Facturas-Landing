import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        ink: "#10233f",
        cloud: "#f6f8fb",
        mint: "#19b886",
        lime: "#d8f263",
        coral: "#f58d72",
        night: "#090f1a",
        graphite: "#111827",
        steel: "#1f2937",
        frost: "#eef7ff",
        cyan: "#67e8f9",
        gold: "#f6d36b",
      },
      boxShadow: {
        soft: "0 24px 80px rgba(16, 35, 63, 0.14)",
        card: "0 18px 45px rgba(16, 35, 63, 0.08)",
        glow: "0 30px 110px rgba(103, 232, 249, 0.16)",
      },
    },
  },
  plugins: [],
};

export default config;
