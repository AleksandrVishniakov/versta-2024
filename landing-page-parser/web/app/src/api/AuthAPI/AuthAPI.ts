import {APIError} from "../APIError";
import {Status} from "../Statuses";

interface User {
    id: number
    email: string
    name: string
    status: string
    isEmailVerified: boolean,
    createdAt: Date
}

class AuthAPI {
    private readonly host: string

    constructor(host: string = "http://localhost:8001") {
        this.host = host
    }

    public async getProfile(email?: string): Promise<User> {

        if (!window.sessionStorage.getItem("accessToken") || window.sessionStorage.getItem("accessToken") === "") {
            await this.refreshTokens()
        }

        const url = email ?
            this.host + "/api/user/email/" + email
            :
            this.host + "/api/user/my_profile"

        const response = await fetch(
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
                await this.refreshTokens()

                return await this.getProfile(email)
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        return await response.json() as User
    }

    public async updateName(name: string): Promise<void> {
        const url = this.host + "/api/user/name"

        const response = await fetch(
            url, {
                method: "PUT",
                body: JSON.stringify({
                    name: name
                }),
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + window.sessionStorage.getItem("accessToken"),
                },
                credentials: "include"
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.refreshTokens()

                return await this.updateName(name)
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }
    }

    public async refreshTokens(): Promise<void> {
        console.log("refreshing tokens...")

        const url = this.host + "/api/tokens/refresh"

        const response = await fetch(
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

    public async authUser(email: string): Promise<number> {
        const url = this.host + "/api/auth?email=" + email

        const response = await fetch(
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

        return await response.json() as number
    }

    public async verifyEmail(email: string, code: string): Promise<void> {
        const url = this.host + "/api/" + email + "/verify?code=" + code

        const response = await fetch(
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
}

export default AuthAPI