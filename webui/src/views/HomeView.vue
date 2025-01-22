<template>
        <div class="sidebar">
                <input v-model="username" placeholder="Enter your name" />
                <button v-if="!securityKey" @click="loginUser">Login</button>
                <button v-else @click="logoutUser">Logout</button>
                <p v-if="msg">{{ msg }}</p>
                <div class="black-bar"></div>
                <div v-if="securityKey" class="user-list">
                        <h3>Other users</h3>
                        <ul>
                                <li v-for="user in otherUsers" :key="user.id">{{ user.name }}</li>
                        </ul>
                </div>
        </div>
</template>

<script>
export default {
        data: function() {
                return {
                        username: '',
                        msg: '',
                        securityKey: null,
                        userId: null,
                        otherUsers: []
                }
        },
        methods: {
                async loginUser() {
                        try {
                                let response = await this.$axios.post("/session", {
                                        name: this.username
                                });
                                this.securityKey = response.data.apiKey;
                                this.userId = response.data.userId;
                                this.msg="Logged in successfully";
                                await this.fetchUsers();
                        } catch (e) {
                                this.msg = "Login failed: " + e.message;
                        }
                },
                logoutUser() {
                        this.username = '';
                        this.securityKey = null;
                        this.userId = null;
                        this.msg = "Logged out successfully";
                        this.otherUsers = [];
                },
                async fetchUsers() {
                        try {
                                let maxIdResponse = await this.$axios.get("/users", {
                                        headers: { Authorization: `Bearer ${this.securityKey}` }
                                });
                                let maxUserId = maxIdResponse.data.maxUserId;
                                let users = await Promise.all(
                                        Array.from({ length: maxUserId }, (_, i) => i + 1)
                                        .filter(id => id !== this.userId)
                                        .map(async (id) => {
                                                let userResponse = await this.$axios.get(`/users/${id}/name`, {
                                                        headers: { Authorization: `Bearer ${this.securityKey}` }
                                                });
                                                return { id, name: userResponse.data.name };
                                        })
                                );
                                this.otherUsers = users;
                        } catch (e) {
                                this.msg = "Failed to fetch user data: " + e.message;
                        }
                }
        }
}
</script>

<style>
.sidebar {
        width: 200px;
        background-color: #f4f4f4;
        padding: 20px;
}

.black-bar {
        width: 100%;
        height: 10px;
        background-color: black;
        margin-top: 10px;
}

.user-list {
        max-height: 200px;
        overflow-y: auto;
}
</style>

