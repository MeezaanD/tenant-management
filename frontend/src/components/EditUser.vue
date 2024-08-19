<template>
	<el-dialog v-model="visible" title="Edit User" width="50%" @close="close">
		<el-form :model="editableUser" label-position="top">
			<el-form-item label="First Name"><el-input v-model="editableUser.firstname" /></el-form-item>
			<el-form-item label="Last Name"><el-input v-model="editableUser.lastname" /></el-form-item>
			<el-form-item label="Email"><el-input v-model="editableUser.email" /></el-form-item>
			<el-form-item label="Role">
				<el-select v-model="editableUser.userRole" placeholder="Select Role">
					<el-option label="Landlord" value="landlord" />
					<el-option label="Tenant" value="tenant" />
				</el-select>
			</el-form-item>
		</el-form>
		<template #footer>
			<div class="dialog-footer">
				<el-button @click="close">Cancel</el-button>
				<el-button type="primary" @click="save">Save</el-button>
			</div>
		</template>
	</el-dialog>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue';

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
	props: {
		user: {
			type: Object as () => User,
			required: true,
		},
		visible: {
			type: Boolean,
			required: true,
		},
	},
	setup(props, { emit }) {
		const editableUser = ref<User>({ ...props.user });

		watch(
			() => props.user,
			(newVal) => {
				editableUser.value = { ...newVal };
			},
			{ deep: true }
		);

		watch(
			() => props.visible,
			(newVal) => {
				if (!newVal) emit('close');
			}
		);

		const close = () => {
			emit('close');
		};

		const save = () => {
			emit('save', editableUser.value);
		};

		return {
			editableUser,
			close,
			save,
		};
	},
});
</script>

<style scoped>
.dialog-footer {
	text-align: right;
}
</style>
