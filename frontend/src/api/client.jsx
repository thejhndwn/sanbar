// src/api/client.js

class ApiClient {
    async request(endpoint, options = {}) {
        const token = localStorage.getItem('sanbarToken');
        const url = `${import.meta.env.VITE_REACT_APP_API_BASE}${endpoint}`;
        const config = {
            ...options,
            headers: {
                'content-type': 'application/json',
                'authorization': token,
                ...options.headers,
            },
        };

        console.log("we are sending shit to:", url)

        const response = await fetch(url, config);
        console.log("got a response", response)

        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            console.log("definitely not this one")
            throw new Error(error.message || response.statusText);
        }

        if (response.status === 204) {
            console.log("we shouldn't really go in here")
            return {};
        }
        console.log("return?")

        return response.json();
    }

    get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    }

    post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    put(endpoint, data) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    delete(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    }
}

export const api = new ApiClient();
