import { contextBridge } from "electron";

contextBridge.exposeInMainWorld("FFB", {
  apiHostMatchmaker: "http://localhost:8081",
  apiHostGateway: "http://localhost:8080",
});
