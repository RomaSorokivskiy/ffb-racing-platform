const fs = require("fs");
const path = require("path");

const SRC_DIR = path.resolve(__dirname, "..", "src", "renderer");
const DST_DIR = path.resolve(__dirname, "..", "dist", "renderer");

function copyOnce() {
  fs.mkdirSync(DST_DIR, { recursive: true });
  const files = fs.readdirSync(SRC_DIR).filter(f => f.endsWith(".html"));
  for (const file of files) {
    const src = path.join(SRC_DIR, file);
    const dst = path.join(DST_DIR, file);
    fs.copyFileSync(src, dst);
    console.log(`[copy-assets] ${src} -> ${dst}`);
  }
}

if (process.argv.includes("--watch")) {
  copyOnce();
  fs.watch(SRC_DIR, { recursive: false }, (_ev, file) => {
    if (file && file.endsWith(".html")) copyOnce();
  });
} else {
  copyOnce();
}
