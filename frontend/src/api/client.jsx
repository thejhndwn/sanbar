// src/api/client.js

class ApiClient {
  async request(endpoint, options = {}) {
      console.log("is it this?")
    const token = localStorage.getItem('sanbarToken');
      console.log("no, but why tho")
    const url = `${process.env.REACT_APP_API_BASE}${endpoint}`;
      console.log("surely ain't this")
    const config = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
          'Authorization': token,
        ...options.headers,
      },
    };

    console.log("we are sending shit to:", url)

    const response = await fetch(url, config);
    console.log("got a response")

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || response.statusText);
    }

    if (response.status === 204) {
      return {};
    }

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
