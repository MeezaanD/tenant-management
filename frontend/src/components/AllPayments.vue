<template>
	<el-container>
		<el-header class="header">
			<h1>Payments</h1>
			<el-button type="primary" class="add-button" @click="openAddDialog">Add Payment</el-button>
		</el-header>

		<el-select v-model="selectedColumns" multiple collapse-tags placeholder="Select Columns" class="column-select"
			@change="saveSelectedColumns">
			<el-option v-for="column in columns" :key="column.prop" :label="column.label" :value="column.prop" />
		</el-select>

		<el-main>
			<el-table :data="payments" class="custom-table">
				<el-table-column v-for="column in visibleColumns" :key="column.prop" :prop="column.prop"
					:label="column.label" :width="column.width" />
				<el-table-column label="Actions" width="200">
					<template v-slot="scope">
						<el-button @click="openEditDialog(scope.row)" type="primary" size="small" class="action-button">
							<el-icon>
								<edit />
							</el-icon>
							Edit
						</el-button>
						<el-button @click="confirmDeletePayment(scope.row.payment_id)" type="danger" size="small"
							class="action-button">
							<el-icon>
								<delete />
							</el-icon>
							Delete
						</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-main>

		<add-payment-for-user v-if="isAddPaymentDialogVisible" :visible="isAddPaymentDialogVisible" :users="users"
			@close="closeAddPaymentDialog" @save="savePayment" />
		<edit-payment v-if="selectedPayment && isEditDialogVisible" :payment="selectedPayment"
			:visible="isEditDialogVisible" @close="closeDialog" @save="savePayment" />
	</el-container>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted } from 'vue';
import { ElMessageBox } from 'element-plus';
import { Edit, Delete } from '@element-plus/icons-vue';
import ApiService from '../api';
import EditPayment from './EditPayment.vue';
import AddPaymentForUser from './AddPaymentForUser.vue';

interface Payment {
	payment_id: number;
	firstname: string;
	agreed_amount: number;
	paid_amount: number;
	remaining_amount: number;
	proof_of_payment: string;
	paid_ontime: boolean;
}

interface User {
	user_id: string;
	firstname: string;
	lastname: string;
}

export default defineComponent({
	components: { EditPayment, AddPaymentForUser },
	setup() {
		const payments = ref<Payment[]>([]);
		const users = ref<User[]>([]);
		const isEditDialogVisible = ref(false);
		const isAddPaymentDialogVisible = ref(false);
		const selectedPayment = ref<Payment | null>(null);

		// Define columns with prop and label
		const columns = ref([
			{ label: 'Payment ID', prop: 'payment_id', width: '120' },
			{ label: 'First Name', prop: 'firstname' },
			{ label: 'Agreed Amount', prop: 'agreed_amount' },
			{ label: 'Paid Amount', prop: 'paid_amount' },
			{ label: 'Remaining Amount', prop: 'remaining_amount' },
			{ label: 'Proof of Payment', prop: 'proof_of_payment' },
			{ label: 'Paid On Time', prop: 'paid_ontime' }
		]);

		// Load selected columns from localStorage or default to all columns
		const selectedColumns = ref<string[]>(JSON.parse(localStorage.getItem('selectedColumns') || JSON.stringify(columns.value.map(col => col.prop))));

		// Compute visible columns based on selectedColumns
		const visibleColumns = computed(() => {
			return columns.value.filter(col => selectedColumns.value.includes(col.prop));
		});

		// Save selected columns to localStorage
		const saveSelectedColumns = () => {
			localStorage.setItem('selectedColumns', JSON.stringify(selectedColumns.value));
		};

		const fetchPayments = async () => {
			try {
				const response = await ApiService.getAllPayments();
				payments.value = response.payments;
			} catch (error) {
				console.error('Error fetching payments:', error);
			}
		};

		const fetchUsers = async () => {
			try {
				const response = await ApiService.getUsers(); // Fetch all users
				users.value = response.user;
			} catch (error) {
				console.error('Error fetching users:', error);
			}
		};

		const openEditDialog = (payment: Payment) => {
			selectedPayment.value = { ...payment };
			isEditDialogVisible.value = true;
		};

		const openAddDialog = () => {
			console.log('Opening Add Payment dialog', isAddPaymentDialogVisible); // Log to confirm button click
			isAddPaymentDialogVisible.value = true;
		};


		const closeDialog = () => {
			isEditDialogVisible.value = false;
		};

		const closeAddPaymentDialog = () => {
			isAddPaymentDialogVisible.value = false;
		};

		const confirmDeletePayment = (paymentID: number) => {
			ElMessageBox.confirm(
				'Are you sure you want to delete this payment?',
				'Warning',
				{
					confirmButtonText: 'Yes, Delete',
					cancelButtonText: 'Cancel',
					type: 'warning',
				}
			)
				.then(() => {
					handleDeletePayment(paymentID);
				})
				.catch(() => {
					console.log('Deletion canceled');
				});
		};

		const handleDeletePayment = async (paymentID: number) => {
			try {
				await ApiService.deletePayment(paymentID.toString());
				fetchPayments();
			} catch (error) {
				console.error('Error deleting payment:', error);
			}
		};

		const savePayment = async () => {
			fetchPayments();
			closeDialog();
			closeAddPaymentDialog();
		};

		onMounted(async () => {
			await fetchPayments();
			await fetchUsers(); // Fetch all users when the component is mounted
		});

		return {
			payments,
			users,
			columns,
			selectedColumns,
			visibleColumns,
			isEditDialogVisible,
			isAddPaymentDialogVisible,
			selectedPayment,
			openEditDialog,
			openAddDialog,
			closeDialog,
			closeAddPaymentDialog,
			confirmDeletePayment,
			handleDeletePayment,
			savePayment,
			saveSelectedColumns,
		};
	},
});
</script>

<style scoped>
.el-main {
	padding: 10px;
	max-height: 80vh;
}
</style>
