<template>
	<el-dialog v-model="visible" title="Add Payment" width="50%" @close="close">
		<el-form :model="newPayment" label-position="top">
			<el-form-item label="Agreed Amount">
				<el-input v-model.number="newPayment.agreed_amount" />
			</el-form-item>
			<el-form-item label="Paid Amount">
				<el-input v-model.number="newPayment.paid_amount" />
			</el-form-item>
			<el-form-item label="Proof of Payment">
				<el-input v-model="newPayment.proof_of_payment" />
			</el-form-item>
			<el-form-item label="Paid On Time">
				<el-switch v-model="newPayment.paid_ontime" />
			</el-form-item>
		</el-form>
		<template #footer>
			<el-button @click="close">Cancel</el-button>
			<el-button type="primary" @click="savePayment">Save</el-button>
		</template>
	</el-dialog>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import ApiService from '../api';

export default defineComponent({
	props: {
		visible: {
			type: Boolean,
			required: true,
		},
	},
	emits: ['close', 'save'],
	setup(_, { emit }) {
		const newPayment = ref({
			agreed_amount: 0,
			paid_amount: 0,
			proof_of_payment: '',
			paid_ontime: false,
		});

		const token = localStorage.getItem('token');
		let userID = '';

		if (token) {
			try {
				const payload = JSON.parse(atob(token.split('.')[1])); // Decode JWT token
				userID = payload.userID;
			} catch (error) {
				console.error('Error decoding token:', error);
			}
		}

		const close = () => {
			emit('close');
		};

		const savePayment = async () => {
			if (!userID) {
				console.error('User ID is required to add a payment.');
				return;
			}

			const remainingAmount = newPayment.value.agreed_amount - newPayment.value.paid_amount;

			try {
				await ApiService.createPayment(userID, {
					...newPayment.value,
					remaining_amount: remainingAmount,
				});
				emit('save');
				close();
			} catch (error) {
				console.error('Error saving payment:', error);
			}
		};

		return {
			newPayment,
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
