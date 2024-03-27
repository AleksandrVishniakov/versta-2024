import React from "react";
import './ChatMessage.css'

interface Message {
    id: number
    message: string
    senderId: number
    receiverId: number
    createdAt: Date
}

const ChatMessage: React.FC<{
    chatterId: number
    msg: Message
}> = (props) => {
    props.msg.createdAt = new Date(props.msg.createdAt)

    return (
        <div
            data-message-sender={props.chatterId === props.msg.senderId ? "user" : "manager"}
            className="chat-message"
        >
            <p className="chat-message-content">{props.msg.message}</p>
            <div className="chat-message-timestamp">{
                props.msg.createdAt.toTimeString().split(" ")[0]
            }</div>
        </div>
    )
}

export default ChatMessage