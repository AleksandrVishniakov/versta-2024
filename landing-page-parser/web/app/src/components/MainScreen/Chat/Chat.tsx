import React, {useEffect, useRef, useState} from "react";
import './Chat.css'
import {IconButton, TextField} from "@mui/material";
import CloseIcon from '@mui/icons-material/Close';
import SendIcon from '@mui/icons-material/Send';
import ChatMessage from "../../common/ChatMessage/ChatMessage";
import DateMessage from "../../common/ChatMessage/DateMessage";

interface Message {
    id: number
    message: string
    senderId: number
    receiverId: number
    createdAt: Date
}

interface ChatProps {
    chatterId: number
    messages: Array<Message>

    open: boolean

    onClose: () => void
    onSend: (msg: string) => void
}

const Chat: React.FC<ChatProps> = (props) => {
    const [msgInputValue, setMsgInputValue] = useState("")

    const messagesEndRef = useRef<null | HTMLDivElement>(null)

    const scrollToBottom = (smooth: boolean = false) => {
        messagesEndRef.current?.scrollIntoView({behavior: smooth ? "smooth" : "auto"})
    }

    useEffect(() => {
        scrollToBottom()
    }, [props.open]);

    useEffect(() => {
        scrollToBottom(true)
    }, [props.messages]);

    let lastDate = new Date("01.01.1980")
    return (
        props.open ?
            <section className={"Chat"}>
                <div className={"ChatTopBar"}>
                    <h3>Поддержка</h3>
                    <IconButton aria-label="close" onClick={props.onClose} style={{height: "fit-content"}}>
                        <CloseIcon/>
                    </IconButton>
                </div>

                <div className={"ChatLog"}>
                    {
                        props.messages && props.messages.length > 0 ?
                            props.messages.map((msg, index) => {
                                if (index !== 0) {
                                    lastDate = props.messages[index - 1].createdAt
                                }

                                return (
                                    <div key={`block-${index}`}>
                                        {
                                            isDiffDates(lastDate, msg.createdAt) ?
                                                <DateMessage date={msg.createdAt} key={`date-for-message-${index}`}/>
                                                : null
                                        }

                                        <ChatMessage
                                            chatterId={props.chatterId}
                                            msg={msg}
                                            key={`message-${index}`}
                                        />
                                    </div>
                                )
                            }) : <p className={"no-messages-text"}>Сообщений пока нет</p>
                    }
                    <div ref={messagesEndRef}/>
                </div>

                <form className={"ChatInput"}
                      onSubmit={(evt) => {
                          evt.preventDefault()

                          if (msgInputValue.length === 0) return

                          props.onSend(msgInputValue)

                          setMsgInputValue("")
                      }}
                >
                    <TextField
                        id="standard-basic"
                        label="Сообщение"
                        variant="standard"
                        className="ChatMessageInput"
                        value={msgInputValue}
                        onChange={(evt) => {
                            setMsgInputValue(evt.target.value)
                        }}
                    />

                    <IconButton
                        aria-label="send"
                        size="large"
                        color="primary"
                        type="submit"
                    >
                        <SendIcon/>
                    </IconButton>
                </form>
            </section>
            : null
    )
}

const isDiffDates = (oldDate: Date, newDate: Date): boolean => {
    const date1 = new Date(oldDate)
    const date2 = new Date(newDate)

    if (date1.getFullYear() !== date2.getFullYear()) {
        return true
    }

    if (date1.getMonth() !== date2.getMonth()) {
        return true
    }

    return date1.getDay() !== date2.getDay()
}

export default Chat