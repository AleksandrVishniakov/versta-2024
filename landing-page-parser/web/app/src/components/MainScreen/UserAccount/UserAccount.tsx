import React from "react";
import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined';
import {Button} from "@mui/material";

interface UserAccountProps {
    userEmail: string
    onOpenProfile: () => void,
    onAuthUser: () => void
}

const UserAccount: React.FC<UserAccountProps> = (props) => {
    return (
        <div className="UserAccount">
            {
                props.userEmail !== "" ?
                    <button onClick={props.onOpenProfile}>
                        <AccountCircleOutlinedIcon/>
                    </button>
                    :
                    <Button
                        variant="contained"
                        onClick={props.onAuthUser}
                    >
                        Вход
                    </Button>
            }
        </div>
    )
}

export default UserAccount