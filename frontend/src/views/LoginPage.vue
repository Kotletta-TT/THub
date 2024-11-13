<template>
  <div class="login-page">
    <form @submit.prevent="login">
      <input v-model="username" placeholder="Username" />
      <input v-model="password" type="password" placeholder="Password" />
      <div class="button-group">
        <button type="submit">Login</button>
        <button type="button" @click="showSignUpDialog = true">Sign Up</button>
      </div>
      <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>
    </form>

    <!-- Всплывающее окно для регистрации -->
    <div v-if="showSignUpDialog" class="dialog">
      <h3>Sign Up</h3>
      <input v-model="signUpData.username" placeholder="Username" />
      <input v-model="signUpData.email" placeholder="Email" type="email" />
      <input v-model="signUpData.password" type="password" placeholder="Password" />
      <input v-model="signUpData.repeatPassword" type="password" placeholder="Repeat Password" />
      <p class="error-message" v-if="signUpErrorMessage">{{ signUpErrorMessage }}</p>
      <div class="dialog-buttons">
        <button @click="submitSignUp">Submit</button>
        <button @click="closeSignUpDialog">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      username: '',
      password: '',
      errorMessage: '',
      showSignUpDialog: false,
      signUpErrorMessage: '',
      signUpData: {
        username: '',
        email: '',
        password: '',
        repeatPassword: ''
      }
    };
  },
  created() {
    const token = localStorage.getItem('token');
    if (token) {
      this.$router.push('/main');
    }
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('/api/user/login', {
          username: this.username,
          password: this.password
        });
        localStorage.setItem('token', response.data.token);
        this.$router.push('/main');
      } catch (error) {
        this.errorMessage = 'Login failed';
      }
    },
    closeSignUpDialog() {
      this.showSignUpDialog = false;
      this.signUpErrorMessage = '';
      this.signUpData = { username: '', email: '', password: '', repeatPassword: '' };
    },
    async submitSignUp() {
      // Проверка заполнения полей и совпадения паролей
      const { username, email, password, repeatPassword } = this.signUpData;
      if (!username || !email || !password || !repeatPassword) {
        this.signUpErrorMessage = 'All fields are required.';
        return;
      }
      if (password !== repeatPassword) {
        this.signUpErrorMessage = 'Passwords do not match.';
        return;
      }

      // Отправка запроса на регистрацию
      try {
        await axios.post('/api/user/signup', {
          username,
          email,
          password
        });
        alert('Registration successful! You can now log in.');
        this.closeSignUpDialog();
      } catch (error) {
        this.signUpErrorMessage = 'Registration failed. Please try again.';
      }
    }
  }
};
</script>

<style>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

form {
  background: #333;
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

input {
  width: 100%;
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: none;
}

.button-group {
  display: flex;
  gap: 10px;
  width: 100%;
}

button {
  background-color: #3b82f6;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  flex: 1;
}

.error-message {
  color: #ff4d4d;
  margin-top: 10px;
  text-align: center;
}

.dialog {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #222;
  padding: 20px;
  border-radius: 8px;
  width: 360px; /* Увеличенная ширина окна для соответствия кнопкам */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

.dialog h3 {
  margin-bottom: 10px;
  text-align: center;
  color: #fff;
}

.dialog input {
  width: 100%; /* Растягиваем поле ввода на всю ширину диалогового окна */
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: none;
  box-sizing: border-box;
}

.dialog-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
}

.dialog button {
  padding: 10px 20px;
}
</style>
