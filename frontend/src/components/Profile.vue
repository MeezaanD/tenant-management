<template>
	<div class="profile-container">
		<h1>User Profile</h1>
		<el-card v-if="user" class="profile-card">
			<el-form :model="user" label-position="top" class="profile-form">
				<el-form-item label="Email">
					<el-input v-model="user.email" disabled />
				</el-form-item>
				<el-form-item label="First Name">
					<el-input v-model="user.firstname" />
				</el-form-item>
				<el-form-item label="Last Name">
					<el-input v-model="user.lastname" />
				</el-form-item>
				<el-form-item label="Joined Date">
					<el-input v-model="user.joined_date" disabled />
				</el-form-item>
				<el-form-item class="form-actions">
					<el-button type="primary" @click="saveUser" class="action-button">Save</el-button>
					<el-button @click="resetForm" class="action-button">Cancel</el-button>
				</el-form-item>
			</el-form>
		</el-card>
		<el-alert v-else title="Loading user information..." type="info"></el-alert>
	</div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import ApiService from '../api';

interface User {
	user_id: string;
	firstname: string;
	lastname: string;
	email: string;
	userRole: string;
	joined_date: string;
}

export default defineComponent({
	setup() {
		const user = ref<User | null>(null);
		const originalUser = ref<User | null>(null);

		onMounted(async () => {
			const token = localStorage.getItem('token');
			if (token) {
				try {
					const payload = JSON.parse(atob(token.split('.')[1]));
					const response = await ApiService.getUser(payload.userID);
					user.value = response.user;
					originalUser.value = { ...response.user };
				} catch (error) {
					console.error('Error fetching user details:', error);
				}
			}
		});

		const saveUser = async () => {
			if (user.value) {
				try {
					await ApiService.updateUser(user.value.user_id, user.value);
					originalUser.value = { ...user.value };
				} catch (error) {
					console.error('Error updating user:', error);
				}
			}
		};

		const resetForm = () => {
			if (originalUser.value) {
				user.value = { ...originalUser.value };
			}
		};

		return {
			user,
			saveUser,
			resetForm,
		};
	},
});
</script>

<style scoped>
.profile-container {
	max-width: 600px;
	margin: 0 auto;
	padding: 20px;
	background-color: #fff;
	border-radius: 10px;
	box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

h1 {
	text-align: center;
	color: var(--primary-color);
	font-family: 'Roboto', sans-serif;
	margin-bottom: 20px;
}

.profile-card {
	background-color: #f9f9f9;
	border-radius: 10px;
	box-shadow: none;
	padding: 20px;
}

.profile-form {
	width: 100%;
}

.el-form-item {
	margin-bottom: 15px;
}

.form-actions {
	display: flex;
	justify-content: space-between;
}

.action-button {
	background-color: var(--primary-color);
	border-radius: 25px;
	padding: 0.75rem 1.5rem;
	font-weight: bold;
	transition: background-color 0.3s ease;
}

.action-button:hover {
	background-color: var(--secondary-color);
}
</style>
