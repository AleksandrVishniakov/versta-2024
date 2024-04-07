import React from "react";
import './TitleSection.css'
import UserAccount from "../UserAccount/UserAccount";

interface TitleSectionProps {
    userAccountProps: {
        userEmail: string
        onOpenProfile: () => void,
        onAuthUser: () => void
    }
}

const TitleSection: React.FC<TitleSectionProps> = (props) => {
    return (
        <section className="TitleSection">

            {/* Кнопка входа в аккаунт. Лучше разместить в top-bar'е или в меню */}
            <UserAccount
                userEmail={props.userAccountProps.userEmail}
                onOpenProfile={props.userAccountProps.onOpenProfile}
                onAuthUser={props.userAccountProps.onAuthUser}
            />

            <h1>Versta-2024</h1>
            <h3>TitleSection</h3>
        </section>
    )
}

export default TitleSection
