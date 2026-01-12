// API Configuration
// Update these URLs to match your backend
export const API_URL: string = 'https://app.sumcrowds.com/api/';
export const WS_URL: string = 'wss://app.sumcrowds.com/ws/';

// Environment helper
declare const __DEV__: boolean;
export const isDev: boolean = __DEV__;

// For development, you might use:
// export const API_URL = 'http://10.0.2.2:8080/'; // Android emulator localhost
// export const WS_URL = 'ws://10.0.2.2:8080/ws/';

