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
    "http://localhost:8001"
)
const ordersAPI = new OrdersAPI(
    "http://localhost:8000",
    authAPI,
)

const chatAPI = new ChatAPI(
    authAPI,
    "http://localhost:8003"
)

chatAPI.preflightChatRequest()
    .catch((error) => {
        console.error("chat preflight request failed with error: " + error)
    })

const appHostField = document.querySelector("#app-host") as HTMLInputElement

console.log(appHostField?.value)

root.render(
    <React.StrictMode>
        <App
            ordersAPI={ordersAPI}
            authAPI={authAPI}
            chatAPI={chatAPI}
        />
    </React.StrictMode>
);
