import { app, BrowserWindow } from "electron";
import { join } from "path";

function createWindow() {
  const win = new BrowserWindow({
    width: 1100,
    height: 780,
    webPreferences: {
      preload: join(__dirname, "preload.js"),
      contextIsolation: true,
      nodeIntegration: false
    }
  });

  // HTML у dist/renderer — копіюється скриптом build:assets
  const htmlPath = join(__dirname, "renderer", "index.html");
  win.loadFile(htmlPath);
}

app.whenReady().then(() => {
  createWindow();
  app.on("activate", () => { if (BrowserWindow.getAllWindows().length === 0) createWindow(); });
});
app.on("window-all-closed", () => { if (process.platform !== "darwin") app.quit(); });
