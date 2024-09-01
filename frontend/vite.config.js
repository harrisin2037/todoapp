import { defineConfig, loadEnv } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { resolve } from 'path'

export default defineConfig(({ command, mode }) => {
    const env = loadEnv(mode, process.cwd(), '')
    console.log('Loaded env:', env)
    console.log('VITE_API_BASE_URL:', env.VITE_API_BASE_URL)

    return {
        plugins: [svelte()],
        define: {
            'import.meta.env.VITE_API_BASE_URL': JSON.stringify(env.VITE_API_BASE_URL),
        },
        build: {
            outDir: 'dist',
            rollupOptions: {
                input: resolve(__dirname, 'index.html'),
            },
        },
        preview: {
            port: 3000,
        },
    }
})