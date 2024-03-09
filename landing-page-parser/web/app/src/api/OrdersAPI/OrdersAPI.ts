import {APIError} from "../APIError";

interface OrderDTO {
    extraInformation: string
}

class OrdersAPI {
    private readonly host: string
    constructor(host: string) {
        this.host = host
    }

    public async newOrder(email: string, order: OrderDTO) {
        const response = await fetch(
            this.host + "/api/order?email="+email, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(order)
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }
    }
}

export default OrdersAPI