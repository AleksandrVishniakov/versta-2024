import {APIError} from "../APIError";
import AuthAPI from "../AuthAPI/AuthAPI";
import {Status} from "../Statuses";

interface ChatTokenDTO {
    chatterId: number
    chatToken: string
}

interface MessageDTO {
    id: number
    message: string
    senderId: number
    receiverId: number
    createdAt: Date
}

interface ChatterDTO {
    id: number
    userId: number
    tempSession: string
    unreadMessagesCount: number
}

export class ChatAPI {
    private readonly authAPI: AuthAPI

    public chatterId: number = 0

    private readonly host: string
    private chatToken: string = ""
    private socket: WebSocket|null = null

    private async accessToken(): Promise<string> {
        let accessToken = window.sessionStorage.getItem("accessToken")
        if (!accessToken || accessToken === "") {
            try {
                await this.authAPI.refreshTokens()
            }catch {}

            accessToken = window.sessionStorage.getItem("accessToken")
        }

        return accessToken ? accessToken : ""
    }

    constructor(authAPI: AuthAPI, host: string = 'http://localhost:8003') {
        this.authAPI = authAPI
        this.host = host
    }

    public async sendMessage(msg: string) {
        if (!this.socket || this.socket.readyState !== this.socket.OPEN) return

        this.socket.send(msg)
    }

    public async disconnectChat() {
        if (!this.socket) return

        this.socket.close()
    }

    public async getMessages(withId?: number): Promise<Array<MessageDTO>> {
        const url = !withId?
            this.host + "/api/messages"
            :
            this.host + "/api/admin/messages?with=" + withId

        const accessToken = await this.accessToken()

        const response = await fetch(
            url, {
                method: "GET",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": !accessToken || accessToken === "" ? "" : "Bearer " + accessToken,
                }
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                return await this.getMessages()
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        return await response.json() as Array<MessageDTO>
    }

    public async getUnreadMessagesCount(withID?: number): Promise<number> {
        const url = withID? this.host + "/api/admin/messages/unread?with=" + withID
            :
            this.host + "/api/messages/unread"

        const accessToken = await this.accessToken()

        const response = await fetch(
            url, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": !accessToken || accessToken === "" ? "" : "Bearer " + accessToken,
                },
                credentials: "include"
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                return await this.getUnreadMessagesCount(withID)
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        return await response.json() as number
    }

    public async readAllMessages(withId?: number) {
        const url = !withId?
            this.host + "/api/messages/read_all"
            :
            this.host + "/api/admin/messages/read_all?with=" + withId

        const accessToken = await this.accessToken()

        const response = await fetch(
            url, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": !accessToken || accessToken === "" ? "" : "Bearer " + accessToken,
                },
                credentials: "include"
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                await this.readAllMessages()
                return
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }
    }

    public async connectChat(onMessage: (msg: MessageDTO) => void, withId?: number) {
        if (!window["WebSocket"]) throw new Error("websocket is unsupported")

        if (this.socket) {
            await this.disconnectChat()
        }

        const connect = async () => {
            let accessToken = ""

            if (withId) {
                accessToken = await this.accessToken()
            }

            const url = !withId ?
                "ws://" + this.host.split("://")[1] + "/api/chat?t=" + this.chatToken
                :
                `ws://${this.host.split("://")[1]}/api/admin/chat?t=${this.chatToken}&with=${withId}&jwt=${accessToken}`

            this.socket = new WebSocket(url)

            this.socket.onmessage = (evt) => {
                console.log(evt.data)

                onMessage(JSON.parse(evt.data) as MessageDTO)
            }

            this.socket.onerror = (error) => {
                console.log("ws connection error:", error)

                this.preflightChatRequest().then(connect)
            }
        }

        if (!this.chatToken || this.chatToken === "") {
            await this.preflightChatRequest()
        }

        await connect()
    }

    public async preflightChatRequest() {
        const url = this.host + "/api/chat/preflight"

        const accessToken = window.sessionStorage.getItem("accessToken")

        const response = await fetch(
            url, {
                method: "GET",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": !accessToken || accessToken === "" ? "" : "Bearer " + accessToken,
                }
            }
        )

        if (!response.ok) {
            const apiError = await response.json() as APIError
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                await this.preflightChatRequest()
                return
            }

            throw new Error(apiError.code + ": " + apiError.message)
        }

        const dto = await response.json() as ChatTokenDTO

        this.chatToken = dto.chatToken
        this.chatterId = dto.chatterId
    }

    public async getAllChatters(): Promise<Array<ChatterDTO>> {
        const url = this.host + "/api/admin/clients"

        const accessToken = await this.accessToken()

        const response = await fetch(
            url, {
                method: "GET",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": !accessToken || accessToken === "" ? "" : "Bearer " + accessToken,
                }
            }
        )

        if (!response.ok) {
            if (response.status === Status.Unauthorized) {
                await this.authAPI.refreshTokens()

                return await this.getAllChatters()
            }

            const apiError = await response.json() as APIError

            throw new Error(apiError.code + ": " + apiError.message)
        }

        return await response.json() as Array<ChatterDTO>
    }
}