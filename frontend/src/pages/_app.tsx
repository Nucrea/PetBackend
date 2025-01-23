import "@/styles/globals.css";
import type { AppProps } from "next/app";
import '@mantine/core/styles.css';
import "./global.css"

export default function App({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />;
}
