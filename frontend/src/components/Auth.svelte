<script>
  import { createEventDispatcher } from "svelte";
  import { API_BASE_URL } from "../config.js";

  const dispatch = createEventDispatcher();

  let isLogin = true;
  let username = "";
  let email = "";
  let password = "";
  let error = "";

  function toggleMode() {
    isLogin = !isLogin;
    error = "";
    username = "";
    email = "";
    password = "";
  }

  async function handleSubmit() {
    const endpoint = isLogin ? "login" : "register";
    const body = isLogin
      ? { username, password }
      : { username, email, password };

    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });

      const data = await response.json();

      if (response.ok) {
        if (isLogin) {
          localStorage.setItem("token", data.token);
          dispatch("login");
        } else {
          toggleMode();
        }
      } else {
        error = data.error || "An error occurred";
      }
    } catch (err) {
      error = "Network error. Please try again.";
    }
  }
</script>

<div class="page-container">
  <div class="auth-container">
    <div class="auth-box">
      <h2>{isLogin ? "Welcome back!" : "Create your account"}</h2>
      <p class="subtitle">
        {isLogin ? "Log in to continue" : "Sign up to get started"}
      </p>
      {#if error}
        <p class="error">{error}</p>
      {/if}
      <form on:submit|preventDefault={handleSubmit}>
        <div class="input-group">
          <label for="username">Username</label>
          <input type="text" id="username" bind:value={username} required />
        </div>
        {#if !isLogin}
          <div class="input-group">
            <label for="email">Email</label>
            <input type="email" id="email" bind:value={email} required />
          </div>
        {/if}
        <div class="input-group">
          <label for="password">Password</label>
          <input type="password" id="password" bind:value={password} required />
        </div>
        <button type="submit" class="submit-btn">
          {isLogin ? "Log In" : "Sign Up"}
        </button>
      </form>
      <p class="toggle-text">
        {isLogin ? "Don't have an account?" : "Already have an account?"}
        <button class="toggle-btn" on:click={toggleMode}>
          {isLogin ? "Sign Up" : "Log In"}
        </button>
      </p>
    </div>
  </div>
</div>

<style>
  .page-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #fafbfc;
  }

  .auth-container {
    width: 100%;
    max-width: 400px;
    padding: 20px;
  }

  .auth-box {
    background-color: #ffffff;
    border-radius: 8px;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
    padding: 40px;
  }

  h2 {
    color: #2c2c2c;
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 8px;
    text-align: center;
  }

  .subtitle {
    color: #777;
    font-size: 16px;
    margin-bottom: 24px;
    text-align: center;
  }

  .input-group {
    margin-bottom: 20px;
  }

  label {
    display: block;
    margin-bottom: 6px;
    color: #2c2c2c;
    font-size: 14px;
    font-weight: 500;
  }

  input {
    width: 100%;
    padding: 10px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    font-size: 16px;
    transition: border-color 0.3s ease;
  }

  input:focus {
    outline: none;
    border-color: #7b68ee;
  }

  .submit-btn {
    width: 100%;
    padding: 12px;
    background-color: #7b68ee;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }

  .submit-btn:hover {
    background-color: #6c5ce7;
  }

  .error {
    color: #e74c3c;
    margin-bottom: 16px;
    text-align: center;
    font-size: 14px;
  }

  .toggle-text {
    margin-top: 20px;
    text-align: center;
    font-size: 14px;
    color: #777;
  }

  .toggle-btn {
    background: none;
    border: none;
    color: #7b68ee;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
  }

  .toggle-btn:hover {
    text-decoration: underline;
  }
</style>
