import {Button, Checkbox, FormControlLabel, TextField} from "@mui/material";
import React, {ChangeEvent, useState} from "react";
import "./FormSection.css"
import EmailInput from "./Inputs/EmailInput";

const FormSection: React.FC = (props) => {
    const [emailInputValue, setEmailInputValue] = useState("")
    const [nameInputValue, setNameInputValue] = useState("")
    const [informationInputValue, setInformationInputValue] = useState("")
    const [isConfirmationCheckboxChecked, setConfirmationCheckboxState] = useState(true)

    const handleInputChange = (value: string, setValue: React.Dispatch<React.SetStateAction<string>>) => {
        setValue(value)
    }

    const toggleCheckbox = (curState: boolean, setState: React.Dispatch<React.SetStateAction<boolean>>) => {
        setState(!curState)
    }

    return (
        <section className="form-section">
            <h2>Оставить заявку</h2>
            <p>После получения заявки, мы свяжемся с Вами для подтверждения по email</p>

            <form>
                <EmailInput
                    value={emailInputValue}
                    onChange={(value: string) => {
                        handleInputChange(value, setEmailInputValue)
                    }}
                />
                <TextField
                    id="name-input"
                    name="name-input"
                    label="Имя"
                    variant="outlined"
                    value={nameInputValue}
                    onChange={(evt: ChangeEvent<HTMLInputElement>) => {
                        handleInputChange(evt.target.value, setNameInputValue)
                    }}
                    fullWidth
                />
                <TextField
                    id="info-input"
                    name="info-input"
                    label="Дополнительная информация"
                    variant="outlined"
                    value={informationInputValue}
                    onChange={(evt: ChangeEvent<HTMLInputElement>) => {
                        handleInputChange(evt.target.value, setInformationInputValue)
                    }}
                    fullWidth
                />

                <div className="form-footer">
                    <FormControlLabel
                        control={
                            <Checkbox checked={isConfirmationCheckboxChecked}/>
                        }
                        label="Я согласен с получением сообщений на почту"
                        onClick={() => {
                            toggleCheckbox(isConfirmationCheckboxChecked, setConfirmationCheckboxState)
                        }}
                    />

                    <Button
                        variant="contained"
                        disabled={!isConfirmationCheckboxChecked}
                        type={"submit"}
                        fullWidth
                    >
                        Подтвердить
                    </Button>
                </div>
            </form>
        </section>
    )
}

export default FormSection