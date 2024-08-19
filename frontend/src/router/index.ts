import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import Home from '../components/Home.vue';
import AllUsers from '../components/AllUsers.vue';
import AllPayments from '../components/AllPayments.vue';
import Profile from '../components/Profile.vue';
import UserPayments from '../components/UserPayments.vue';

const routes = [
	{
		path: '/', //Login page
		component: Login,
		meta: { hideNav: true }
	},
	{
		path: '/register',
		component: Register,
		meta: { hideNav: true }
	},
	{
		path: '/profile',
		name: 'Profile',
		component: Profile,
	},
	{
		path: '/home',
		component: Home
	},
	{
		path: '/users',
		component: AllUsers
	},
	{
		path: '/payments',
		component: AllPayments
	},
	{
		path: '/user-payments',
		name: 'UserPayments',
		component: UserPayments,
	},
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

export default router;
