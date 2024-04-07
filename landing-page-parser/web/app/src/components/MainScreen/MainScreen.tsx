import OrdersAPI from "../../api/OrdersAPI/OrdersAPI";
import React, {useEffect, useState} from "react";
import FormSection from "./FormSection/FormSection";
import UserAccount from "./UserAccount/UserAccount";
import {ChatAPI} from "../../api/ChatAPI/ChatAPI";
import {Badge, Fab} from "@mui/material";
import ChatIcon from '@mui/icons-material/Chat';
import './MainScreen.css'
import Chat from "./Chat/Chat";
import {UserStatus} from "../../api/AuthAPI/Statuses";
import TitleSection from "./TitleSection/TitleSection";
import InfoSection from "./InfoSection/InfoSection";

interface Message {
    id: number
    message: string
    senderId: number
    receiverId: number
    createdAt: Date
}

interface MainScreenProps {
    ordersAPI: OrdersAPI
    chatAPI: ChatAPI

    userEmail: string
    userStatus: UserStatus
    handleError: (message: string) => void
    onAuth: (email: string) => void

    onOpenUsersChat: () => void
    onOpenProfile: () => void,
    onAuthUser: () => void
}

const MainScreen: React.FC<MainScreenProps> = (props) => {
    const [chatterId, setChatterId] = useState(0)
    const [messages, setMessages] = useState<Message[]>([])
    const [chatOpen, setChatOpen] = useState(false)
    const [unreadMessages, setUnreadMessages] = useState(0)

    useEffect(() => {
        props.chatAPI.preflightChatRequest()
            .then(() => {
                setChatterId(props.chatAPI.chatterId)
            })
            .catch((error) => {
                props.handleError("chat preflight request failed with error: " + error)
            })
    }, [props]);



    const getUnreadMessages = () => {
        props.chatAPI.getUnreadMessagesCount()
            .then((count) => {
                setUnreadMessages(count)
            })
            .catch((error) => {
                props.handleError("get unread messages error: " + error)
            })
    }

    useEffect(() => {
        getUnreadMessages()

        const interval = setInterval(() => {
            getUnreadMessages()
        }, 10000)

        return () => clearInterval(interval)
    }, [props])

    const onChatOpen = () => {
        if (props.userStatus === UserStatus.StatusAdmin) {
            props.onOpenUsersChat()
            return
        }

        props.chatAPI.getMessages()
            .then((messages) => {
                setMessages(messages)
            })
            .catch((error) => {
                props.handleError("get messages error: " + error)
            })

        props.chatAPI.readAllMessages()
            .then(() => {
                setUnreadMessages(0)
            })
            .catch((error) => {
                props.handleError("read all messages error: " + error)
            })

        setChatOpen(true)
        props.chatAPI.connectChat((msg) => {
            setMessages(messages => [...messages, msg])
        }).catch((error) => {
            props.handleError("ws connection error:" + error)
        })
    }

    return (
        <main className="MainScreen">
            <Chat
                open={chatOpen}
                messages={messages}
                chatterId={chatterId}

                onSend={(msg: string) => {
                    props.chatAPI.sendMessage(msg).catch((error) => {
                        props.handleError("ws error:" + error)
                    })
                }}
                onClose={() => {
                    setChatOpen(false)
                    props.chatAPI.disconnectChat().catch((error) => {
                        props.handleError("ws connection error:" + error)
                    })
                }}
            />

            {
                !chatOpen ?
                    <Badge badgeContent={unreadMessages} color="primary" className={"fab-right-bottom"}>
                        <Fab color="default" aria-label="add" onClick={() => {
                            onChatOpen()
                        }}>
                            <ChatIcon/>
                        </Fab>
                    </Badge>
                    : null
            }


            <TitleSection
                userAccountProps={{
                    userEmail: props.userEmail,
                    onOpenProfile: props.onOpenProfile,
                    onAuthUser: props.onAuthUser
                }}
            />

            <InfoSection />

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
