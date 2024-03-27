import React, {useEffect, useState} from 'react';
import './App.css';
import OrdersAPI from "./api/OrdersAPI/OrdersAPI";
import MainScreen from "./components/MainScreen/MainScreen";
import authAPI from "./api/AuthAPI/AuthAPI";
import ProfileScreen from "./components/ProfileScreen/ProfileScreen";
import LoginScreen from "./components/LoginScreen/LoginScreen";
import {ChatAPI} from "./api/ChatAPI/ChatAPI";
import {UserStatus} from "./api/AuthAPI/Statuses";
import ChatScreen from "./components/ChatScreen/ChatScreen";

enum Screens {
    Main,
    Profile,
    Login,
    ChatWithUsers
}


interface AppProps {
    ordersAPI: OrdersAPI
    authAPI: authAPI
    chatAPI: ChatAPI
}

const App: React.FC<AppProps> = (props) => {
    const [userEmail, setUserEmail] = useState("")
    const [_, setUserId] = useState(0)
    const [currentScreen, setCurrentScreen] = useState(Screens.Main)
    const [userStatus, setUserStatus] = useState(UserStatus.StatusUser)

    const getUserProfile = async () => {
        try {
            return await props.authAPI.getProfile()
        } catch (e: any) {
            handleError(e.toString())
        }
    }

    useEffect(() => {

        getUserProfile().then((user) => {
            if (!user) return

            if (user.status === UserStatus.StatusAdmin) {
                setUserStatus(UserStatus.StatusAdmin)
            }

            setUserEmail(user.email)
        })
    }, [props.authAPI]);

    const handleError = (message: string) => {
        console.error(message)
    }

    const handleAuth = (email: string) => {
        getUserProfile().then((user) => {
            if (!user) return

            if (user.status === UserStatus.StatusAdmin) {
                setUserStatus(UserStatus.StatusAdmin)
            }

            setUserEmail(user.email)
        })
    }

    const appNavigation = () => {
        switch (currentScreen) {
            case Screens.Main:
                return (
                    <MainScreen
                        ordersAPI={props.ordersAPI}
                        chatAPI={props.chatAPI}

                        userEmail={userEmail}
                        handleError={handleError}
                        onAuth={handleAuth}
                        userStatus={userStatus}
                        onOpenUsersChat={userStatus === UserStatus.StatusAdmin ? () => {
                            setCurrentScreen(Screens.ChatWithUsers)
                        } : () => {
                        }}

                        onAuthUser={() => {
                            setCurrentScreen(Screens.Login)
                        }}
                        onOpenProfile={() => {
                            setCurrentScreen(Screens.Profile)
                        }}
                    />
                )

            case Screens.Profile:
                return (
                    <ProfileScreen
                        authAPI={props.authAPI}
                        ordersAPI={props.ordersAPI}
                        handleError={handleError}
                        userStatus={userStatus}
                        onOpenUsersChat={userStatus === UserStatus.StatusAdmin ? () => {
                            setCurrentScreen(Screens.ChatWithUsers)
                        } : () => {
                        }}
                        onBack={() => {
                            setCurrentScreen(Screens.Main)
                        }}
                        onLogout={() => {
                            setCurrentScreen(Screens.Login)
                        }}
                    />
                )

            case Screens.Login:
                return (
                    <LoginScreen
                        authAPI={props.authAPI}
                        handleError={handleError}
                        onBack={() => {
                            setCurrentScreen(Screens.Main)
                        }}
                        onAuth={(email: string, userId: number) => {
                            handleAuth(email)
                            setUserId(userId)
                            setCurrentScreen(Screens.Profile)
                        }}
                    />
                )

            case Screens.ChatWithUsers:
                return (
                    <ChatScreen
                        chatAPI={props.chatAPI}
                        onBack={() => {
                            setCurrentScreen(Screens.Profile)
                        }}
                        handleError={handleError}
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
