<template>
  <div class="app-container">
    <!-- Left Panel: Sidebar -->
    <div class="sidebar">
      <!-- Login / Logout Section -->
      <input v-model="username" placeholder="Enter your name" />
      <button v-if="!securityKey" @click="loginUser">Login</button>
      <button v-else @click="logoutUser">Logout</button>
      <p v-if="msg">{{ msg }}</p>

      <!-- Change Name Section -->
      <div v-if="securityKey" class="change-name">
        <input v-model="newName" placeholder="Change your name" />
        <button @click="changeName">Change Name</button>
      </div>

      <div class="black-bar"></div>

      <!-- Conversation List Section -->
      <div v-if="securityKey" class="conversation-list">
        <h3>Conversations</h3>
        <div v-if="conversations.length > 0">
          <ul>
            <li
              v-for="conversation in conversations"
              :key="conversation.id"
              @click="selectConversation(conversation)"
              :class="{ selected: selectedConversationDetails && selectedConversationDetails.id === conversation.id }"
            >
              {{ conversation.name }}
            </li>
          </ul>
        </div>
        <div v-else>
          <p>You have no conversations.</p>
        </div>
        <!-- Add New Conversation Button -->
        <button @click="showNewConversation = true">Add New Conversation</button>
      </div>

      <!-- New Conversation Creation Section -->
      <div v-if="securityKey && showNewConversation" class="new-conversation">
        <h3>Start a New Conversation</h3>
        <input v-model="userSearch" placeholder="Search users by name" />
        <ul>
          <li v-for="user in filteredUsers" :key="user.id">
            <label>
              <input type="checkbox" :value="user.id" v-model="selectedUserIds" />
              {{ user.name }}
            </label>
          </li>
        </ul>
        <button @click="createNewConversation">Create Conversation</button>
        <button @click="cancelNewConversation">Cancel</button>
      </div>
    </div>

    <!-- Right Panel: Conversation Details -->
    <div class="conversation-details" v-if="selectedConversationDetails">
      <h1>{{ selectedConversationDetails.name }}</h1>
      <h3>
        Members:
        <span
          v-for="(member, index) in selectedConversationDetails.memberNames"
          :key="member.id"
        >
          {{ member.name }}<span v-if="index < selectedConversationDetails.memberNames.length - 1">, </span>
        </span>
      </h3>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      // User authentication and profile info
      username: "",
      newName: "",
      msg: "",
      securityKey: null,
      userId: null,

      // Data lists
      conversations: [],
      otherUsers: [],

      // New conversation creation UI
      showNewConversation: false,
      userSearch: "",
      selectedUserIds: [],

      // Details for the selected conversation
      selectedConversationDetails: null,
    };
  },
  computed: {
    // Filters the full list of other users by the search term.
    filteredUsers() {
      if (!this.userSearch) return this.otherUsers;
      return this.otherUsers.filter(user =>
        user.name.toLowerCase().includes(this.userSearch.toLowerCase())
      );
    },
  },
  methods: {
    // Log in the user, then fetch conversations and the list of other users.
    async loginUser() {
      try {
        const response = await this.$axios.post("/session", { name: this.username });
        if (!response.data.apiKey || !response.data.userId) {
          throw new Error("Invalid API response");
        }
        this.securityKey = response.data.apiKey;
        this.userId = response.data.userId;
        this.msg = "Logged in successfully";

        // Fetch conversations and user list (for conversation creation & member names)
        await this.fetchConversations();
        await this.fetchUsers();
      } catch (e) {
        this.msg = "Login failed: " + e.message;
      }
    },
    // Logout: clear all data.
    logoutUser() {
      this.username = "";
      this.newName = "";
      this.securityKey = null;
      this.userId = null;
      this.conversations = [];
      this.otherUsers = [];
      this.selectedUserIds = [];
      this.userSearch = "";
      this.showNewConversation = false;
      this.selectedConversationDetails = null;
      this.msg = "Logged out successfully";
    },
    // Change the user’s name.
    async changeName() {
      try {
        if (!this.securityKey || !this.userId || !this.newName) {
          this.msg = "All fields are required to change your name.";
          return;
        }
        const response = await this.$axios.put(
          `/users/${this.userId}/name`,
          { name: this.newName },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json",
            },
          }
        );
        if (response.status === 204) {
          this.msg = "Name changed successfully!";
          this.newName = "";
        } else {
          this.msg = `Unexpected response: ${response.status}`;
        }
      } catch (error) {
        this.msg = "Failed to change name: " + (error.response?.data?.error || error.message);
      }
    },
    // Fetch the conversations list for the logged-in user.
    async fetchConversations() {
      if (!this.securityKey) {
        this.msg = "Authorization key is missing. Please log in.";
        return;
      }
      try {
        const convIdsResponse = await this.$axios.get("/conversations", {
          headers: { Authorization: `Bearer ${this.securityKey}` },
        });
        let conversationIds = convIdsResponse.data.conversations;
        // If it's not an array (for example, null), default to an empty array.
        if (!Array.isArray(conversationIds)) {
          console.log("Received conversation data:", convIdsResponse.data);
          conversationIds = [];
        }
        const convs = [];
        for (const convId of conversationIds) {
          try {
            const convNameResponse = await this.$axios.get(`/conversations/${convId}/name`, {
              headers: { Authorization: `Bearer ${this.securityKey}` },
            });
            if (!convNameResponse.data || !convNameResponse.data.name) {
              throw new Error(`Invalid response for conversation ${convId}`);
            }
            convs.push({ id: convId, name: convNameResponse.data.name });
          } catch (err) {
            console.error(`Failed to fetch conversation ${convId}:`, err);
            convs.push({ id: convId, name: "Unknown" });
          }
        }
        this.conversations = convs;
      } catch (e) {
        this.msg = "Failed to fetch conversations: " + e.message;
      }
    },
    // Fetch the list of all users (excluding the logged-in user) for conversation creation.
    async fetchUsers() {
      if (!this.securityKey) {
        this.msg = "Authorization key is missing. Please log in.";
        return;
      }
      try {
        const maxIdResponse = await this.$axios.get("/users", {
          headers: { Authorization: `Bearer ${this.securityKey}` },
        });
        const maxUserId = maxIdResponse.data.maxUserId;
        if (!maxUserId) {
          throw new Error("Invalid maxUserId received");
        }
        const users = [];
        for (let id = 1; id <= maxUserId; id++) {
          if (id === this.userId) continue;
          try {
            const userResponse = await this.$axios.get(`/users/${id}/name`, {
              headers: { Authorization: `Bearer ${this.securityKey}` },
            });
            if (!userResponse.data || !userResponse.data.name) {
              throw new Error(`Invalid response for user ${id}`);
            }
            users.push({ id, name: userResponse.data.name });
          } catch (err) {
            console.error(`Failed to fetch user ${id}:`, err);
            users.push({ id, name: "Unknown" });
          }
        }
        this.otherUsers = users;
      } catch (e) {
        this.msg = "Failed to fetch user data: " + e.message;
      }
    },
    // Create a new conversation using the selected user IDs (plus the logged-in user’s id).
    async createNewConversation() {
      if (!this.selectedUserIds.length) {
        this.msg = "Please select at least one user.";
        return;
      }
      try {
        // Include the logged-in user's id along with the selected users.
        const conversationUserIds = [this.userId, ...this.selectedUserIds];
        const response = await this.$axios.put(
          "/new_conversation",
          { userIds: conversationUserIds },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json",
            },
          }
        );
        this.msg = "New conversation created successfully!";
        // Clear the new conversation UI.
        this.selectedUserIds = [];
        this.userSearch = "";
        this.showNewConversation = false;
        // Refresh the conversation list to include the new conversation.
        await this.fetchConversations();
      } catch (error) {
        this.msg =
          "Failed to create new conversation: " +
          (error.response?.data?.error || error.message);
      }
    },
    cancelNewConversation() {
      this.showNewConversation = false;
      this.selectedUserIds = [];
      this.userSearch = "";
    },
    // When a conversation is clicked, fetch its details and map participant IDs to names.
    async selectConversation(conversation) {
      try {
        const response = await this.$axios.get(`/conversations/${conversation.id}`, {
          headers: { Authorization: `Bearer ${this.securityKey}` },
        });
        const participants = response.data.participants || [];
        // Map each participant ID to a user object (using otherUsers list or fallback to the logged-in user).
        const memberNames = participants.map((participantId) => {
          let user = this.otherUsers.find((u) => u.id === participantId);
          if (!user && participantId === this.userId) {
            user = { id: participantId, name: this.username };
          }
          if (!user) {
            user = { id: participantId, name: "Unknown" };
          }
          return user;
        });
        this.selectedConversationDetails = {
          id: conversation.id,
          name: conversation.name,
          participants,
          memberNames,
        };
      } catch (error) {
        this.msg =
          "Failed to fetch conversation details: " +
          (error.response?.data?.error || error.message);
      }
    },
  },
};
</script>

<style>
.app-container {
  display: flex;
  min-height: 100vh;
}

/* Left Panel (Sidebar) */
.sidebar {
  width: 300px;
  background-color: #f4f4f4;
  padding: 20px;
  overflow-y: auto;
}
.black-bar {
  width: 100%;
  height: 10px;
  background-color: black;
  margin-top: 10px;
}
.conversation-list {
  margin-top: 15px;
}
.conversation-list ul {
  list-style: none;
  padding: 0;
}
.conversation-list li {
  cursor: pointer;
  padding: 5px;
  border-bottom: 1px solid #ccc;
}
.conversation-list li.selected {
  background-color: #ddd;
}
.new-conversation {
  margin-top: 15px;
  max-height: 200px;
  overflow-y: auto;
}

/* Right Panel (Conversation Details) */
.conversation-details {
  flex: 1;
  padding: 20px;
  margin-left: 60px;
}
</style>

