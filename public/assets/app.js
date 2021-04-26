var app = new Vue({
    el: '#app',
    data: {
        ws: null,
        serverUrl: "ws://localhost:8080/ws",
        roomInput: null,
        user: {
            name: "alameddin"
        },
        rooms: [
            {"name":  "_notif", "messages": []}
        ],
        users: []
    },
    mounted: function () {
        this.connectToWebsocket();
    },
    methods: {
        connect() {
            this.connectToWebsocket();
        },
        connectToWebsocket() {
            // Pass the name paramter when connecting.
            this.rooms[0].name = this.user.name + "_notif"
            this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
            this.ws.addEventListener('open', (event) => {
                this.onWebsocketOpen(event)
            });
            this.ws.addEventListener('message', (event) => {
                console.log("event")
                console.log(event.data)
            });
        },
        onWebsocketOpen() {
            console.log("connected to WS!");
        },
    }
})