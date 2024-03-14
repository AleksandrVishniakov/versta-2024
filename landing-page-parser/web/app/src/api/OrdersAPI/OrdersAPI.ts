import {APIError} from "../APIError";
import AuthAPI from "../AuthAPI/AuthAPI";
import {Status} from "../Statuses";

interface OrderRequestDTO {
    extraInformation: string
}

interface OrderResponseDTO {
    id: number,
    userId: number,
    extraInformation: string,
    status: number
}

class OrdersAPI {
    private readonly host: string
    private readonly authAPI: AuthAPI
    constructor(
        host: string,
        authAPI: AuthAPI
    ) {
        this.host = host
        this.authAPI = authAPI
    }

    public async newOrder(email: string, order: OrderRequestDTO): Promise<number> {
        const response = await fetch(
            this.host + "/api/order?email="+email, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(order),
                credentials: "include"
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        return await response.json() as number
    }

    public async verifyOrder(orderId: number, email: string, verificationCode: string) {
        const url = this.host + "/api/order/" + orderId + "/verify?email=" + email + "&code=" + verificationCode

        const response = await fetch (
            url, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include"
            }
        )

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        const accessToken = await response.json() as string

        window.sessionStorage.setItem("accessToken", accessToken)
    }

    public async getAllOrders():Promise<Array<OrderResponseDTO>> {
        const url = this.host + "/api/orders"

        const response = await fetch (
            url, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + window.sessionStorage.getItem("accessToken"),
                },
                credentials: "include"
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                return await this.getAllOrders()
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        return await response.json() as Array<OrderResponseDTO>
    }

    public async deleteOrder(orderId: number): Promise<void> {
        const url = this.host + "/api/order/" + orderId

        const response = await fetch (
            url, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + window.sessionStorage.getItem("accessToken"),
                },
                credentials: "include"
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                return await this.deleteOrder(orderId)
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }
    }
}

export default OrdersAPI