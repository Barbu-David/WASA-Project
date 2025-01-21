<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			username: '',
			securityKey: null,
			userId: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async loginUser() {
			try {
				let response = await this.$axios.post("http://localhost:3000/session", {
					name: this.username
				});
				this.securityKey = response.data.apiKey;
				this.userId = response.data.userID;
				alert("Login successful!");
			} catch (e) {
				this.errormsg = "Login failed: " + e;
			}
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<h1 class="h2">Home page</h1>
		<input v-model="username" placeholder="Enter your name" />
		<button @click="loginUser">Login</button>
		<p v-if="securityKey">Security Key: {{ securityKey }}</p>
		<p v-if="userId">User ID: {{ userId }}</p>
		<p v-if="errormsg" style="color:red">{{ errormsg }}</p>
	</div>
</template>

<style>
</style>

