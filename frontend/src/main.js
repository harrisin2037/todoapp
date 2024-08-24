import App from './App.svelte';
import { API_BASE_URL } from './config.js';
import { isLoading } from './loadingStore';

const app = new App({
    target: document.body,
    props: {
        apiBaseUrl: API_BASE_URL,
    }
});

setTimeout(() => {
    isLoading.set(false);
}, 3000);

export default app;
