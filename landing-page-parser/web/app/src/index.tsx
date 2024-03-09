import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import OrdersAPI from "./api/OrdersAPI/OrdersAPI";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const ordersAPI = new OrdersAPI("http://localhost:8000")

root.render(
  <React.StrictMode>
    <App
        ordersAPI={ordersAPI}
    />
  </React.StrictMode>
);