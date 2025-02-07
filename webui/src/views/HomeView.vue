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
        <div v-if="securityKey" class="change-name">
            <input v-model="newName" placeholder="Change your name" />
            <button @click="changeName">Change Name</button>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            username: '',
            msg: '',
            securityKey: null,
            userId: null,
            newName: '', 
            otherUsers: []
        };
    },
    methods: {
        async loginUser() {
            try {
                const response = await this.$axios.post("/session", {
                    name: this.username
                });

                if (!response.data.apiKey || !response.data.userId) {
                    throw new Error("Invalid API response");
                }

                this.securityKey = response.data.apiKey;
                this.userId = response.data.userId;
                this.msg = "Logged in successfully";
                await this.fetchUsers();
            } catch (e) {
                this.msg = "Login failed: " + e.message;
            }
        },

        logoutUser() {
            this.username = '';
            this.securityKey = null;
            this.userId = null;
            this.otherUsers = [];
            this.msg = "Logged out successfully";
        },

        async changeName() {
            try {
                // Ensure required fields are present
                if (!this.securityKey || !this.userId || !this.newName) {
                    this.msg = "All fields are required to change your name.";
                    return;
                }

                // Prepare the API request
                const response = await this.$axios.put(`/users/${this.userId}/name`, 
                    { name: this.newName }, // Request body
                    {
                        headers: {
                            Authorization: `Bearer ${this.securityKey}`, // Security key in header
                            "Content-Type": "application/json"
                        }
                    }
                );

                if (response.status === 204) {
                    this.msg = "Name changed successfully!";
                    this.newName = ""; // Clear the input field after success
                } else {
                    this.msg = `Unexpected response: ${response.status}`;
                }
            } catch (error) {
                // Handle errors gracefully
                this.msg = `Failed to change name: ${error.response?.data?.error || error.message}`;
            }
        },

        async fetchUsers() {
            if (!this.securityKey) {
                this.msg = "Authorization key is missing. Please log in.";
                return;
            }

            try {
                // Get max user ID
                const maxIdResponse = await this.$axios.get("/users", {
                    headers: { Authorization: `Bearer ${this.securityKey}` }
                });

                const maxUserId = maxIdResponse.data.maxUserId;
                if (!maxUserId) {
                    throw new Error("Invalid maxUserId received");
                }

                // Fetch each user's name
                const users = [];
                for (let id = 1; id <= maxUserId; id++) {
                    if (id === this.userId) continue; // Skip self

                    try {
                        const userResponse = await this.$axios.get(`/users/${id}/name`, {
                            headers: { Authorization: `Bearer ${this.securityKey}` }
                        });

                        if (!userResponse.data || !userResponse.data.name) {
                            throw new Error(`Invalid response for user ${id}`);
                        }

                        users.push({ id, name: userResponse.data.name });
                    } catch (err) {
                        console.error(`Failed to fetch user ${id}:`, err);
                        users.push({ id, name: "Unknown" }); // Handle missing users gracefully
                    }
                }

                this.otherUsers = users;
            } catch (e) {
                this.msg = "Failed to fetch user data: " + e.message;
            }
        }
    }
};
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

