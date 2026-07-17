import * as esbuild from 'esbuild'
import { spawnSync, spawn } from 'child_process'
import { createHash } from 'crypto'
import { readFileSync, writeFileSync, mkdirSync, copyFileSync, rmSync } from 'fs'
import { basename } from 'path'

const watch = process.argv.includes('--watch')
const outdir = 'public/assets'
const manifest = {}

// Resolve the tailwindcss CLI entry point directly — avoids symlink issues in npm scripts.
const twCLI = new URL('../node_modules/@tailwindcss/cli/dist/index.mjs', import.meta.url).pathname

mkdirSync(outdir, { recursive: true })

// --- CSS via tailwindcss CLI ---

function buildCSS() {
  const outPath = `${outdir}/main.css`
  const args = [twCLI, '-i', 'internal/assets/src/css/main.css', '-o', outPath]
  if (!watch) args.push('--minify')
  const result = spawnSync(process.execPath, args, { stdio: 'inherit' })
  if (result.error) throw result.error
  if (result.status !== 0) process.exit(result.status ?? 1)
  if (watch) {
    manifest['main.css'] = 'main.css'
  } else {
    const css = readFileSync(outPath)
    const hash = createHash('sha256').update(css).digest('hex').slice(0, 8)
    const hashedName = `main-${hash}.css`
    copyFileSync(outPath, `${outdir}/${hashedName}`)
    rmSync(outPath)
    manifest['main.css'] = hashedName
  }
}

// --- JS via esbuild ---

function buildJS() {
  if (watch) {
    esbuild.buildSync({
      entryPoints: ['internal/assets/src/js/main.js'],
      bundle: true,
      outfile: `${outdir}/main.js`,
    })
    manifest['main.js'] = 'main.js'
  } else {
    const result = esbuild.buildSync({
      entryPoints: ['internal/assets/src/js/main.js'],
      bundle: true,
      minify: true,
      entryNames: '[name]-[hash]',
      outdir,
      metafile: true,
    })
    for (const [outputPath, output] of Object.entries(result.metafile.outputs)) {
      if (output.entryPoint) {
        manifest[basename(output.entryPoint)] = basename(outputPath)
      }
    }
  }
}

function writeManifest() {
  writeFileSync(`${outdir}/manifest.json`, JSON.stringify(manifest, null, 2))
  console.log('manifest:', manifest)
}

buildCSS()
buildJS()
writeManifest()

if (watch) {
  console.log('starting watch mode...')

  // Tailwind in watch mode — spawns node directly with the resolved CLI path.
  const tw = spawn(
    process.execPath,
    [twCLI, '-i', 'internal/assets/src/css/main.css', '-o', `${outdir}/main.css`, '--watch'],
    { stdio: 'inherit' }
  )
  tw.on('exit', code => process.exit(code ?? 0))

  // esbuild JS watch
  const ctx = await esbuild.context({
    entryPoints: ['internal/assets/src/js/main.js'],
    bundle: true,
    outfile: `${outdir}/main.js`,
  })
  await ctx.watch()
}
