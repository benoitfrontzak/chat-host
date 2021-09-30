class Chat_API {
    async getMessages() {
        const url = '/api/getMessages'
        const response = await fetch(url)
        const data = await response.json()
        return data
    }
}