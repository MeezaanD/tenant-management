<template>
	<div class="menu-container">
		<el-menu :default-active="activeIndex" class="el-menu-vertical" mode="vertical" :collapse="false">
			<template v-if="userRole === 'landlord'">
				<el-menu-item index="1">
					<router-link to="/home">
						<el-icon>
							<house />
						</el-icon>
						<span>Home</span>
					</router-link>
				</el-menu-item>
				<el-menu-item index="2">
					<router-link to="/users">
						<el-icon>
							<user />
						</el-icon>
						<span>Users</span>
					</router-link>
				</el-menu-item>
				<el-menu-item index="3">
					<router-link to="/payments">
						<el-icon><credit-card /></el-icon>
						<span>Payments</span>
					</router-link>
				</el-menu-item>
			</template>

			<el-menu-item index="4">
				<router-link to="/user-payments">
					<el-icon>
						<money />
					</el-icon>
					<span>My Payments</span>
				</router-link>
			</el-menu-item>

			<el-menu-item index="5" v-if="hasToken">
				<router-link to="/profile">
					<el-icon>
						<user />
					</el-icon>
					<span>Profile</span>
				</router-link>
			</el-menu-item>

			<el-menu-item index="6" v-if="hasToken" class="logout-item">
				<span class="logout-button" @click="confirmLogout">
					<el-icon><switch-filled /></el-icon>
					<span>Logout</span>
				</span>
			</el-menu-item>
		</el-menu>

		<div class="user-info">
			<div v-if="hasToken">
				<p><el-icon><user /></el-icon>{{ userEmail }}</p>
			</div>
			<div v-else>
				<el-button type="secondary" @click="redirectToLogin">Login</el-button>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import ApiService from '../api';
import { ElMessageBox } from 'element-plus';
import { House, User, CreditCard, Money, SwitchFilled } from '@element-plus/icons-vue';

const activeIndex = ref('1');
const router = useRouter();

const userEmail = ref<string | null>(null);
const userRole = ref<string | null>(null); 

const hasToken = computed(() => !!localStorage.getItem('token'));

const confirmLogout = () => {
	ElMessageBox.confirm(
		'Are you sure you want to log out?',
		'Confirmation',
		{
			confirmButtonText: 'Yes, Logout',
			cancelButtonText: 'Cancel',
			type: 'warning',
		}
	)
		.then(() => {
			handleLogout();
		})
		.catch(() => {
			console.log('Logout cancelled');
		});
};

const handleLogout = () => {
	if (hasToken.value) {
		ApiService.logout();
		router.push('/');
	}
};

const redirectToLogin = () => {
	router.push('/');
};

if (hasToken.value) {
	const token = localStorage.getItem('token');
	if (token) {
		try {
			const payload = JSON.parse(atob(token.split('.')[1])); // Decode the JWT token to get user info
			userEmail.value = payload.email || 'Unknown Email';
			userRole.value = payload.userRole || 'tenant';
		} catch (error) {
			console.error('Error decoding token:', error);
		}
	}
}
</script>

<style scoped>
.menu-container {
	display: flex;
	flex-direction: column;
	height: 100vh;
	width: 100%;
	border-right: 1px solid #ebeef5;
	overflow: hidden;
	background-color: #f0f2f5;
}

.el-menu-vertical {
	flex: 1;
	width: 100%;
	overflow-y: auto;
	padding: 10px;
}

.el-menu-item {
	margin-bottom: 10px;
	display: flex;
	align-items: center;
	padding-left: 16px;
	border-radius: 8px;
	transition: background-color 0.3s, transform 0.3s;
}

.el-menu-item:hover,
.el-menu-item.is-active {
	background-color: var(--secondary-color);
	transform: translateY(-2px);
}

.el-menu-item a,
.logout-button {
	text-decoration: none;
	color: inherit;
	width: 100%;
	text-align: left;
	display: flex;
	align-items: center;
}

.logout-button {
	cursor: pointer;
}

.user-info {
	padding: 1rem;
	text-align: center;
	font-weight: bold;
	color: #333;
	background-color: #fff;
	border-top: 1px solid #ebeef5;
}

.user-info p {
	margin: 0;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.user-info button {
	width: 100%;
}

.el-icon {
	margin-right: 10px;
}
</style>
