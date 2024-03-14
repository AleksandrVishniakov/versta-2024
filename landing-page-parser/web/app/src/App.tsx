import React, {useEffect, useState} from 'react';
import './App.css';
import OrdersAPI from "./api/OrdersAPI/OrdersAPI";
import MainScreen from "./components/MainScreen/MainScreen";
import authAPI from "./api/AuthAPI/AuthAPI";
import ProfileScreen from "./components/ProfileScreen/ProfileScreen";
import LoginScreen from "./components/LoginScreen/LoginScreen";

enum Screens {
    Main,
    Profile,
    Login,
}

interface AppProps {
    ordersAPI: OrdersAPI
    authAPI: authAPI
}

const App: React.FC<AppProps> = (props) => {
    const [userEmail, setUserEmail] = useState("")
    const [currentScreen, setCurrentScreen] = useState(Screens.Main)

    useEffect(() => {
        const getUserProfile = async () => {
            try {
                return await props.authAPI.getProfile()
            }
            catch (e: any) {
                handleError(e.toString())
            }
        }

        getUserProfile().then((user) => {
            if (!user) return

            setUserEmail(user.email)
        })
    }, [props.authAPI]);

    const handleError = (message: string) => {
        console.error(message)
    }

    const handleAuth = (email: string) => {
        setUserEmail(email)
    }

    const appNavigation = () => {
        switch (currentScreen) {
            case Screens.Main:
                return (
                    <MainScreen
                        ordersAPI={props.ordersAPI}
                        userEmail={userEmail}
                        handleError={handleError}
                        onAuth={handleAuth}

                        onAuthUser={()=>{setCurrentScreen(Screens.Login)}}
                        onOpenProfile={()=>{setCurrentScreen(Screens.Profile)}}
                    />
                )

            case Screens.Profile:
                return (
                    <ProfileScreen
                        authAPI={props.authAPI}
                        ordersAPI={props.ordersAPI}
                        handleError={handleError}
                        onBack={()=>{setCurrentScreen(Screens.Main)}}
                        onLogout={()=>{setCurrentScreen(Screens.Login)}}
                    />
                )

            case Screens.Login:
                return (
                    <LoginScreen
                        authAPI={props.authAPI}
                        handleError={handleError}
                        onBack={()=>{setCurrentScreen(Screens.Main)}}
                        onAuth={(email)=>{
                            handleAuth(email)
                            setCurrentScreen(Screens.Profile)
                        }}
                    />
                )
        }
    }

    return (
        <div className="App">
            {appNavigation()}
        </div>
    );
}

export default App;
