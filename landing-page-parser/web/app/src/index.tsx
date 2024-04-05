import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import OrdersAPI from "./api/OrdersAPI/OrdersAPI";
import AuthAPI from "./api/AuthAPI/AuthAPI";
import {ChatAPI} from "./api/ChatAPI/ChatAPI";

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

const authAPI = new AuthAPI(
    `http://${document.location.hostname}:8001`
)
const ordersAPI = new OrdersAPI(
    `http://${document.location.hostname}:8000`,
    authAPI,
)

const chatAPI = new ChatAPI(
    authAPI,
    `http://${document.location.hostname}:8003`
)

chatAPI.preflightChatRequest()
    .catch((error) => {
        console.error("chat preflight request failed with error: " + error)
    })

let colorScheme: "light" | "dark" = "light"

if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    colorScheme = "dark"
}

root.render(
    <React.StrictMode>
        <App
            ordersAPI={ordersAPI}
            authAPI={authAPI}
            chatAPI={chatAPI}
            theme={colorScheme}
        />
    </React.StrictMode>
);
