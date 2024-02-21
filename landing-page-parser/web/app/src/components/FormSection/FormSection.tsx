import React, { useState } from "react";
import "./FormSection.css"
import EmailInput from "./Inputs/EmailInput";
const FormSection: React.FC = (props) => {
    const [emailInputValue, setEmailIputValue] = useState("")
    
    const handleInputChange = (value: string, setValue: React.Dispatch<React.SetStateAction<string>>) => {
        setValue(value)
        console.log(value)
    }
    
    return (
        <section className="form-section">
            <h2>Оставить заявку</h2>
            <p>После получения заявки, мы свяжемся с Вами для подтверждения по email</p>
            
            <form>
                <EmailInput
                    value={emailInputValue}
                    onChange={(value: string) => {
                        handleInputChange(value, setEmailIputValue)
                    }}
                />
            </form>
        </section>
    )
}

export default FormSection