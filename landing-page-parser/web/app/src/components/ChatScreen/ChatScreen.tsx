import React, {useEffect, useState} from "react";
import {ChatAPI} from "../../api/ChatAPI/ChatAPI";
import './ChatScreen.css'
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import {Button} from "@mui/material";
import SidePanel from "./SidePanel/SidePanel";
import Chat from "./Chat/Chat";

interface Chatter {
    id: number
    userId: number
    tempSession: string
    unreadMessagesCount: number
}

interface ChatScreenProps {
    chatAPI: ChatAPI

    onBack: () => void
    handleError: (msg: string) => void
}

interface Message {
    id: number
    message: string
    senderId: number
    receiverId: number
    createdAt: Date
}

const ChatScreen: React.FC<ChatScreenProps> = (props) => {
    const [chatters, setChatters] = useState(new Array<Chatter>())
    const [selectedUserId, setSelectedUserId] = useState(0)
    const [messages, setMessages] = useState<Message[]>([])
    const [chatterId, setChatterId] = useState(0)

    useEffect(() => {
        setSelectedUserId(0)
    }, []);

    useEffect(() => {
        props.chatAPI.preflightChatRequest()
            .then(() => {
                setChatterId(props.chatAPI.chatterId)
            })
            .catch((error) => {
                props.handleError("chat preflight request failed with error: " + error)
            })
    }, [props]);

    useEffect(() => {
        if (selectedUserId !== 0) {
            props.chatAPI.readAllMessages(selectedUserId)
                .catch((error) => {
                    props.handleError("read all messages error: " + error)
                })

            props.chatAPI.connectChat((msg) => {
                setMessages(messages => [...messages, msg])
            }, selectedUserId).catch((error) => {
                props.handleError("ws connection error:" + error)
            })
        }

        props.chatAPI.getMessages(selectedUserId)
            .then((messages) => {
                setMessages(messages)
            })
            .catch((error) => {
                props.handleError("get messages error: " + error)
            })
    }, [selectedUserId]);

    const getAllChatters = () => {
        props.chatAPI.getAllChatters()
            .then((c) => {
                setChatters(c)
            })
            .catch((error) => {
                props.handleError("get chatters error: " + error)
            })
    }

    useEffect(() => {
        getAllChatters()

        const interval = setInterval(() => {
            getAllChatters()
        }, 10000)

        return () => clearInterval(interval)
    }, [])

    return (
        <main className="ChatScreen">
            <div className="ChatScreen-top-bar">
                <Button
                    variant="text"
                    startIcon={<ArrowBackIcon/>}
                    onClick={() => {
                        props.onBack()
                        props.chatAPI.disconnectChat().catch((error)=>{
                            props.handleError("chat closing error: " + error)
                        })
                    }}
                    className="back-btn"
                >
                    Назад
                </Button>
                <h3 className="chat-title">Чат с пользователями</h3>
            </div>


            <section className="chat-section">
                <div
                    style={{width: "100%"}}
                >
                    {
                        !selectedUserId || selectedUserId === 0 ?
                            <p>Пользователь не выбран</p>
                            :
                            <Chat
                                chatterId={chatterId}
                                messages={messages}
                                open={true}
                                onClose={() => {
                                    props.chatAPI.disconnectChat().catch((error) => {
                                        props.handleError("ws connection error:" + error)
                                    })
                                }}
                                onSend={
                                    (msg: string) => {
                                        props.chatAPI.sendMessage(msg).catch((error) => {
                                            props.handleError("ws error:" + error)
                                        })
                                    }
                                }
                            />
                    }
                </div>

                <SidePanel
                    chatters={chatters}
                    onSelectChatter={(id: number) => {
                        setSelectedUserId(id)
                    }}
                />

            </section>


        </main>
    )
}

export default ChatScreen