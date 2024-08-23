<script>
  import { navigate } from "svelte-routing";
  import { API_BASE_URL } from "../config";

  let email = "";
  let password = "";
  let error = "";

  async function handleLogin() {
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (response.ok) {
      localStorage.setItem("token", data.token);
      navigate("/dashboard");
    } else {
      error = data.error;
    }
  }
</script>

<div class="auth-container">
  <div class="auth-box">
    <h2>Welcome back!</h2>
    <p class="subtitle">Log in to continue</p>
    {#if error}
      <p class="error">{error}</p>
    {/if}
    <form on:submit|preventDefault={handleLogin}>
      <div class="input-group">
        <label for="email">Email</label>
        <input type="email" id="email" bind:value={email} required />
      </div>
      <div class="input-group">
        <label for="password">Password</label>
        <input type="password" id="password" bind:value={password} required />
      </div>
      <button type="submit" class="submit-btn">Log In</button>
    </form>
    <p class="toggle-text">
      Don't have an account?
      <a href="/register" class="toggle-link">Sign Up</a>
    </p>
  </div>
</div>
