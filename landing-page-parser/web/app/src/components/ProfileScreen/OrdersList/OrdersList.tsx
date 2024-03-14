import React, {useState} from "react";
import Order from "./Order";
import './OrdersList.css'
import AgreeAlert from "../../common/AgreeAlert/AgreeAlert";
import {Button} from "@mui/material";

interface OrderDTO {
    id: number,
    userId: number,
    extraInformation: string,
    status: number,
}

interface OrdersListProps {
    orders: Array<OrderDTO>
    onDelete: (orderId: number) => void

    onBack:()=>void
}

const OrdersList: React.FC<OrdersListProps> = (props) => {
    const [isAgreeAlertOpen, setAgreeAlertOpen] = useState(false)
    const [deletingOrderId, setDeletingOrderId] = useState(0)

    return (
        <section className="OrdersList">
            {
                props.orders && props.orders.length > 0 ?
                    props.orders.map((order, index) => {
                        return (
                            <Order
                                order={{
                                    id: order.id,
                                    data: order.extraInformation,
                                    status: order.status,
                                }}
                                key={index}

                                onDelete={(id) => {
                                    setAgreeAlertOpen(true)
                                    setDeletingOrderId(id)
                                }}
                            />
                        )
                    })
                    :
                    <div className="no-orders">
                        <p className="no-orders-title">У вас пока нет заказов. Самое время сделать первый!</p>
                        <Button
                            variant="contained"
                            onClick={props.onBack}
                        >
                            Сделать заказ
                        </Button>
                    </div>
            }

            <AgreeAlert
                title={"Вы уверены, что хотите удалить заказ #" + deletingOrderId + "?"}
                text="Ваш заказ будет удален без возможности восстановления"
                onAgree={() => {
                    props.onDelete(deletingOrderId)
                    setDeletingOrderId(0)
                }}
                onClose={() => {
                    setAgreeAlertOpen(false)
                }}
                open={isAgreeAlertOpen}
            />
        </section>
    )
}

export default OrdersList