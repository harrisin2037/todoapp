let API_BASE_URL;

if (typeof import.meta !== 'undefined' && import.meta.env) {
    API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    console.log('Successfully read VITE_API_BASE_URL:', API_BASE_URL);
} else if (typeof process !== 'undefined' && process.env) {
    API_BASE_URL = process.env.VITE_API_BASE_URL;
    console.log('Successfully read VITE_API_BASE_URL from process.env:', API_BASE_URL);
} else {
    API_BASE_URL = 'http://localhost/api';
    console.log('Using fallback API_BASE_URL:', API_BASE_URL);
}

export { API_BASE_URL };
