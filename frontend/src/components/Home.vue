<template>
	<el-container>
		<el-header class="header">
			<h1>Welcome to the Home Page</h1>
		</el-header>

		<el-card class="payment-card" shadow="hover">
			<el-row>
				<el-col :span="24">
					<h2>Payments Still Due</h2>
				</el-col>
				<el-col :span="24">
					<div v-if="pendingPayments.length > 0">
						<ul class="payment-list">
							<li v-for="payment in pendingPayments" :key="payment.payment_id" class="payment-item">
								<div class="payment-name"><strong>Name:</strong> {{ payment.firstname }} {{ payment.lastname }}</div>
								<div class="payment-amount"><strong>Remaining Amount:</strong> R{{ payment.remaining_amount }}</div>
							</li>
						</ul>
					</div>
					<div v-else>
						<p>All payments are settled.</p>
					</div>
				</el-col>
			</el-row>
		</el-card>

		<el-row>
			<el-col :span="6">
				<el-statistic title="Outstanding Payments" :value="outstandingPaymentsCount"
					class="large-title-statistic" />
			</el-col>
		</el-row>
	</el-container>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import ApiService from '../api';

interface Payment {
	payment_id: number;
	firstname: string;
	lastname: string;
	remaining_amount: number;
}

export default defineComponent({
	name: 'Home',
	setup() {
		const pendingPayments = ref<Payment[]>([]);
		const outstandingPaymentsCount = ref<number>(0);

		const fetchPendingPayments = async () => {
			try {
				const response = await ApiService.getAllPayments();
				const payments = response.payments.filter((payment: Payment) => payment.remaining_amount > 0);
				pendingPayments.value = payments;
				outstandingPaymentsCount.value = payments.length;
			} catch (error) {
				console.error('Error fetching pending payments:', error);
			}
		};

		onMounted(fetchPendingPayments);

		return {
			pendingPayments,
			outstandingPaymentsCount,
		};
	},
});
</script>

<style scoped>
.payment-card {
	background-color: #f9f9f9;
	border-radius: 10px;
	box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
	overflow-y: scroll;
	height: 100%;
	max-height: 80vh;
}

.payment-list {
	list-style-type: none;
	padding: 0;
}

.payment-list li {
	margin: 10px 0;
	padding: 10px;
	border: 1px solid #ccc;
	border-radius: 8px;
	background-color: #fff;
	display: flex;
	justify-content: space-between;
}

.payment-name {
	flex: 1;
}

.payment-amount {
	flex: 1;
	text-align: start;
	max-width: 15rem;
	width: 100%;
}

.el-row {
	display: flex;
	justify-content: end;
	margin-top: 20px;
}

.large-title-statistic .el-statistic__title {
	font-size: 1.8rem;
	font-weight: bold;
	color: var(--primary-color);
}

.large-title-statistic .el-statistic__content {
	font-size: 1.5rem;
	color: var(--secondary-color);
}
</style>
