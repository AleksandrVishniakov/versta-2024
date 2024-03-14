import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import OrdersAPI from "./api/OrdersAPI/OrdersAPI";
import AuthAPI from "./api/AuthAPI/AuthAPI";

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

const authAPI = new AuthAPI(
    "http://localhost:8001"
)
const ordersAPI = new OrdersAPI(
    "http://localhost:8000",
    authAPI,
)

root.render(
    <React.StrictMode>
        <App
            ordersAPI={ordersAPI}
            authAPI={authAPI}
        />
    </React.StrictMode>
);
