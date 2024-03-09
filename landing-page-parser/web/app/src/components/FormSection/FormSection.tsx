import {Button, Checkbox, FormControlLabel, TextField} from "@mui/material";
import React, {ChangeEvent, useState} from "react";
import "./FormSection.css"
import EmailInput from "./Inputs/EmailInput";
import OrdersAPI from "../../api/OrdersAPI/OrdersAPI";

interface FormSectionProps {
    ordersAPI: OrdersAPI
    handleError: (message: string) => void
}

const FormSection: React.FC<FormSectionProps> = (props) => {
    const [isDataSubmitting, setDataSubmitting] = useState(false)
    const [emailInputValue, setEmailInputValue] = useState("")
    const [informationInputValue, setInformationInputValue] = useState("")
    const [isConfirmationCheckboxChecked, setConfirmationCheckboxState] = useState(true)

    const handleInputChange = (value: string, setValue: React.Dispatch<React.SetStateAction<string>>) => {
        setValue(value)
    }

    const toggleCheckbox = (curState: boolean, setState: React.Dispatch<React.SetStateAction<boolean>>) => {
        setState(!curState)
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (isDataSubmitting) return

        setDataSubmitting(true)

        try {
            await props.ordersAPI.newOrder(emailInputValue, {
                extraInformation: informationInputValue,
            })
        } catch (e: any) {
            props.handleError(e.toString())
        } finally {
            setDataSubmitting(false)
        }
    }

    return (
        <section className="form-section">
            <h2>Оставить заявку</h2>
            <p>После получения заявки, мы свяжемся с Вами для подтверждения по email</p>

            <form
                onSubmit={handleSubmit}
                className={
                    isDataSubmitting ? "submitting" : ""
                }
            >
                <EmailInput
                    value={emailInputValue}
                    onChange={(value: string) => {
                        handleInputChange(value, setEmailInputValue)
                    }}
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