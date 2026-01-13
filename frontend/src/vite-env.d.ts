/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APIURL: string
  readonly VITE_WSURL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
