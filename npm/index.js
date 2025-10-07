#!/usr/bin/env node
const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');

const exe = process.platform === 'win32' ? 'chainenv.exe' : 'chainenv';
const binPath = path.join(__dirname, exe);

if (!fs.existsSync(binPath)) {
  console.error(`Binary not found: ${binPath}`);
  process.exit(1);
}

if (process.platform !== 'win32') {
  try { fs.chmodSync(binPath, 0o755); } catch {}
}

const child = spawn(binPath, process.argv.slice(2), { stdio: 'inherit', windowsHide: true });
child.on('exit', (code) => process.exit(code == null ? 1 : code));
child.on('error', (err) => { console.error(err.message); process.exit(1); });