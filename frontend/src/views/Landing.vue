<script setup>
import { onBeforeMount } from 'vue';
import router from '../router/index.js'
onBeforeMount(() => {
	if (localStorage.getItem('access-token') !== null) {
		console.log('User is already logged. Redirecting...')
		router.push('/home')
	}
})
</script>

<template>
	<div id="container">
		<div>
			<h1>Chat<span id="box-word">Box</span></h1>
			<h3>Fullstack Chat Application made with Go and Vue</h3>
			<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et
				dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
				ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
				fugiat nulla pariatur</p>
		</div>
		<img src="../assets/chat_landing_draw.svg" alt="chat-draw" width="400">
	</div>
</template>

<style scoped>
#container {
	display: flex;
	align-items: center;
	justify-content: space-around;
	margin: 60px auto;
	max-width: 1100px;
}

#container h1 {
	font-size: 80px;
	font-weight: bold;
	letter-spacing: 12px;
}

#box-word {
	color: #2fa3b3;
}

#container img {
	width: 100%;
	max-width: 400px;
	margin-top: 45px;
}

#container p {
	margin-right: 10px;
	color: darkgray;
	line-height: 24px;
}

@media (max-width: 768px) {
	#container {
		flex-direction: column;
		margin: 20px;
	}

	#container img {
		max-width: 300px;
	}
}
</style>