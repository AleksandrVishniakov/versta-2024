import React, {ChangeEvent, useEffect, useState} from "react";
import {InputAdornment, TextField} from "@mui/material";
import TagIcon from '@mui/icons-material/Tag';

interface VerificationCodeInputProps {
    label?: string
    className?: string
    name?: string
    id?: string
    value?: string
    onChange?: (code: string) => void
    disabled?: boolean
}

const VerificationCodeInput: React.FC<VerificationCodeInputProps> = (props) => {
    const [value, setValue] = useState(props.value ? props.value : "")

    const handleInputValueChange = (evt: ChangeEvent<HTMLInputElement>) => {
        setValue(evt.target.value)

        if (props.onChange) {
            props.onChange(evt.target.value)
        }
    }

    useEffect(() => {
        if (!props.value) return

        setValue(props.value)
    }, [props.value]);

    return (
        <TextField
            label={props.label ? props.label : "Verification Code"}
            id={props.id ? props.id : "v-code-input"}
            name={props.name ? props.name : "v-code-input"}
            type="number"
            InputProps={{
                startAdornment: (
                    <InputAdornment position="start">
                        <TagIcon />
                    </InputAdornment>
                ),
            }}
            variant="outlined"
            fullWidth
            value={value}
            onChange={handleInputValueChange}
            className={props.className}
            disabled={props.disabled ? props.disabled : false}
        />
    )
}

export default VerificationCodeInput