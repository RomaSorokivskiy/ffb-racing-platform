export {};

declare global {
  interface Window {
    FFB?: {
      apiHostMatchmaker: string;
      apiHostGateway: string;
    };
  }
}
