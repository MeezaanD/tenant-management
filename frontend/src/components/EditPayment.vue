<template>
	<el-dialog v-model="localVisible" title="Edit Payment" width="50%" @close="handleClose">
		<el-form :model="localPayment" label-position="top">
			<el-form-item label="Agreed Amount">
				<el-input v-model.number="localPayment.agreed_amount" disabled />
			</el-form-item>
			<el-form-item label="Paid Amount">
				<el-input v-model.number="localPayment.paid_amount" />
			</el-form-item>
			<el-form-item label="Remaining Amount">
				<el-input v-model.number="localPayment.remaining_amount" disabled />
			</el-form-item>
			<el-form-item label="Proof of Payment">
				<el-input v-model="localPayment.proof_of_payment" />
			</el-form-item>
			<el-form-item label="Paid On Time">
				<el-switch v-model="localPayment.paid_ontime" />
			</el-form-item>
		</el-form>
		<template #footer>
			<el-button @click="handleClose">Cancel</el-button>
			<el-button type="primary" @click="savePayment">Save</el-button>
		</template>
	</el-dialog>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue';
import ApiService from '../api';

interface Payment {
	payment_id: number;
	agreed_amount: number;
	paid_amount: number;
	remaining_amount: number;
	proof_of_payment: string;
	paid_ontime: boolean;
}

export default defineComponent({
	props: {
		payment: {
			type: Object as () => Payment,
			required: true,
		},
		visible: {
			type: Boolean,
			required: true,
		},
	},
	emits: ['close', 'save'],
	setup(props, { emit }) {
		const localPayment = ref<Payment>({ ...props.payment });
		const localVisible = ref(props.visible);

		// Watchers to ensure remaining_amount is always calculated correctly
		watch(
			() => localPayment.value.paid_amount,
			(newPaidAmount) => {
				localPayment.value.remaining_amount = localPayment.value.agreed_amount - newPaidAmount;
			}
		);

		watch(
			() => props.visible,
			(newValue) => {
				localVisible.value = newValue;
			}
		);

		watch(
			() => props.payment,
			(newPayment) => {
				localPayment.value = { ...newPayment };
			}
		);

		const handleClose = () => {
			emit('close');
		};

		const savePayment = async () => {
			try {
				await ApiService.updatePayment(
					localPayment.value.payment_id.toString(),
					localPayment.value
				);
				emit('save');
			} catch (error) {
				console.error('Error updating payment:', error);
			}
		};

		return {
			localPayment,
			localVisible,
			handleClose,
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
