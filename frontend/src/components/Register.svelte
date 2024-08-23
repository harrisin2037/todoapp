<script>
  import { navigate } from "svelte-routing";
  import { API_BASE_URL } from "../config";

  let username = "";
  let email = "";
  let password = "";
  let error = "";

  async function handleRegister() {
    const response = await fetch(`${API_BASE_URL}/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, email, password }),
    });

    const data = await response.json();

    if (response.ok) {
      navigate("/login");
    } else {
      error = data.error;
    }
  }
</script>

<div class="auth-container">
  <div class="auth-box">
    <h2>Create your account</h2>
    <p class="subtitle">Sign up to get started</p>
    {#if error}
      <p class="error">{error}</p>
    {/if}
    <form on:submit|preventDefault={handleRegister}>
      <div class="input-group">
        <label for="username">Username</label>
        <input type="text" id="username" bind:value={username} required />
      </div>
      <div class="input-group">
        <label for="email">Email</label>
        <input type="email" id="email" bind:value={email} required />
      </div>
      <div class="input-group">
        <label for="password">Password</label>
        <input type="password" id="password" bind:value={password} required />
      </div>
      <button type="submit" class="submit-btn">Sign Up</button>
    </form>
    <p class="toggle-text">
      Already have an account?
      <a href="/login" class="toggle-link">Log In</a>
    </p>
  </div>
</div>
