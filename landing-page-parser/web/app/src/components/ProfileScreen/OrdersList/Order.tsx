import React from "react";
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import './Order.css'
import SaveIcon from "@mui/icons-material/Save";
import {IconButton} from "@mui/material";

interface OrderProps {
    id: number,
    data: string,
    status: number
}

const Order: React.FC<{
    order: OrderProps,
    onDelete: (orderId: number)=>void
}> = (props) => {
    const handleOrderDelete = () => {
        props.onDelete(props.order.id)
    }

    return (
        <div
            className="Order"
            key={props.order.id}
            id={"order-" + props.order.id}
            data-order-status={props.order.status}
        >
            <div>
                <IconButton
                    aria-label="delete"
                    size="medium"
                    onClick={handleOrderDelete}
                >
                    <DeleteForeverIcon fontSize="inherit"/>
                </IconButton>
            </div>

            <div className="order-info">
                <h4>{"Заказ #" + props.order.id}</h4>
                <p>{props.order.data}</p>
            </div>
        </div>
    )
}

export default Order