<template>
  <div class="app-container">
    <!-- Left Panel: Sidebar -->
    <div class="sidebar">
      <!-- Login / Logout Section -->
      <input v-model="username" placeholder="Enter your name" />
      <button v-if="!securityKey" @click="loginUser">Login</button>
      <button v-else @click="logoutUser">Logout</button>
      <p v-if="msg">{{ msg }}</p>

    <div v-if="userPhoto" class="user-photo">
      <img :src="userPhoto" alt="Your Profile Photo" />
      <button @click="openFileDialog">Change Photo</button>
   </div>
   <!-- Hidden file input to trigger the file selector -->
   <input
      type="file"
      ref="photoInput"
      accept="image/gif"
      style="display: none"
      @change="handlePhotoChange"
    />
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
            <div class="conversation-item">
              <!-- Display conversation photo if available -->
              <img
                v-if="conversation.photo"
                :src="conversation.photo"
                class="conversation-photo"
                alt="Conversation Photo"
              />
              <div class="conversation-info">
                <div class="conversation-name">{{ conversation.name }}</div>
                <div class="preview">
                  <span>{{ conversation.preview }}</span>
                </div>
              </div>
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

      <!--Conversation Creation Section -->
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
            <div class="conversation-header">
                <img
                    v-if="selectedConversationDetails.photo"
                    :src="selectedConversationDetails.photo"
                    class="conversation-photo-large"
                    alt="Conversation Photo"
                />
                <div class="conversation-header-info">
                    <h1>{{ selectedConversationDetails.name }}</h1>
                    <div class="conversation-preview">
                        <i v-if="selectedConversationDetails.photoPreview" class="photo-icon">[Photo]</i>
                    </div>
                </div>
            </div>
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
                <!-- Option 4: Change Group Photo -->
                <button @click="openGroupPhotoDialog">Change Group Photo</button>
                <!-- Hidden file input to trigger file selector -->
                <input
                    type="file"
                    ref="groupPhotoInput"
                    accept="image/gif"
                    style="display: none"
                    @change="handleGroupPhotoChange"
                />
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


      <!-- Message Sending Section -->
      <div class="message-sending" v-if="securityKey && selectedConversationDetails">
        <!-- Reply Banner: shows when a reply is in progress -->
        <div v-if="replyTo" class="reply-banner">
          Replying to {{ getSenderName(replyTo.senderId) }}: "{{ replyTo.stringContent }}"
          <button @click="cancelReply" style="margin-left: 5px;">Cancel Reply</button>
        </div>
        <input
          v-model="newMessage"
          placeholder="Type your message here"
          @keyup.enter="sendMessage"
        />
        <button @click="sendMessage">Send</button>
      </div>

    <!-- Messages Section -->
    <div class="messages" v-if="securityKey && selectedConversationDetails">
      <h3>Messages</h3>
      <div v-if="reversedMessages.length === 0">No messages in conversation</div>
      <ul v-else>
        <li v-for="(message, index) in reversedMessages" :key="index">
          <!-- Display the sender's photo if available -->
          <img
            v-if="getSenderPhoto(message.senderId)"
            :src="getSenderPhoto(message.senderId)"
            class="sender-photo"
            alt="Sender Photo"
          />
          <!-- Display sender name only for messages from others -->
          <span v-if="message.senderId !== userId">
            <strong>{{ getSenderName(message.senderId) }}</strong>
          </span>
          <!-- For the logged-in user's messages, no sender name is shown -->
          <span v-else></span>
          ({{ message.timestamp ? formatTimestamp(message.timestamp) : 'No timestamp' }}):
          <br />
          <!-- Display message content in blue for the logged-in user's messages, normal otherwise -->
          <span v-if="message.senderId === userId" style="color: blue;">
            {{ message.stringContent || '[No content]' }}
          </span>
          <span v-else>
            {{ message.stringContent || '[No content]' }}
          </span>
          <!-- Display checkmark for the logged-in user's message if available -->
          <span v-if="message.senderId === userId && message.checkmark" style="margin-left: 5px; color: green;">
            {{ message.checkmark }}
          </span>
          <!-- Display forwarded tag if the message has been forwarded -->
          <span v-if="message.forwarded" class="forwarded-tag">Forwarded</span>
          <!-- Action buttons: Delete, Forward, and Reply -->
          <button
            @click="deleteMessage(message)"
            style="background-color: red; color: white; border: none; margin-left: 10px; cursor: pointer;"
          >
            Delete
          </button>
          <button
            @click="openForwardUI(message)"
            style="background-color: green; color: white; border: none; margin-left: 5px; cursor: pointer;"
          >
            Forward
          </button>
          <button
            @click="initiateReply(message)"
            style="background-color: orange; color: white; border: none; margin-left: 5px; cursor: pointer;"
          >
            Reply
          </button>
        </li>
      </ul>
    </div>
    </div>

    <!-- Forward Message Modal -->
    <div v-if="showForwardUI" class="modal">
      <div class="modal-content">
        <h3>Forward Message</h3>
        <p>Select a conversation or user to forward the message:</p>
        <div class="forward-section">
          <h4>Existing Conversations</h4>
          <ul>
            <li v-for="conv in conversations" :key="'conv-' + conv.id">
              <span>{{ conv.name }}</span>
              <button @click="forwardMessageToConversation(conv.id)">Forward</button>
            </li>
          </ul>
        </div>
        <div class="forward-section">
          <h4>Users without Conversation</h4>
          <ul>
            <li v-for="user in nonConversationalUsers" :key="'user-' + user.id">
              <span>{{ user.name }}</span>
              <button @click="forwardMessageToUser(user)">Forward</button>
            </li>
          </ul>
        </div>
        <button @click="cancelForward" class="cancel-button">Cancel</button>
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
      userPhoto: null,

      // Data lists
      conversations: [],
      otherUsers: [],

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
      userRefreshIntervalId: null,

      // Forwarding message UI
      showForwardUI: false,
      messageToForward: null,
      // Store the source conversation id when opening the forward modal
      sourceConversationIdForForward: null,

      replyTo: null,
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
    reversedMessages() {
      return this.messages.slice().reverse();
    },
    nonConversationalUsers() {
      // Compute the list of users with whom there is no one-on-one conversation.
      const oneOnOneUserIds = [];
      this.conversations.forEach(conv => {
        if (!conv.isGroup) {
          const otherId = conv.participants.find(id => id !== this.userId);
          if (otherId) {
            oneOnOneUserIds.push(otherId);
          }
        }
      });
      return this.otherUsers.filter(user => !oneOnOneUserIds.includes(user.id));
    }
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
        this.userId = Number(response.data.userId);
        this.msg = "Logged in successfully";

	await this.fetchMyPhoto()
        await this.fetchUsers();
        await this.fetchConversations();

        this.conversationIntervalId = setInterval(() => {
          this.fetchConversations();
        }, 10000);

     this.userRefreshIntervalId = setInterval(() => {
      this.fetchUsers();
    }, 10000);
      } catch (e) {
        this.msg = "Login failed: " + e.message;
      }
    },

  openFileDialog() {
    this.$refs.photoInput.click();
  },

  async handlePhotoChange(event) {
    const file = event.target.files[0];
    if (!file) return;

    if (file.type !== "image/gif") {
      this.msg = "Please select a valid GIF image.";
      return;
    }

    try {
      await this.$axios.put(`/users/${this.userId}/photo`, file, {
        headers: {
          Authorization: `Bearer ${this.securityKey}`,
          "Content-Type": "image/gif",
        },
      });
      this.msg = "Photo updated successfully!";
      // Refresh the displayed photo by fetching it again.
      await this.fetchMyPhoto();
    } catch (error) {
      console.error("Failed to update photo:", error);
      this.msg =
        "Failed to update photo: " +
        (error.response?.data?.error || error.message);
    }
  },

    async fetchMyPhoto() {
      try {
        const response = await this.$axios.get(`/users/${this.userId}/photo`, {
          headers: { Authorization: `Bearer ${this.securityKey}` },
          responseType: 'blob' // Important: get the binary data
        });
        // Convert the Blob to an object URL that can be used as the image source
        this.userPhoto = URL.createObjectURL(response.data);
      } catch (error) {
        console.error("Failed to fetch my photo:", error);
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
      if (this.userRefreshIntervalId) {
        clearInterval(this.userRefreshIntervalId);
        this.userRefreshIntervalId = null;
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
      this.replyTo = null;
      this.msg = "Logged out successfully";
      this.userPhoto = null;
    },
    // Change the logged-in user's name.
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
                for (const convId of conversationIds) {
                    try {
                        const detailsResponse = await this.$axios.get(`/conversations/${convId}`, {
                            headers: { Authorization: `Bearer ${this.securityKey}` },
                        });
                        const details = detailsResponse.data;
                        const numericParticipants = (details.participants || []).map(p => Number(p));
                        let convName = "";
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
                            const otherId = numericParticipants.find(id => id !== this.userId);
                            const otherUser = this.otherUsers.find(user => user.id === otherId);
                            convName = otherUser ? otherUser.name : "Unknown";
                        }
                        const preview = details.preview || "";
                        const photoPreview = details.photo_preview || false;
                        let convPhoto = null;
                        if (details.is_group) {
                            // For groups, fetch the group photo.
                            try {
                                const photoResponse = await this.$axios.get(`/conversations/${convId}/photo`, {
                                    headers: { Authorization: `Bearer ${this.securityKey}` },
                                    responseType: 'blob'
                                });
                                convPhoto = URL.createObjectURL(photoResponse.data);
                            } catch (err) {
                                console.error(`Failed to fetch group photo for conversation ${convId}:`, err);
                            }
                        } else {
                            // For one-on-one conversations, use the other user's photo.
                            const otherId = numericParticipants.find(id => id !== this.userId);
                            const otherUser = this.otherUsers.find(user => user.id === otherId);
                            convPhoto = otherUser ? otherUser.photo : null;
                        }
                        convs.push({
                            id: convId,
                            name: convName,
                            preview: preview,
                            photoPreview: photoPreview,
                            isGroup: details.is_group,
                            participants: numericParticipants,
                            photo: convPhoto,
                            timestamp: details.timestamp
                        });
                    } catch (err) {
                        console.error(`Failed to fetch conversation ${convId}:`, err);
                        convs.push({
                            id: convId,
                            name: "Unknown",
                            preview: "",
                            photoPreview: false,
                            isGroup: true,
                            participants: [],
                            photo: null,
                            timestamp: null
                        });
                    }
                }
                // Sort conversations so that the latest (by timestamp) appears first
                convs.sort((a, b) => {
                    const timeA = a.timestamp ? new Date(a.timestamp).getTime() : 0;
                    const timeB = b.timestamp ? new Date(b.timestamp).getTime() : 0;
                    return timeB - timeA;
                });
                this.conversations = convs;
            } catch (e) {
                this.msg = "Failed to fetch conversations: " + e.message;
            }
        },
  
 
  openGroupPhotoDialog() {
    this.$refs.groupPhotoInput.click();
  },

  // Handles file selection and uploads the new group photo.
  async handleGroupPhotoChange(event) {
    const file = event.target.files[0];
    if (!file) return;

    if (file.type !== "image/gif") {
      this.msg = "Please select a valid GIF image.";
      return;
    }

    try {
      await this.$axios.put(
        `/conversations/${this.selectedConversationDetails.id}/photo`,
        file,
        {
          headers: {
            Authorization: `Bearer ${this.securityKey}`,
            "Content-Type": "image/gif"
          }
        }
      );
      this.msg = "Group photo updated successfully!";
      // Refresh the conversation photo in both the details and the conversation list.
      await this.fetchConversationPhoto(this.selectedConversationDetails.id);
    } catch (error) {
      console.error("Failed to update group photo:", error);
      this.msg = "Failed to update group photo: " + (error.response?.data?.error || error.message);
    }
  },

  // Fetches and updates the conversation photo for the given conversation ID.
  async fetchConversationPhoto(convId) {
    try {
      const response = await this.$axios.get(`/conversations/${convId}/photo`, {
        headers: { Authorization: `Bearer ${this.securityKey}` },
        responseType: 'blob'
      });
      const newPhoto = URL.createObjectURL(response.data);
      // Update the selected conversation if it matches.
      if (this.selectedConversationDetails && this.selectedConversationDetails.id === convId) {
        this.selectedConversationDetails.photo = newPhoto;
      }
      // Also update the conversation in the list.
      const conv = this.conversations.find(c => c.id === convId);
      if (conv) {
        conv.photo = newPhoto;
      }
    } catch (err) {
      console.error("Failed to fetch conversation photo:", err);
    }
  },

  getSenderPhoto(senderId) {
    if (senderId === this.userId) {
      return this.userPhoto;
    } else {
      const user = this.otherUsers.find(u => u.id === senderId);
      return user ? user.photo : null;
    }
  },
    // Fetch users.
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
            const userObj = { id, name: userResponse.data.name, photo: null };

            // Fetch the user's photo
            try {
              const photoResponse = await this.$axios.get(`/users/${id}/photo`, {
                headers: { Authorization: `Bearer ${this.securityKey}` },
                responseType: 'blob'
              });
              userObj.photo = URL.createObjectURL(photoResponse.data);
            } catch (photoErr) {
              console.error(`Failed to fetch photo for user ${id}:`, photoErr);
            }
            users.push(userObj);
          } catch (e) {
            console.error(`Failed to fetch data for user ${id}:`, e);
            users.push({ id, name: "Unknown", photo: null });
          }
    }
    this.otherUsers = users;
      } catch (e) {
        this.msg = "Failed to fetch user data: " + e.message;
      }
    },


    // Create a new conversation.
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
    // Select a conversation.
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
    // Fetch messages for a conversation.
    async fetchMessages(conversationId, messageIds) {
      try {
        const requests = messageIds.map(messageId =>
          this.$axios.get(`/conversations/${conversationId}/messages/${messageId}`, {
            headers: { Authorization: `Bearer ${this.securityKey}` },
          })
        );
        const responses = await Promise.all(requests);
        this.messages = responses.map((response, index) => {
          return { id: messageIds[index], ...response.data };
        });
      } catch (error) {
        this.msg = "Failed to fetch messages: " +
          (error.response?.data?.error || error.message);
      }
    },
    // Refresh messages.
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
    // Send a new message.
    async sendMessage() {
      if (!this.newMessage.trim()) return;
      // If replying to a message, prepend the reply info.
      let messageContent = this.newMessage.trim();
      if (this.replyTo) {
        messageContent = `Replying to ${this.getSenderName(this.replyTo.senderId)} message: ${this.replyTo.stringContent} : ${messageContent}`;
      }
      try {
        await this.$axios.post(
          `/conversations/${this.selectedConversationDetails.id}`,
          { message: messageContent },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json",
            },
          }
        );
        this.newMessage = "";
        this.replyTo = null; // Clear the reply context after sending.
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
    // Confirm adding members.
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
    // Cancel Add Members UI.
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
    // Confirm changing group name.
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
    // Cancel Change Group Name UI.
    cancelChangeGroupName() {
      this.showChangeGroupName = false;
      this.newGroupName = "";
    },
    // Delete a message.
    async deleteMessage(message) {
      try {
        await this.$axios.delete(
          `/conversations/${this.selectedConversationDetails.id}/messages/${message.id}`,
          {
            headers: { Authorization: `Bearer ${this.securityKey}` },
          }
        );
        this.messages = this.messages.filter(m => m.id !== message.id);
        this.msg = "Message deleted successfully.";
      } catch (error) {
        this.msg =
          "Failed to delete message: " +
          (error.response?.data?.error || error.message);
      }
    },
    // --------------- Forward Message Methods ---------------
    // Open the forward modal and store the source conversation id.
    openForwardUI(message) {
      this.messageToForward = message;
      // Store the conversation id from which the message originates.
      this.sourceConversationIdForForward = this.selectedConversationDetails.id;
      this.showForwardUI = true;
    },
    // Cancel forwarding.
    cancelForward() {
      this.messageToForward = null;
      this.sourceConversationIdForForward = null;
      this.showForwardUI = false;
    },
    // Forward the message to an existing conversation.
    async forwardMessageToConversation(targetConversationId) {
      if (!this.messageToForward) return;
      try {
        await this.$axios.post(
          `/conversations/${this.sourceConversationIdForForward}/messages/${this.messageToForward.id}`,
          { targetConversationId },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json"
            }
          }
        );
        this.msg = "Message forwarded successfully.";
        this.cancelForward();
      } catch (error) {
        this.msg = "Failed to forward message: " + (error.response?.data?.error || error.message);
      }
    },
    // Forward the message by first creating a new conversation with the selected user.
    async forwardMessageToUser(user) {
      if (!this.messageToForward) return;
      try {
        // Create a new conversation with the current user and the selected user.
        const conversationUserIds = [this.userId, user.id];
        await this.$axios.put(
          "/new_conversation",
          { userIds: conversationUserIds },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json"
            }
          }
        );
        await this.fetchConversations();
        // Find the newly created one-on-one conversation.
        const newConversation = this.conversations.find(conv => {
          return (
            !conv.isGroup &&
            conv.participants.includes(user.id) &&
            conv.participants.includes(this.userId)
          );
        });
        if (!newConversation) {
          throw new Error("New conversation not found.");
        }
        await this.$axios.post(
          `/conversations/${this.sourceConversationIdForForward}/messages/${this.messageToForward.id}`,
          { targetConversationId: newConversation.id },
          {
            headers: {
              Authorization: `Bearer ${this.securityKey}`,
              "Content-Type": "application/json"
            }
          }
        );
        this.msg = "Message forwarded successfully to new conversation.";
        this.cancelForward();
      } catch (error) {
        this.msg = "Failed to forward message: " + (error.response?.data?.error || error.message);
      }
    },
    // Format a timestamp.
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleString();
    },
    // Get the sender name.
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
    // --------------- Reply Message Methods ---------------
    // Initiate a reply by storing the message being replied to.
    initiateReply(message) {
      this.replyTo = message;
    },
    // Cancel the reply.
    cancelReply() {
      this.replyTo = null;
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
.conversation-item {
  display: flex;
  align-items: center;
}
.conversation-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}
.conversation-info {
  display: flex;
  flex-direction: column;
}

.conversation-photo-large {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
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

/* Forward Modal Styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 4px;
  width: 90%;
  max-width: 400px;
}
.forward-section {
  margin-bottom: 15px;
}
.forward-section ul {
  list-style: none;
  padding: 0;
}
.forward-section li {
  margin-bottom: 5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.cancel-button {
  background-color: gray;
  color: white;
  border: none;
  padding: 5px 10px;
  cursor: pointer;
}

.user-photo img {
  display: block;
  width: 50px; 
  height: auto;
  margin-bottom: 10px;
  border-radius: 50%;
}

/* Forwarded message tag */
.forwarded-tag {
  font-size: 0.8em;
  padding: 2px 4px;
  border: 1px solid orange;
  border-radius: 3px;
  color: orange;
  margin-left: 5px;
}

.reply-preview {
  background-color: #f0f8ff;
  border-left: 3px solid #2196F3;
  padding: 8px;
  margin-bottom: 8px;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.reply-preview button {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  margin-left: 10px;
  padding: 2px 5px;
}

.reply-preview button:hover {
  color: #333;
}

.sender-photo {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 5px;
  vertical-align: middle;
}



</style>
