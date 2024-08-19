<template>
	<div id="app" class="app-container">
		<nav-menu class="nav-menu" v-if="!$route.meta.hideNav" />
		<div class="content">
			<router-view />
		</div>
	</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { useRoute } from 'vue-router';
import './assets/style.css';
import NavMenu from './components/NavMenu.vue';

export default defineComponent({
	name: 'App',
	components: {
		NavMenu,
	},
	setup() {
		const route = useRoute();

		return {
			route,
		};
	},
});
</script>

<style>
#app {
	width: 100%;
	height: 100%;
}

.app-container {
	display: flex;
	min-height: 100vh;
	width: 100%;
}

.nav-menu {
	flex: 0 0 12%;
	width: 100%;
	background-color: #f0f2f5;
}

.content {
	width: 100%;
	overflow: hidden;
}
</style>
