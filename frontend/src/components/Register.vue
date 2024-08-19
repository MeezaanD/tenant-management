<template>
	<div class="register-container">
		<h1>Register</h1>
		<el-form @submit.prevent="handleRegister" label-position="top" :model="user">
			<el-form-item label="First Name" :label-width="formLabelWidth">
				<el-input v-model="user.firstname" placeholder="Enter your first name" required />
			</el-form-item>
			<el-form-item label="Last Name" :label-width="formLabelWidth">
				<el-input v-model="user.lastname" placeholder="Enter your last name" required />
			</el-form-item>
			<el-form-item label="Email" :label-width="formLabelWidth">
				<el-input v-model="user.email" type="email" placeholder="Enter your email" required />
			</el-form-item>
			<el-form-item label="Password" :label-width="formLabelWidth">
				<el-input v-model="user.password" type="password" placeholder="Enter your password" required />
			</el-form-item>
			<el-form-item label="User Role" :label-width="formLabelWidth">
				<el-select v-model="user.user_role" placeholder="Select your role" required>
					<el-option label="Landlord" value="landlord"></el-option>
					<el-option label="Tenant" value="tenant"></el-option>
				</el-select>
			</el-form-item>
			<el-button type="primary" @click="handleRegister" class="register-button">Register</el-button>
		</el-form>
		<p>
			Already have an account?
			<router-link to="/" class="link">Login here</router-link>
		</p>
	</div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';
import ApiService from '../api';

export default defineComponent({
	setup() {
		const user = ref({
			firstname: '',
			lastname: '',
			email: '',
			password: '',
			user_role: '',
			joined_date: new Date().toISOString().slice(0, 10), // Default to current date
		});

		const formLabelWidth = '120px';
		const router = useRouter();

		const handleRegister = async () => {
			try {
				await ApiService.register({
					firstname: user.value.firstname,
					lastname: user.value.lastname,
					email: user.value.email,
					password: user.value.password,
					user_role: user.value.user_role,
					joined_date: user.value.joined_date,
				});
				router.push('/user-payments');
			} catch (error) {
				console.error('Error registering:', error);
			}
		};

		return {
			user,
			formLabelWidth,
			handleRegister,
		};
	},
});
</script>

<style scoped>
.register-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	min-height: 100vh;
	background-color: #f9f9f9;
	padding: 40px;
	box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
	border-radius: 10px;
}

.register-container h1 {
	font-family: 'Roboto', sans-serif;
	font-size: 2rem;
	color: var(--primary-color);
	margin-bottom: 20px;
	text-align: center;
}

.el-form {
	width: 100%;
	max-width: 400px;
}

.register-button {
	width: 100%;
	background-color: var(--primary-color);
	border-radius: 25px;
	padding: 12px 20px;
	font-size: 1rem;
	font-weight: bold;
	margin-top: 20px;
}

.register-button:hover {
	background-color: var(--secondary-color);
}

.link {
	color: var(--accent-color);
	text-decoration: none;
	margin-top: 20px;
}

.link:hover {
	text-decoration: underline;
}
</style>
