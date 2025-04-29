import { defineConfig } from 'vite'
import { resolve } from 'path'
import { readdirSync } from 'fs'

const htmlFiles = readdirSync('./pages')
    .filter(file => file.endsWith('.html'))
    .reduce((entries, file) => {
        const name = file.replace('.html', '')
        entries[name] = resolve(__dirname, 'pages', file)
        return entries
    }, {})

export default defineConfig({
    build: {
        rollupOptions: {
            input: htmlFiles
        },
        outDir: '../dir',
        emptyOutDir: true
    }
})
