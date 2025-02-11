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
              <!-- Display the computed conversation name and the preview message -->
              <div>{{ conversation.name }}</div>
              <div class="preview">
                <span>{{ conversation.preview }}</span>
              </div>
            </li>
          </ul>
        </div>
        <div v-else>
          <p>You have no conversations.</p>
        </div>
        <!-- Only show Add New Conversation button when no conversation is selected -->
        <button @click="showNewConversation = true">
          Add New Conversation
        </button>
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

    <!-- Right Panel: Conversation Details with Group Options -->
    <div class="conversation-details" v-if="selectedConversationDetails">
      <h1>{{ selectedConversationDetails.name }}</h1>
      <h3 v-if="selectedConversationDetails.isGroup">
        Members:
        <span
          v-for="(member, index) in selectedConversationDetails.memberNames"
          :key="member.id"
        >
          {{ member.name }}<span v-if="index < selectedConversationDetails.memberNames.length - 1">, </span>
        </span>
      </h3>

      <!-- Group Options: shown only if this is a group conversation -->
      <div class="group-options" v-if="selectedConversationDetails.isGroup">
        <h3>Options</h3>
        <!-- Option 1: Leave Group (red button) -->
        <button class="leave-group" style="background-color: red; color: white;" @click="leaveGroup">
          Leave Group
        </button>
        <!-- Option 2: Add Members -->
        <button @click="toggleAddMembers">Add Members</button>
        <!-- Option 3: Change Group Name -->
        <button @click="toggleChangeGroupName">Change Group Name</button>

        <!-- UI for Adding New Members -->
        <div v-if="showAddMembers">
          <h4>Add New Members</h4>
          <input v-model="newMembersSearch" placeholder="Search users" />
          <ul>
            <li v-for="user in filteredUsersForGroup" :key="user.id">
              <label>
                <input type="checkbox" :value="user.id" v-model="selectedNewMemberIds" />
                {{ user.name }}
              </label>
            </li>
          </ul>
          <button @click="confirmAddMembers">Confirm Add Members</button>
          <button @click="cancelAddMembers">Cancel</button>
        </div>

        <!-- UI for Changing Group Name -->
        <div v-if="showChangeGroupName">
          <h4>Change Group Name</h4>
          <input v-model="newGroupName" placeholder="Enter new group name" />
          <button @click="confirmChangeGroupName">Confirm Change</button>
          <button @click="cancelChangeGroupName">Cancel</button>
        </div>
      </div>

      <!-- Messages Section -->
