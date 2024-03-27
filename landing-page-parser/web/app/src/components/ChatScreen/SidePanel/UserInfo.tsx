import React from "react";
import './UserInfo.css'
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import {Badge, IconButton} from "@mui/material";

interface Chatter {
    id: number
    userId: number
    tempSession: string
    unreadMessagesCount: number
}

interface UserInfoProps {
    chatter: Chatter
    onClick: (id: number)=>void
}

const UserInfo: React.FC<UserInfoProps> = (props) => {
    return (
        <div className="UserInfo" onClick={()=>{props.onClick(props.chatter.id)}}>

            <p className="user-text">
                {`#${props.chatter.id}: ${props.chatter.userId} - ${props.chatter.tempSession}`}
            </p>

            <div style={{display:"flex", gap:"10px", alignItems:"center"}}>
                <IconButton aria-label="expand more" style={{color: "#939393"}}>
                    <MoreHorizIcon/>
                </IconButton>
                <Badge badgeContent={props.chatter.unreadMessagesCount} color="primary">
                </Badge>
            </div>
        </div>

    )
}

export default UserInfo