import { ElMessage } from 'element-plus';

class ApiService {
	private baseURL: string;

	constructor(baseURL: string) {
		this.baseURL = baseURL;
	}

	private async apiRequest(endpoint: string, options: RequestInit) {
		const response = await fetch(`${this.baseURL}${endpoint}`, options);

		if (!response.ok) {
			let errorData;
			try {
				errorData = await response.json();
			} catch (err) {
				errorData = { message: 'Failed to parse error response' };
			}
			const errorMessage = errorData.message || 'Unknown error occurred';
			throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorMessage}`);
		}

		return response.json();
	}

	async getUsers() {
		return this.apiRequestWithMessage('/users', { method: 'GET' }, 'Users retrieved successfully');
	}

	async getUser(userID: string) {
		this.ensureIdExists(userID, 'User ID is required');
		return this.apiRequestWithMessage(`/user/${userID}`, { method: 'GET' }, 'User retrieved successfully');
	}

	async register(user: {
		firstname: string;
		lastname: string;
		email: string;
		password: string;
		user_role: string;
		joined_date: string;
	}) {
		this.ensureValidUser(user);
		await this.apiRequest('/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(user),
		});
		ElMessage.success('User registered successfully');
		return this.login(user.email, user.password);
	}

	async login(email: string, password: string) {
		this.ensureValidLogin(email, password);
		try {
			const data = await this.apiRequest('/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password }),
			});
			ElMessage.success(data.message);
			localStorage.setItem('token', data.token);
			return data;
		} catch (error: any) {
			// TODO: Handle specific error messages
			ElMessage.error('Login failed: Incorrect email or password');
			console.error('Login error:', error);

			throw error;
		}
	}

	logout() {
		try {
			localStorage.removeItem('token');
			ElMessage.success('Logout successful');
		} catch (error) {
			ElMessage.error('Error during logout');
			throw error;
		}
	}

	async updateUser(userID: string, updatedData: {
		firstname: string;
		lastname: string;
		email: string;
		userRole: string;
	}) {
		this.ensureIdExists(userID, 'User ID is required');
		return this.apiRequestWithMessage(`/user/${userID}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(updatedData),
		}, 'User updated successfully');
	}

	async deleteUser(userID: string) {
		this.ensureIdExists(userID, 'User ID is required');
		return this.apiRequestWithMessage(`/user/${userID}`, { method: 'DELETE' }, 'User deleted successfully');
	}

	async getAllPayments() {
		return this.apiRequestWithMessage('/payments', { method: 'GET' }, 'Payments retrieved successfully');
	}

	async getPaymentsByUser(userID: string) {
		this.ensureIdExists(userID, 'User ID is required');
		return this.apiRequest(`/user/${userID}/payments`, { method: 'GET' });
	}

	async getPayment(userID: string, paymentID: string) {
		this.ensureIdExists(userID, 'User ID is required');
		this.ensureIdExists(paymentID, 'Payment ID is required');
		return this.apiRequestWithMessage(`/user/${userID}/payment/${paymentID}`, { method: 'GET' }, 'Payment retrieved successfully');
	}

	async createPayment(userID: string, paymentData: object) {
		this.ensureIdExists(userID, 'User ID is required');
		return this.apiRequestWithMessage(`/user/${userID}/payments`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(paymentData),
		}, 'Payment created successfully');
	}

	async updatePayment(paymentID: string, updatedData: object) {
		this.ensureIdExists(paymentID, 'Payment ID is required');
		return this.apiRequestWithMessage(`/payment/${paymentID}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(updatedData),
		}, 'Payment updated successfully');
	}


	async deletePayment(paymentID: string) {
		this.ensureIdExists(paymentID, 'Payment ID is required');
		return this.apiRequestWithMessage(`/payment/${paymentID}`, { method: 'DELETE' }, 'Payment deleted successfully');
	}

	// Helper Methods
	private async apiRequestWithMessage(endpoint: string, options: RequestInit, successMessage: string) {
		try {
			const response = await this.apiRequest(endpoint, options);
			ElMessage.success(successMessage);
			return response;
		} catch (error) {
			ElMessage.error(error || 'An error occurred');
			throw error;
		}
	}

	private ensureIdExists(id: string, errorMessage: string) {
		if (!id) throw new Error(errorMessage);
	}

	private ensureValidUser(user: any) {
		if (!user.email || !user.password) {
			throw new Error('Email and password are required');
		}
	}

	private ensureValidLogin(email: string, password: string) {
		if (!email || !password) {
			throw new Error('Email and password are required');
		}
	}
}

export default new ApiService('http://localhost:5000');
