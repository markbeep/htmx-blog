import { defineConfig, presetWind3, transformerDirectives } from "unocss";

export default defineConfig({
  cli: {
    entry: {
      patterns: ["components/**/*.templ", "static/uno.css"],
      outFile: "static/main.css",
    },
  },
  presets: [presetWind3()],
  transformers: [transformerDirectives()],
  theme: {
    fontFamily: {
      gothic: "Gothic, sans-serif",
    },
  },
});
