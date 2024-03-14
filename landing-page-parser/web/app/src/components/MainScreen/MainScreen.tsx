import OrdersAPI from "../../api/OrdersAPI/OrdersAPI";
import React from "react";
import FormSection from "./FormSection/FormSection";
import UserAccount from "./UserAccount/UserAccount";

interface MainScreenProps {
    ordersAPI: OrdersAPI
    userEmail: string
    handleError: (message: string) => void
    onAuth: (email: string) => void

    onOpenProfile: ()=>void,
    onAuthUser: ()=>void
}

const MainScreen: React.FC<MainScreenProps> = (props) => {

    return (
        <main>
            <UserAccount
                userEmail={props.userEmail}
                onOpenProfile={props.onOpenProfile}
                onAuthUser={props.onAuthUser}
            />

            <FormSection
                ordersAPI={props.ordersAPI}
                handleError={props.handleError}
                userEmail={props.userEmail}
                onAuth={props.onAuth}
            />
        </main>
    )
}

export default MainScreen
