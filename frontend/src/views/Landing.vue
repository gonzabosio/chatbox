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
		<div id="left-content">
			<h1>Chat<span id="box-word">Box</span></h1>
			<h3 style="font-style: italic;">Your Secure and Simple Way to Stay Connected</h3>
			<p>Stay in touch with your friends, family, and colleagues with a fast,
				reliable, and easy-to-use messaging platform designed to keep you connected at all times.</p>
		</div>
		<img src="../assets/chat_landing_draw.svg" alt="chat-draw" width="400">
	</div>
</template>

<style scoped>
#container {
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 60px auto;
	max-width: 1100px;
	padding: 18px;
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
	margin-top: 40px;
}

#container p {
	margin-right: 10px;
	margin-bottom: 30px;
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