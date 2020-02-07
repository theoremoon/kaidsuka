import axios from 'axios';
import env from './env';


const API = axios.create({
  baseURL: env.SERVER_ADDRESS,
  withCredentials: true,
});

export default API;
