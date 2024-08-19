<template>
	<div class="login-container">
		<h1>Login</h1>
		<div class="form">
			<el-form @submit.prevent="handleLogin">
				<el-form-item>
					<el-input v-model="email" type="email" placeholder="Email" required @keyup.enter="handleLogin" />
				</el-form-item>
				<el-form-item>
					<el-input v-model="password" type="password" placeholder="Password" required
						@keyup.enter="handleLogin" />
				</el-form-item>
				<el-button type="primary" @click="handleLogin" class="login-button">Login</el-button>
			</el-form>
		</div>
		<div class="text">
			<p>
				Don't have an account?
				<router-link to="/register" class="link">Register here</router-link>
			</p>
			<p>
				<router-link to="/home" class="link">Skip login</router-link>
			</p>
		</div>
	</div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';
import ApiService from '../api';

export default defineComponent({
	setup() {
		const email = ref('');
		const password = ref('');
		const router = useRouter();

		const handleLogin = async () => {
			try {
				await ApiService.login(email.value, password.value);
				router.push('/user-payments');
			} catch (error) {
				console.error('Login error:', error);
			}
		};

		return {
			email,
			password,
			handleLogin,
		};
	},
});
</script>

<style scoped>
.login-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	min-height: 100vh;
	background-color: #f9f9f9;
	box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
	border-radius: 10px;
}

.form {
	width: 100%;
	max-width: 20rem;
	margin-bottom: 20px;
}

.login-container h1 {
	font-family: 'Roboto', sans-serif;
	font-size: 2rem;
	color: var(--primary-color);
	margin-bottom: 20px;
	text-align: center;
}

.login-button {
	width: 100%;
	background-color: var(--primary-color);
	border-radius: 25px;
	padding: 1rem;
	font-size: 1rem;
	font-weight: bold;
	transition: background-color 0.3s ease;
}

.login-button:hover {
	background-color: var(--secondary-color);
}

.text {
	text-align: center;
}

.text p {
	margin-top: 10px;
	font-family: 'Open Sans', sans-serif;
	color: var(--primary-color);
}

.link {
	color: var(--accent-color);
	text-decoration: none;
}

.link:hover {
	text-decoration: underline;
}
</style>
