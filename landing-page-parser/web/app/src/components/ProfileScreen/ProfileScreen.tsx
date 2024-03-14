import React, {useEffect, useState} from "react";
import AuthAPI from "../../api/AuthAPI/AuthAPI";
import ProfileHeader from "./ProfileHeader/ProfileHeader";
import './ProfileScreen.css'
import OrdersAPI from "../../api/OrdersAPI/OrdersAPI";
import OrdersList from "./OrdersList/OrdersList";
import {Button} from "@mui/material";
import ArrowBackIcon from '@mui/icons-material/ArrowBack';

interface ProfileScreenProps {
    authAPI: AuthAPI
    ordersAPI: OrdersAPI

    handleError: (message: string)=>void

    onBack:()=>void
    onLogout:()=>void
}

interface Order {
    id: number,
    userId: number,
    extraInformation: string,
    status: number
}

const ProfileScreen: React.FC<ProfileScreenProps> = (props) => {
    const [userName, setUserName] = useState("")
    const [email, setEmail] = useState("")
    const [createdAt, setCreatedAt] = useState(new Date())
    const [orders, setOrders] = useState(new Array<Order>())

    useEffect(() => {
        getUserProfile().then((user) => {
            if (!user) return

            setEmail(user.email)
            setUserName(user.name)
            setCreatedAt(user.createdAt)
        })

        getOrders().then((orders)=> {
            if (!orders) return

            setOrders(orders)
        })
    }, [props]);

    const getUserProfile = async () => {
        try {
            return await props.authAPI.getProfile()
        }
        catch (e: any) {
            props.handleError(e.toString())
        }
    }

    const getOrders = async () => {
        try {
            return await props.ordersAPI.getAllOrders()
        }
        catch (e: any) {
            props.handleError(e.toString())
        }
    }

    const updateName = async(name: string)=> {
        try {
            return await props.authAPI.updateName(name)
        }
        catch (e: any) {
            props.handleError(e.toString())
        }
    }

    const deleteOrder = async(orderId: number)=> {
        try {
            return await props.ordersAPI.deleteOrder(orderId)
        }
        catch (e: any) {
            props.handleError(e.toString())
        }
    }

    const handleUpdateName = (name: string)=> {
        updateName(name).then(()=>{
            setUserName(name)
        })
    }

    const handleOrderDelete = (orderId: number) => {
        deleteOrder(orderId).then(()=>{
            console.log("order #" + orderId + " successfully deleted")

            getOrders().then((orders)=> {
                if (!orders) return

                setOrders(orders)
            })
        })
    }

    return (
        <main className="ProfileScreen">
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

            <ProfileHeader
                email={email}
                name={userName}
                createdAt={createdAt}
                onUpdateName={handleUpdateName}
                onLogout={props.onLogout}
            />

            <OrdersList
                orders={orders}
                onDelete={handleOrderDelete}
                onBack={props.onBack}
            />
        </main>
    )
}

export default ProfileScreen