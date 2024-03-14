import {Button, Checkbox, FormControlLabel, TextField} from "@mui/material";
import React, {ChangeEvent, useEffect, useState} from "react";
import "./FormSection.css"
import EmailInput from "../../common/Inputs/EmailInput";
import OrdersAPI from "../../../api/OrdersAPI/OrdersAPI";
import VerificationCodeInput from "../../common/Inputs/VerificationCodeInput";

interface FormSectionProps {
    ordersAPI: OrdersAPI
    userEmail: string
    handleError: (message: string) => void
    onAuth:(email: string) => void
}

const FormSection: React.FC<FormSectionProps> = (props) => {
    const [isDataSubmitting, setDataSubmitting] = useState(false)
    const [emailInputValue, setEmailInputValue] = useState(props.userEmail)
    const [informationInputValue, setInformationInputValue] = useState("")
    const [verificationCodeInputValue, setVerificationCodeInputValue] = useState("")
    const [isConfirmationCheckboxChecked, setConfirmationCheckboxState] = useState(true)
    const [orderingStep, setOrderingStep] = useState(0)
    const [orderId, setOrderId] = useState(0)

    useEffect(() => {
        setEmailInputValue(props.userEmail)
    }, [props.userEmail]);

    const handleInputChange = (value: string, setValue: React.Dispatch<React.SetStateAction<string>>) => {
        setValue(value)
    }

    const toggleCheckbox = (curState: boolean, setState: React.Dispatch<React.SetStateAction<boolean>>) => {
        setState(!curState)
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (isDataSubmitting) return

        if (orderingStep === 0) {
            await handleOrderRequest()
        } else if (orderingStep === 1) {
            await handleOrderVerification()
        }
    }

    const handleOrderRequest = async () => {
        setDataSubmitting(true)

        try {
            const orderId = await props.ordersAPI.newOrder(emailInputValue, {
                extraInformation: informationInputValue,
            })

            setOrderingStep(1)
            setOrderId((orderId))
        } catch (e: any) {
            props.handleError(e.toString())
        } finally {
            setDataSubmitting(false)
        }
    }

    const handleOrderVerification = async () => {
        setDataSubmitting(true)

        if (orderId <= 0) return

        try {
            await props.ordersAPI.verifyOrder(
                orderId,
                emailInputValue,
                verificationCodeInputValue
            )

            setOrderingStep(0)
            setOrderId(0)
            setVerificationCodeInputValue("")
            setInformationInputValue("")

            props.onAuth(emailInputValue)
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
                    isDataSubmitting ? "partially-transparent" : ""
                }
                id={"order-form"}
            >
                <EmailInput
                    value={emailInputValue}
                    onChange={(value: string) => {
                        handleInputChange(value, setEmailInputValue)
                    }}

                    className={orderingStep === 1 ? "partially-transparent" : ""}
                    disabled={props.userEmail !== "" || orderingStep === 1}
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
                    className={orderingStep === 1 ? "partially-transparent" : ""}
                    disabled={orderingStep === 1}
                />

                <FormControlLabel
                    control={
                        <Checkbox checked={isConfirmationCheckboxChecked}/>
                    }
                    label="Я согласен с получением сообщений на почту"
                    onClick={() => {
                        toggleCheckbox(isConfirmationCheckboxChecked, setConfirmationCheckboxState)
                    }}
                    className={orderingStep === 1 ? "partially-transparent" : ""}
                    disabled={orderingStep === 1}
                />

                {
                    orderingStep === 1 ?
                        <VerificationCodeInput
                            id="verification-input"
                            name="verification-input"
                            label="6-значный код подтверждения"
                            value={verificationCodeInputValue}
                            onChange={(code: string)=>{
                                setVerificationCodeInputValue(code)
                            }}
                        />
                        : null
                }

                <Button
                    variant="contained"
                    disabled={!isConfirmationCheckboxChecked}
                    type={"submit"}
                    fullWidth
                >
                    Подтвердить
                </Button>
            </form>
        </section>
    )
}

export default FormSection