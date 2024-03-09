import React from 'react';
import './App.css';
import FormSection from "./components/FormSection/FormSection";
import OrdersAPI from "./api/OrdersAPI/OrdersAPI";

interface AppProps {
    ordersAPI: OrdersAPI
}

const App: React.FC<AppProps> = (props) => {
    const handleError = (message: string) => {
        console.error(message)
    }

    return (
        <div className="App">
            <FormSection
                ordersAPI={props.ordersAPI}
                handleError={handleError}
            />
        </div>
    );
}

export default App;
