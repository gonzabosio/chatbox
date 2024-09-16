<script setup>
import { RouterLink } from 'vue-router';
import router from '../router';
import axios from 'axios';
import { ref } from 'vue';

let username = ''
let password = ''
const formError = ref('')
const baseUrl = import.meta.env.VITE_BACK_BASE_URL
async function userLogin() {
	if (username.length > 12) {
		formError.value = 'Username must contain 12 characters or fewer'
	} else if (password.length > 20) {
		formError.value = 'Password must contain 20 characters or fewer'
	} else {
		await axios({
			method: 'post',
			url: `${baseUrl}/signin`,
			data: {
				name: username,
				password: password
			}
		}).then(res => {
			console.log('Message: ' + res.data.message)
			console.log('Status: ' + res.status)
			console.log('StatusText: ' + res.statusText)
			localStorage.setItem('user-id', res.data.user.id)
			localStorage.setItem('access-token', res.data.access_token)
			localStorage.setItem('session-id', res.data.session_id)
			router.push('/home')
		}).catch(err => {
			username = ''
			password = ''
			console.log('Login failed: ' + err.response.data.message)
			console.log('Error details: ' + err.response.data.error)
			formError.value = err.response.data.message
		})
	}
}
</script>

<template>
	<div id="container">
		<div class="card">
			<h2>Login</h2>
			<form @submit.prevent="userLogin" class="card-form">
				<div class="input">
					<input type="text" v-model="username" class="input-field" required />
					<label class="input-lbl">Username</label>
				</div>
				<div class="input">
					<input type="password" v-model="password" class="input-field" required />
					<label class="input-lbl">Password</label>
				</div>
				<button type="submit" class="submit-btn">Sign in</button>
				<p id="form-error">{{ formError }}</p>
				<div class="below-form">
					<p class="question">Don't have an account?</p>
					<RouterLink to='/register' class="to-register-link">Create an account</RouterLink>
				</div>
			</form>
		</div>
	</div>
</template>

<style scoped>
#container {
	display: flex;
	align-items: center;
	justify-content: center;
}

input {
	appearance: none;
	border-radius: 0;
}

.card {
	margin: 2rem auto;
	display: flex;
	flex-direction: column;
	width: 100%;
	max-width: 425px;
	background-color: #2e2e2e;
	border-radius: 10px;
	box-shadow: 0 10px 20px 0 rgba(#999, .25);
	padding: .75rem;

	h2 {
		display: flex;
		letter-spacing: 4px;
		align-items: center;
		justify-content: center;
	}
}

.card-form {
	padding: 1rem 1rem;
}

.input {
	display: flex;
	flex-direction: column-reverse;
	position: relative;
	padding-top: 1.5rem;

	&+.input {
		margin-top: 1.5rem;
	}
}

.input-lbl {
	color: rgb(120, 141, 155);
	position: absolute;
	top: 1.5rem;
	transition: .25s ease;
}

.input-field {
	border: 0;
	z-index: 1;
	background-color: transparent;
	border-bottom: 2px solid darkgray;
	font: inherit;
	font-size: 1.125rem;
	padding: .25rem 0;

	&:focus,
	&:valid {
		outline: 0;
		border-bottom-color: #84bac5;

		&+.input-lbl {
			color: #84bac5;
			transform: translateY(-1.5rem);
		}
	}
}

.submit-btn {
	font: inherit;
	font-size: 1.25rem;
	padding: 0.5em;
	margin-top: 40px;
	width: 100%;
	font-weight: 600;
	background-color: #6dadb9;
	border-radius: 6px;
	border: 0;
	cursor: pointer;
	transition: 0.3s;

	&:hover {
		background-color: #84bac5;
	}

	&:active {
		transform: scale(.96)
	}
}

#form-error {
	color: rgb(195, 70, 70);
}

.below-form {
	display: flex;
	align-items: center;

	.question {
		margin-right: 6px;
	}

	.to-register-link {
		color: #85c3cf;
		text-decoration: none;

		&:hover {
			text-decoration: underline;
		}
	}
}
</style>
