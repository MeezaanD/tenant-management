<template>
	<el-container>
		<el-header class="header">
			<h1>Your Payments</h1>
			<el-button type="primary" class="add-button" @click="openAddDialog">Add Payment</el-button>
		</el-header>
		<el-main>
			<el-table :data="payments" class="custom-table">
				<el-table-column prop="payment_id" label="Payment ID" width="120" />
				<el-table-column prop="firstname" label="First Name" />
				<el-table-column prop="agreed_amount" label="Agreed Amount" />
				<el-table-column prop="paid_amount" label="Paid Amount" />
				<el-table-column prop="remaining_amount" label="Remaining Amount" />
				<el-table-column prop="proof_of_payment" label="Proof of Payment" />
				<el-table-column prop="paid_ontime" label="Paid On Time">
					<template v-slot="scope">
						<el-tag :type="scope.row.paid_ontime ? 'success' : 'danger'">
							{{ scope.row.paid_ontime ? 'Yes' : 'No' }}
						</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="payment_date" label="Payment Date" />
				<el-table-column label="Actions" width="200">
					<template v-slot="scope">
						<el-button @click="confirmDeletePayment(scope.row.payment_id)" type="danger" size="small"
							class="action-button">Delete</el-button>
					</template>
				</el-table-column>
			</el-table>
		</el-main>

		<add-payment v-if="isAddDialogVisible" :visible="isAddDialogVisible" @close="closeAddDialog"
			@save="savePayment" />
	</el-container>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import { ElMessageBox } from 'element-plus';
import ApiService from '../api';
import AddPayment from './AddPayment.vue';

interface Payment {
	firstname: string;
	payment_id: number;
	agreed_amount: number;
	paid_amount: number;
	remaining_amount: number;
	proof_of_payment: string;
	paid_ontime: boolean;
	payment_date: string;
}

export default defineComponent({
	components: { AddPayment },
	setup() {
		const payments = ref<Payment[]>([]);
		const userID = ref<string | null>(null);
		const isAddDialogVisible = ref(false);

		const fetchPayments = async () => {
			if (!userID.value) return;
			try {
				const paymentsResponse = await ApiService.getPaymentsByUser(userID.value);
				payments.value = paymentsResponse;
			} catch (error) {
				console.error('Error fetching payments:', error);
			}
		};

		const openAddDialog = () => {
			isAddDialogVisible.value = true;
		};

		const closeAddDialog = () => {
			isAddDialogVisible.value = false;
		};

		const savePayment = async () => {
			await fetchPayments(); // Refresh payments after adding
			closeAddDialog();
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

		onMounted(async () => {
			const token = localStorage.getItem('token');
			if (token) {
				try {
					const payload = JSON.parse(atob(token.split('.')[1]));
					userID.value = payload.userID;
					await fetchPayments();
				} catch (error) {
					console.error('Error decoding token:', error);
				}
			}
		});

		return {
			payments,
			isAddDialogVisible,
			openAddDialog,
			closeAddDialog,
			savePayment,
			confirmDeletePayment,
		};
	},
});
</script>

<style scoped>
.el-main {
	padding: 10px;
	max-height: 80vh;
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

.el-table th,
.el-table td {
	border: 1px solid #ebeef5;
}

.el-table th {
	background-color: var(--primary-color);
	color: white;
}

.el-table .el-table__row:hover>td {
	background-color: #f0f8ff;
}
</style>
