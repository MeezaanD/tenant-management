<template>
	<el-dialog v-model="visible" title="Add Payment" @close="close">
		<el-form :model="paymentData" label-position="top">
			<el-form-item label="Select User">
				<el-select v-model="paymentData.user_id" placeholder="Select User">
					<el-option v-for="user in users" :key="user.user_id" :label="`${user.firstname} ${user.lastname}`" :value="user.user_id" />
				</el-select>
			</el-form-item>
			<el-form-item label="Agreed Amount">
				<el-input v-model.number="paymentData.agreed_amount" />
			</el-form-item>
			<el-form-item label="Paid Amount">
				<el-input v-model.number="paymentData.paid_amount" />
			</el-form-item>
			<el-form-item label="Proof of Payment">
				<el-input v-model="paymentData.proof_of_payment" />
			</el-form-item>
			<el-form-item label="Paid On Time">
				<el-switch v-model="paymentData.paid_ontime" />
			</el-form-item>
			<el-form-item label="Payment Date">
				<el-date-picker v-model="paymentData.payment_date" type="date" placeholder="Select payment date" />
			</el-form-item>
		</el-form>
		<template #footer>
			<div class="dialog-footer">
				<el-button @click="close">Cancel</el-button>
				<el-button type="primary" @click="savePayment">Save</el-button>
			</div>
		</template>
	</el-dialog>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue';
import ApiService from '../api';

interface Payment {
	user_id: string;
	agreed_amount: number;
	paid_amount: number;
	proof_of_payment: string;
	paid_ontime: boolean;
	payment_date?: string;
}

interface User {
	user_id: string;
	firstname: string;
	lastname: string;
}

export default defineComponent({
	props: {
		visible: {
			type: Boolean,
			required: true,
		},
		users: {
			type: Array as () => User[],
			required: true,
		},
	},
	emits: ['close', 'save'],
	setup(props, { emit }) {
		const paymentData = ref<Payment>({
			user_id: '',
			agreed_amount: 0,
			paid_amount: 0,
			proof_of_payment: '',
			paid_ontime: false,
			payment_date: ''
		});

		// Get user ID from token
		let userID = '';
		const token = localStorage.getItem('token');
		if (token) {
			try {
				const payload = JSON.parse(atob(token.split('.')[1]));
				userID = payload.userID;
			} catch (error) {
				console.error('Error decoding token:', error);
			}
		}

		watch(() => props.visible, (newValue) => {
			if (!newValue) {
				paymentData.value = {
					user_id: '',
					agreed_amount: 0,
					paid_amount: 0,
					proof_of_payment: '',
					paid_ontime: false,
					payment_date: '',
				};
			}
		});

		const close = () => {
			emit('close');
		};

		const savePayment = async () => {
			if (!userID) {
				console.error('User ID is required to add a payment.');
				return;
			}

			const paymentDate = paymentData.value.payment_date || new Date().toISOString().split('T')[0];

			const remainingAmount = paymentData.value.agreed_amount - paymentData.value.paid_amount;

			try {
				await ApiService.createPayment(userID, {
					...paymentData.value,
					remaining_amount: remainingAmount,
					payment_date: paymentDate,
				});
				emit('save');
				close();
			} catch (error) {
				console.error('Error saving payment:', error);
			}
		};

		return {
			paymentData,
			close,
			savePayment,
		};
	},
});
</script>

<style scoped>
.dialog-footer {
	text-align: right;
}
</style>
