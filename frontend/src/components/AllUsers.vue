<template>
	<el-container>
		<el-header class="header">
			<h1>Users</h1>
		</el-header>

		<el-main>
			<el-input
				v-model="searchQuery"
				placeholder="Search by first name"
				class="search-input"
				clearable />
				
			<el-table :data="filteredUsers" class="custom-table">
				<el-table-column prop="user_id" label="User ID" width="120" />
				<el-table-column prop="firstname" label="First Name" />
				<el-table-column prop="lastname" label="Last Name" />
				<el-table-column prop="email" label="Email" />
				<el-table-column prop="user_role" label="Role" />
				<el-table-column prop="joined_date" label="Joined Date" />
				<el-table-column label="Actions" width="200">
					<template v-slot="scope">
						<el-button @click="openEditDialog(scope.row)" type="primary" size="small" class="action-button">Edit</el-button>
						<el-button @click="confirmDeleteUser(scope.row.user_id)" type="danger" size="small" class="action-button">Delete</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-main>

		<edit-user
			v-if="selectedUser && isEditDialogVisible"
			:user="selectedUser"
			:visible="isEditDialogVisible"
			@close="closeDialog"
			@save="saveUser" />
	</el-container>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue';
import { ElMessageBox } from 'element-plus';
import ApiService from '../api';
import EditUser from './EditUser.vue';

interface User {
	user_id: string;
	firstname: string;
	lastname: string;
	email: string;
	userRole: string;
	joined_date: string;
	password?: string;
}

export default defineComponent({
	components: { EditUser },
	setup() {
		const users = ref<User[]>([]);
		const searchQuery = ref('');
		const isEditDialogVisible = ref(false);
		const selectedUser = ref<User | null>(null);

		const fetchUsers = async () => {
			try {
				users.value = (await ApiService.getUsers()).user;
			} catch (error) {
				console.error('Error fetching users:', error);
			}
		};

		const filteredUsers = computed(() => {
			return users.value.filter(user =>
				user.firstname.toLowerCase().includes(searchQuery.value.toLowerCase())
			);
		});

		const openEditDialog = (user: User) => {
			selectedUser.value = { ...user };
			isEditDialogVisible.value = true;
		};

		const closeDialog = () => {
			isEditDialogVisible.value = false;
		};

		const confirmDeleteUser = (userID: string) => {
			ElMessageBox.confirm(
				'Are you sure you want to delete this user?',
				'Warning',
				{
					confirmButtonText: 'Yes, Delete',
					cancelButtonText: 'Cancel',
					type: 'warning',
				}
			)
				.then(() => {
					handleDeleteUser(userID);
				})
				.catch(() => {
					console.log('Deletion canceled');
				});
		};

		const handleDeleteUser = async (userID: string) => {
			try {
				await ApiService.deleteUser(userID);
				fetchUsers();
			} catch (error) {
				console.error('Error deleting user:', error);
			}
		};

		const saveUser = async (updatedUser: User) => {
			try {
				await ApiService.updateUser(updatedUser.user_id, updatedUser);
				fetchUsers();
				closeDialog();
			} catch (error) {
				console.error('Error updating user:', error);
			}
		};

		fetchUsers();

		return {
			users,
			searchQuery,
			filteredUsers,
			isEditDialogVisible,
			selectedUser,
			openEditDialog,
			closeDialog,
			confirmDeleteUser,
			handleDeleteUser,
			saveUser,
		};
	},
});
</script>

<style scoped>
.search-input {
	margin-bottom: 20px;
	border-radius: 25px; 
}

.custom-table {
	width: 100%;
	border-radius: 10px; 
	overflow: hidden;
}

.action-button {
	border-radius: 25px; 
	font-weight: bold;
	transition: background-color 0.3s ease;
}

.action-button:hover {
	background-color: var(--secondary-color); 
}

.el-table th, .el-table td {
	border: 1px solid #ebeef5;
}

.el-table th {
	background-color: var(--primary-color);
	color: white;
}

.el-table .el-table__row:hover > td {
	background-color: #f0f8ff;
}
</style>
