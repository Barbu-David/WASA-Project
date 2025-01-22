<template>
        <div class="sidebar">
                <input v-model="username" placeholder="Enter your name" />
                <button v-if="!securityKey" @click="loginUser">Login</button>
                <button v-else @click="logoutUser">Logout</button>
                <p v-if="msg">{{ msg }}</p>
                <div class="black-bar"></div>
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
                }
        },
        methods: {
                async loginUser() {
                        try {
                                let response = await this.$axios.post("http://localhost:3000/session", {
                                        name: this.username
                                });
                                this.securityKey = response.data.apiKey;
                                this.userId = response.data.userId;
                                this.msg="Logged in as " +this.username;
                        } catch (e) {
                                this.msg = "Login failed: " + e;
                        }
                },
                logoutUser() {
                        this.username = '';
                        this.securityKey = null;
                        this.userId = null;
                        this.msg = "Logged out successfully";
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
</style>

