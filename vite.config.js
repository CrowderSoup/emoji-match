import { defineConfig } from 'vite';

export default defineConfig({
  root: './public',
  build: {
    outDir: '../dist',
    assetsDir: '.',
    rollupOptions: {
      input: {
        main: './public/index.html',
      },
    },
  },
});
