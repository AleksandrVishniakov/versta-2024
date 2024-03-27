import React, {useState} from "react";
import './SidePanel.css'
import KeyboardDoubleArrowLeftIcon from '@mui/icons-material/KeyboardDoubleArrowLeft';
import KeyboardDoubleArrowRightIcon from '@mui/icons-material/KeyboardDoubleArrowRight';
import {IconButton} from "@mui/material";
import UserInfo from "./UserInfo";
import {ChatAPI} from "../../../api/ChatAPI/ChatAPI";

interface Chatter {
    id: number
    userId: number
    tempSession: string
    unreadMessagesCount: number
}

interface SidePanelProps {
    chatters: Array<Chatter>
    onSelectChatter: (id: number)=>void
}

const SidePanel: React.FC<SidePanelProps> = (props) => {
    const [open, setOpen] = useState(true)

    return (
        <div className={`SidePanel ${open?"":"closed"}`}>
            <div className="SidePanel-top-bar">
                <IconButton aria-label="close panel" style={{color: "#939393"}}
                            className="close-btn"
                            onClick={()=>{setOpen(!open)}}
                >
                    {
                        open?
                            <KeyboardDoubleArrowLeftIcon/>
                            :
                            <KeyboardDoubleArrowRightIcon/>
                    }
                </IconButton>
                {
                    open ?
                        <p className="side-panel-title">
                            Пользователи
                        </p>
                        :null
                }

            </div>

            <div className="clients-list">
                {
                    open && props.chatters && props.chatters.length > 0 ?
                        props.chatters.map((chatter, index)=> {
                            return (
                                <UserInfo
                                    chatter={chatter}
                                    onClick={props.onSelectChatter}
                                    key={`chatter-${index}`}
                                />
                            )
                        }) : null
                }
            </div>

        </div>
    )
}

export default SidePanel