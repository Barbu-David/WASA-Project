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
		}
	}
}
</script>

<template>
	<div>
		<h1 class="h2">Log in</h1>
		<input v-model="username" placeholder="Enter your name" />
		<button @click="loginUser">Login</button>
		<p v-if="msg">{{ msg }}</p>
	</div>
</template>

<style>
</style>
