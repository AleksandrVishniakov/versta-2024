import React, {ChangeEvent, useEffect, useState} from "react";
import './ProfileHeader.css'
import {Button, IconButton} from "@mui/material";
import SaveIcon from '@mui/icons-material/Save';

interface ProfileHeaderProps {
    email: string
    name: string
    createdAt: Date

    onUpdateName: (name: string) => void
    onLogout: () => void
}


const ProfileHeader: React.FC<ProfileHeaderProps> = (props) => {
    const [userName, setUserName] = useState(props.name)
    const [userNameInputValue, setUserNameInputValue] = useState(props.name)
    const [isDataSubmitting, setDataSubmitting] = useState(false)

    useEffect(() => {
        setUserName(props.name)
        setUserNameInputValue(props.name)
    }, [props.name]);

    const handleInputChange = (value: string, setValue: React.Dispatch<React.SetStateAction<string>>) => {
        setValue(value)
    }

    const handleNameFormSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (isDataSubmitting) return
        setDataSubmitting(true)

        props.onUpdateName(userNameInputValue)

        setDataSubmitting(false)
    }

    return (
        <section className="ProfileHeader">
            <form className="update-name-form"
                  onSubmit={handleNameFormSubmit}
            >
                <input className="name-input"
                       value={userNameInputValue}
                       placeholder="ваше имя"
                       onChange={(evt: ChangeEvent<HTMLInputElement>) => {
                           handleInputChange(evt.target.value, setUserNameInputValue)
                       }}
                />

                {
                    (userNameInputValue !== userName) &&
                    (userNameInputValue !== "") ?
                        <IconButton
                            aria-label="delete"
                            size="large"
                            style={{color: "#006aff"}}
                            type="submit"
                            disabled={isDataSubmitting &&
                                userNameInputValue === userName &&
                                userNameInputValue === ""
                            }
                        >
                            <SaveIcon fontSize="inherit"/>
                        </IconButton>
                        : null
                }
            </form>

            <p>{props.email}</p>
            <p>Создан: <span>{props.createdAt.toString()}</span></p>

            <Button
                variant="outlined"
                color="error"
                onClick={props.onLogout}
                style={{width: "fit-content"}}
            >
                Выйти
            </Button>
        </section>
    )

}

export default ProfileHeader