// Simple dev runner for Windows/macOS/Linux.
// Starts: matchmaker (Go), gateway (Go), web (Next.js dev), launcher (Electron dev)
const { spawn } = require("child_process");
const path = require("path");

const ROOT = path.resolve(__dirname, "..");
const procs = [];

function run(name, cmd, args, opts = {}) {
  const p = spawn(cmd, args, {
    cwd: opts.cwd || ROOT,
    stdio: "inherit",
    shell: process.platform === "win32",
  });
  procs.push({ name, p });
  p.on("exit", (code) => console.log(`[${name}] exited with ${code}`));
}

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
function shutdown() {
  console.log("\nshutting down...");
  procs.forEach(({ p }) => {
    if (p && !p.killed) p.kill();
  });
  process.exit(0);
}

// 1) matchmaker
run("matchmaker", "go", ["run", "./cmd/matchmaker"], {
  cwd: path.join(ROOT, "server", "matchmaker"),
});

// 2) gateway
run("gateway", "go", ["run", "./cmd/gateway"], {
  cwd: path.join(ROOT, "server", "gateway"),
});

// 3) web (Next.js dev)
run("web", "npm", ["run", "dev"], { cwd: path.join(ROOT, "web") });

// 4) launcher (Electron dev)
run("launcher", "npm", ["run", "dev"], { cwd: path.join(ROOT, "launcher") });
