import {InputAdornment, TextField } from "@mui/material";
import React, {ChangeEvent, useState } from "react";
import AlternateEmailIcon from "@mui/icons-material/AlternateEmail";

interface EmailInputProps {
    label?: string
    name?: string
    id?: string
    value?: string
    onChange?: (email: string) => void
}
const EmailInput: React.FC<EmailInputProps> = (props) => {
    const [value, setValue] = useState(props.value ? props.value : "")
    
    const handleInputValueChange = (evt: ChangeEvent<HTMLInputElement>) => {
        setValue(evt.target.value)
        
        if (props.onChange) {
            props.onChange(evt.target.value)
        }
    }
    
    return (
        <TextField
            label={props.label ? props.label : "Email"}
            id={props.id ? props.id : "email-input"}
            name={props.name ? props.name : "email-input"}
            type="email"
            InputProps={{
            startAdornment: (
                <InputAdornment position="start">
                    <AlternateEmailIcon />
                </InputAdornment>
                ),
            }}
            variant="outlined"
            fullWidth
            value={value}
            onChange={handleInputValueChange}
        />
    )
}

export default EmailInput