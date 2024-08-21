import App from './App.svelte';
import { API_BASE_URL } from './config.js';

const app = new App({
    target: document.body,
    props: {
        apiBaseUrl: API_BASE_URL
    }
});

export default app;