<!-- Messages Section -->
    <div class="messages">
      <h3>Messages</h3>
      <div v-if="messages.length === 0">No messages in conversation</div>
      <ul v-else>
        <li v-for="(message, index) in messages" :key="index">
      <!-- Display sender name differently based on whether it's the logged-in user -->
          <span v-if="message.senderId === userId">
            <strong style="color: blue;">{{ getSenderName(message.senderId) }}</strong>
          </span>
          <span v-else>
            <strong>{{ getSenderName(message.senderId) }}</strong>
          </span>
          ({{ message.timestamp ? formatTimestamp(message.timestamp) : 'No timestamp' }}):
          <br />
          <span>{{ message.stringContent || '[No content]' }}</span>
          <span v-if="message.senderId === userId && message.checkmark" 
                style="margin-left: 5px; color: green;">
            {{ message.checkmark }}
          </span>
          <button
            @click="deleteMessage(message)"
            style="background-color: red; color: white; border: none; margin-left: 10px; cursor: pointer;"
          >
            Delete
          </button>
        </li>
      </ul>
    </div>

      <!-- Message Sending Section -->
      <div class="message-sending">
        <input
          v-model="newMessage"
          placeholder="Type your message here"
          @keyup.enter="sendMessage"
        />
        <button @click="sendMessage">Send</button>
      </div>
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

      // For group options:
      showAddMembers: false,
      newMembersSearch: "",
      selectedNewMemberIds: [],
      showChangeGroupName: false,
      newGroupName: "",
      // Messages for the selected conversation
      messages: [],

      // New message text for sending a message
      newMessage: "",

      // Intervals for auto-refresh
      conversationIntervalId: null,
      messageIntervalId: null,
    };
  },
  computed: {
    filteredUsers() {
      if (!this.userSearch) return this.otherUsers;
      return this.otherUsers.filter(user =>
        user.name.toLowerCase().includes(this.userSearch.toLowerCase())
      );
    },
    filteredUsersForGroup() {
      if (!this.selectedConversationDetails) return [];
      const currentMembers = this.selectedConversationDetails.participants || [];
      if (!this.newMembersSearch) {
        return this.otherUsers.filter(user => !currentMembers.includes(user.id));
      }
      return this.otherUsers.filter(user =>
        !currentMembers.includes(user.id) &&
        user.name.toLowerCase().includes(this.newMembersSearch.toLowerCase())
      );
    },
  },
  methods: {
    // Log in the user and start conversation auto-refresh.
    async loginUser() {
      try {
        const response = await this.$axios.post("/session", { name: this.username });
        if (!response.data.apiKey || !response.data.userId) {
          throw new Error("Invalid API response");
        }
        this.securityKey = response.data.apiKey;
        // Convert userId to a number for correct comparisons.
        this.userId = Number(response.data.userId);
        this.msg = "Logged in successfully";

        await this.fetchUsers();
        await this.fetchConversations();

        // Auto-refresh conversations every 10 seconds.
        this.conversationIntervalId = setInterval(() => {
          this.fetchConversations();
        }, 10000);
      } catch (e) {
        this.msg = "Login failed: " + e.message;
      }
    },
    // Clear all data and intervals.
    logoutUser() {
      if (this.conversationIntervalId) {
        clearInterval(this.conversationIntervalId);
        this.conversationIntervalId = null;
      }
      if (this.messageIntervalId) {
        clearInterval(this.messageIntervalId);
        this.messageIntervalId = null;
      }
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
      this.showAddMembers = false;
      this.newMembersSearch = "";
      this.selectedNewMemberIds = [];
      this.showChangeGroupName = false;
      this.newGroupName = "";
      this.messages = [];
      this.newMessage = "";
      this.msg = "Logged out successfully";
    },
    // Change the logged-in user's name and update the current conversation if needed.
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
          this.username = this.newName;
          this.newName = "";
          if (this.selectedConversationDetails) {
            this.selectedConversationDetails.memberNames = this.selectedConversationDetails.memberNames.map(member => {
              if (member.id === this.userId) {
                return { ...member, name: this.username };
              }
              return member;
            });
          }
        } else {
          this.msg = `Unexpected response: ${response.status}`;
        }
      } catch (error) {
        this.msg = "Failed to change name: " + (error.response?.data?.error || error.message);
      }
    },
    // Fetch conversations by first getting the conversation details (for the preview string) and then computing the name.
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
        if (!Array.isArray(conversationIds)) {
          console.log("Received conversation data:", convIdsResponse.data);
          conversationIds = [];
        }
        const convs = [];
        // For each conversation id, get the conversation details (which include preview, participants, is_group)
        for (const convId of conversationIds) {
          try {
            const detailsResponse = await this.$axios.get(`/conversations/${convId}`, {
              headers: { Authorization: `Bearer ${this.securityKey}` },
            });
            const details = detailsResponse.data;
            const numericParticipants = (details.participants || []).map(p => Number(p));
            let convName = "";
            // If this is a group conversation, try to fetch its name from the dedicated endpoint.
            if (details.is_group) {
              try {
                const nameResponse = await this.$axios.get(`/conversations/${convId}/name`, {
                  headers: { Authorization: `Bearer ${this.securityKey}` },
                });
                convName = nameResponse.data.name || "Group Conversation";
              } catch (err) {
                convName = "Group Conversation";
              }
            } else {
              // For one-on-one conversations, compute the name using the other participant.
              const otherId = numericParticipants.find(id => id !== this.userId);
              const otherUser = this.otherUsers.find(user => user.id === otherId);
              convName = otherUser ? otherUser.name : "Unknown";
            }
            // Use the preview string (latest message) from the conversation details.
            const preview = details.preview || "";
            // Store photo_preview as received (ignored in the UI for now).
            const photoPreview = details.photo_preview || false;
            convs.push({
              id: convId,
              name: convName,
              preview: preview,
              photoPreview: photoPreview,
              isGroup: details.is_group,
              participants: numericParticipants,
            });
          } catch (err) {
            console.error(`Failed to fetch conversation ${convId}:`, err);
            convs.push({ id: convId, name: "Unknown", preview: "", photoPreview: false, isGroup: true, participants: [] });
          }
        }
        this.conversations = convs;
      } catch (e) {
        this.msg = "Failed to fetch conversations: " + e.message;
      }
    },
    // Fetch users for conversation creation.
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
    // Create a new conversation with the selected users.
    async createNewConversation() {
      if (!this.selectedUserIds.length) {
        this.msg = "Please select at least one user.";
        return;
      }
      try {
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
        this.selectedUserIds = [];
        this.userSearch = "";
        this.showNewConversation = false;
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
    // Select a conversation, convert participant IDs to numbers, and start auto-refreshing messages.
    async selectConversation(conversation) {
      try {
        const response = await this.$axios.get(`/conversations/${conversation.id}`, {
          headers: { Authorization: `Bearer ${this.securityKey}` },
        });
        const details = response.data;
        const numericParticipants = (details.participants || []).map(p => Number(p));
        const memberNames = numericParticipants.map((participantId) => {
          let user = this.otherUsers.find(u => u.id === participantId);
          if (!user && participantId === this.userId) {
            user = { id: participantId, name: this.username };
          }
          if (!user) {
            user = { id: participantId, name: "Unknown" };
          }
          return user;
        });
        let convName = "";
        if (details.is_group) {
          try {
            const nameResponse = await this.$axios.get(`/conversations/${conversation.id}/name`, {
              headers: { Authorization: `Bearer ${this.securityKey}` },
            });
            convName = nameResponse.data.name || "Group Conversation";
          } catch (err) {
            convName = "Group Conversation";
          }
        } else {
          const otherId = numericParticipants.find(id => id !== this.userId);
          const otherUser = this.otherUsers.find(user => user.id === otherId);
          convName = otherUser ? otherUser.name : "Unknown";
        }
        this.selectedConversationDetails = {
          id: conversation.id,
          name: convName,
          participants: numericParticipants,
          memberNames,
          isGroup: details.is_group,
          preview: details.preview,
          photoPreview: details.photo_preview,
        };

        const messageIds = details.messages || [];
        if (messageIds.length) {
          await this.fetchMessages(conversation.id, messageIds);
        } else {
          this.messages = [];
        }

        // Clear any previous message refresh interval and start a new one.
        if (this.messageIntervalId) {
          clearInterval(this.messageIntervalId);
        }
        this.messageIntervalId = setInterval(() => {
          this.refreshMessages();
        }, 5000);
      } catch (error) {
        this.msg =
          "Failed to fetch conversation details: " +
          (error.response?.data?.error || error.message);
      }
    },
    // Fetch messages for a given conversation.
    async fetchMessages(conversationId, messageIds) {
      try {
        const requests = messageIds.map(messageId =>
          this.$axios.get(`/conversations/${conversationId}/messages/${messageId}`, {
            headers: { Authorization: `Bearer ${this.securityKey}` },
          })
       );
        const responses = await Promise.all(requests);
        // Attach the corresponding message id to each response object.
        this.messages = responses.map((response, index) => {
          return { id: messageIds[index], ...response.data };
        });
      } catch (error) {
        this.msg = "Failed to fetch messages: " +
          (error.response?.data?.error || error.message);
      }
    },


    // Refresh messages for the currently selected conversation.
    async refreshMessages() {
      if (!this.selectedConversationDetails) return;
      try {
        const response = await this.$axios.get(`/conversations/${this.selectedConversationDetails.id}`, {
          headers: { Authorization: `Bearer ${this.securityKey}` },
        });
        const messageIds = response.data.messages || [];
        if (messageIds.length) {
          await this.fetchMessages(this.selectedConversationDetails.id, messageIds);
        } else {
          this.messages = [];
        }
      } catch (error) {
        console.error("Failed to refresh messages: ", error);
      }
    },
    // Send a new message and then refresh the conversation.
    async sendMessage() {
      if (!this.newMessage.trim()) return;
      try {
        await this.$axios.post(
          `/conversations/${this.selectedConversationDetails.id}`,
          { message: this.newMessage },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json",
            },
          }
        );
        this.newMessage = "";
        // Refresh the conversation (and messages)
        await this.selectConversation({
          id: this.selectedConversationDetails.id,
          name: this.selectedConversationDetails.name,
        });
      } catch (error) {
        this.msg =
          "Failed to send message: " +
          (error.response?.data?.error || error.message);
      }
    },
    // Leave the group conversation.
    async leaveGroup() {
      try {
        await this.$axios.delete(`/conversations/${this.selectedConversationDetails.id}/members`, {
          headers: { Authorization: `Bearer ${this.securityKey}` },
        });
        this.msg = "Successfully left the group!";
        await this.fetchConversations();
        this.selectedConversationDetails = null;
        this.messages = [];
      } catch (error) {
        if (error.response && error.response.status === 403) {
          this.msg = "Not a group";
        } else {
          this.msg = "Failed to leave group: " + (error.response?.data?.error || error.message);
        }
      }
    },
    // Toggle the Add Members UI.
    toggleAddMembers() {
      this.showAddMembers = !this.showAddMembers;
      this.newMembersSearch = "";
      this.selectedNewMemberIds = [];
    },
    // Confirm and add new members to a group.
    async confirmAddMembers() {
      if (!this.selectedNewMemberIds.length) {
        this.msg = "Please select at least one user to add.";
        return;
      }
      try {
        await this.$axios.put(
          `/conversations/${this.selectedConversationDetails.id}/members`,
          { userIds: this.selectedNewMemberIds },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json",
            },
          }
        );
        this.msg = "Members added successfully!";
        await this.selectConversation({
          id: this.selectedConversationDetails.id,
          name: this.selectedConversationDetails.name,
        });
        this.cancelAddMembers();
      } catch (error) {
        if (error.response && error.response.status === 403) {
          this.msg = "Not a group";
        } else {
          this.msg = "Failed to add members: " + (error.response?.data?.error || error.message);
        }
      }
    },
    // Cancel the Add Members UI.
    cancelAddMembers() {
      this.showAddMembers = false;
      this.newMembersSearch = "";
      this.selectedNewMemberIds = [];
    },
    // Toggle the Change Group Name UI.
    toggleChangeGroupName() {
      this.showChangeGroupName = !this.showChangeGroupName;
      this.newGroupName = "";
    },
    // Confirm changing the group name.
    async confirmChangeGroupName() {
      if (!this.newGroupName) {
        this.msg = "New group name is required.";
        return;
      }
      try {
        await this.$axios.put(
          `/conversations/${this.selectedConversationDetails.id}/name`,
          { name: this.newGroupName },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json",
            },
          }
        );
        this.msg = "Group name changed successfully!";
        this.selectedConversationDetails.name = this.newGroupName;
        this.cancelChangeGroupName();
        await this.fetchConversations();
      } catch (error) {
        if (error.response && error.response.status === 403) {
          this.msg = "Not a group";
        } else {
          this.msg = "Failed to change group name: " + (error.response?.data?.error || error.message);
        }
      }
    },
    // Cancel the Change Group Name UI.
    cancelChangeGroupName() {
      this.showChangeGroupName = false;
      this.newGroupName = "";
    },
    

        // Delete a message that was sent by the logged-in user.

	// Delete a message that was sent by the logged-in user.
    async deleteMessage(message) {
      try {
        await this.$axios.delete(
          `/conversations/${this.selectedConversationDetails.id}/messages/${message.id}`,
          {
            headers: { Authorization: `Bearer ${this.securityKey}` },
          }
        );
    
        // Remove the deleted message from the local messages array.
        this.messages = this.messages.filter(m => m.id !== message.id);
        this.msg = "Message deleted successfully.";
      } catch (error) {
        this.msg =
          "Failed to delete message: " +
          (error.response?.data?.error || error.message);
      }
    },


   // Format a timestamp for display.
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleString();
    },
    // Get the sender name from a sender ID.
    getSenderName(senderId) {
      if (senderId === this.userId) {
        return this.username;
      }
      if (this.selectedConversationDetails && this.selectedConversationDetails.memberNames) {
        const member = this.selectedConversationDetails.memberNames.find(m => m.id === senderId);
        if (member) return member.name;
      }
      const foundUser = this.otherUsers.find(user => user.id === senderId);
      return foundUser ? foundUser.name : "Unknown";
    },
  },
  beforeDestroy() {
    if (this.conversationIntervalId) {
      clearInterval(this.conversationIntervalId);
    }
    if (this.messageIntervalId) {
      clearInterval(this.messageIntervalId);
    }
  },
};
</script>

<style>
/* Sidebar Styles */
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
.preview {
  font-size: 0.8em;
  color: #555;
}

/* Conversation Details Styles */
.conversation-details {
  flex: 1;
  padding: 20px;
  margin-left: 60px;
}
.group-options {
  margin-top: 20px;
  padding-top: 10px;
  border-top: 1px solid #ccc;
}
.group-options h3 {
  margin-bottom: 10px;
}
.group-options button {
  margin-right: 5px;
  margin-bottom: 5px;
  padding: 5px 10px;
  cursor: pointer;
}
.leave-group {
  background-color: red;
  color: white;
  border: none;
}

/* Messages Styles */
.messages {
  margin-top: 20px;
}
.messages ul {
  list-style: none;
  padding: 0;
}
.messages li {
  padding: 5px 0;
  border-bottom: 1px solid #ccc;
}

/* Message Sending Styles */
.message-sending {
  margin-top: 20px;
  display: flex;
  align-items: center;
}
.message-sending input {
  flex: 1;
  padding: 8px;
  margin-right: 10px;
}
.message-sending button {
  padding: 8px 16px;
}
</style>

