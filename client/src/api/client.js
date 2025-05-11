import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8000/api/v1',
  withCredentials: true,    // чтобы проставлялись jwt_access и jwt_refresh
});

export default api;
