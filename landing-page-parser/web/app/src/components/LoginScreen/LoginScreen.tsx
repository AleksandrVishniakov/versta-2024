import React, {useState} from "react";
import EmailInput from "../common/Inputs/EmailInput";
import {Button, Checkbox, FormControlLabel} from "@mui/material";
import VerificationCodeInput from "../common/Inputs/VerificationCodeInput";
import './LoginScreen.css'
import authAPI from "../../api/AuthAPI/AuthAPI";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";

interface LoginScreenProps {
    authAPI: authAPI
    handleError: (message: string) => void
    onAuth:(email: string) => void

    onBack:()=>void
}

const LoginScreen: React.FC<LoginScreenProps> = (props) => {
    const [authStep, setAuthStep] = useState(0)
    const [isDataSubmitting, setDataSubmitting] = useState(false)
    const [emailInputValue, setEmailInputValue] = useState("")
    const [verificationCodeInputValue, setVerificationCodeInputValue] = useState("")
    const [isConfirmationCheckboxChecked, setConfirmationCheckboxState] = useState(true)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (isDataSubmitting) return

        if (authStep === 0) {
            await handleAuth()
        } else if (authStep === 1) {
            await handleEmailVerification()
        }
    }

    const handleAuth = async () => {
        setDataSubmitting(true)

        try {
            const userId = await props.authAPI.authUser(emailInputValue)
            console.log("user #" + userId, "successfully authorized")

            setAuthStep(1)
        } catch (e: any) {
            props.handleError(e.toString())
        } finally {
            setDataSubmitting(false)
        }
    }

    const handleEmailVerification = async () => {
        setDataSubmitting(true)

        try {
            await props.authAPI.verifyEmail(
                emailInputValue,
                verificationCodeInputValue
            )

            setAuthStep(0)
            setVerificationCodeInputValue("")

            props.onAuth(emailInputValue)
        } catch (e: any) {
            props.handleError(e.toString())
        } finally {
            setDataSubmitting(false)
        }
    }

    const toggleCheckbox = (curState: boolean, setState: React.Dispatch<React.SetStateAction<boolean>>) => {
        setState(!curState)
    }


    const handleInputChange = (value: string, setValue: React.Dispatch<React.SetStateAction<string>>) => {
        setValue(value)
    }

    return (
        <main className="LoginScreen">
            <Button
                variant="text"
                startIcon={<ArrowBackIcon />}
                onClick={props.onBack}
                style = {{
                    width: "fit-content"
                }}
            >
                На главную
            </Button>

            <h2>Вход в аккаунт</h2>
            <p>Войдите в аккаунт, используя email или создайте новый</p>

            <form
                onSubmit={handleSubmit}
                className={
                    isDataSubmitting ? "partially-transparent" : ""
                }
                id={"auth-form"}
            >
                <EmailInput
                    value={emailInputValue}
                    onChange={(value: string) => {
                        handleInputChange(value, setEmailInputValue)
                    }}

                    className={authStep === 1 ? "partially-transparent" : ""}
                    disabled={authStep === 1}
                />

                <FormControlLabel
                    control={
                        <Checkbox checked={isConfirmationCheckboxChecked}/>
                    }
                    label="Я согласен с получением сообщений на почту"
                    onClick={() => {
                        toggleCheckbox(isConfirmationCheckboxChecked, setConfirmationCheckboxState)
                    }}
                    className={authStep === 1 ? "partially-transparent" : ""}
                    disabled={authStep === 1}
                />

                {
                    authStep === 1 ?
                        <VerificationCodeInput
                            id="verification-input"
                            name="verification-input"
                            label="6-значный код подтверждения"
                            value={verificationCodeInputValue}
                            onChange={(code: string) => {
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
        </main>
    )
}

export default LoginScreen