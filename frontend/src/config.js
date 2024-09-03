console.log('import.meta:', import.meta);
console.log('import.meta.env:', import.meta.env);

let API_BASE_URL;

if (import.meta.env.VITE_API_BASE_URL) {
    API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    console.log('Successfully read VITE_API_BASE_URL:', API_BASE_URL);
} else if (import.meta.env.MODE === 'development') {
    API_BASE_URL = 'http://localhost/api';
    console.log('Using development API_BASE_URL:', API_BASE_URL);
} else {
    API_BASE_URL = 'http://localhost/api';
    console.log('Using fallback API_BASE_URL:', API_BASE_URL);
}

export { API_BASE_URL };
