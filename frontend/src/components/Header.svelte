<script>
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  function handleLogout() {
    dispatch("logout");
  }
</script>

<svelte:head>
  <link
    href="https://fonts.googleapis.com/icon?family=Material+Icons"
    rel="stylesheet"
  />
</svelte:head>

<header>
  <div class="logo">
    <span class="material-icons">check_circle</span>
    <h1>Todo App</h1>
  </div>
  <button on:click={handleLogout} class="logout-button">
    <span class="material-icons">logout</span>
    <span class="button-text">Logout</span>
  </button>
</header>

<style>
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5em 1em;
    background-color: #ffffff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    height: 60px;
  }

  .logo {
    display: flex;
    align-items: center;
  }

  .logo .material-icons {
    font-size: 24px;
    color: #7b68ee;
    margin-right: 8px;
  }

  h1 {
    font-size: 20px;
    font-weight: 600;
    color: #2c3e50;
    margin: 0;
  }

  .logout-button {
    display: flex;
    align-items: center;
    background-color: transparent;
    border: none;
    cursor: pointer;
    color: #6b6b6b;
    font-size: 14px;
    padding: 8px 12px;
    border-radius: 4px;
    transition: all 0.3s ease;
  }

  .logout-button:hover {
    background-color: #f0f0f0;
    color: #7b68ee;
  }

  .logout-button .material-icons {
    font-size: 18px;
    margin-right: 4px;
  }

  .material-icons {
    font-family: "Material Icons";
    font-weight: normal;
    font-style: normal;
    font-size: 24px;
    line-height: 1;
    letter-spacing: normal;
    text-transform: none;
    display: inline-block;
    white-space: nowrap;
    word-wrap: normal;
    direction: ltr;
    -webkit-font-feature-settings: "liga";
    -webkit-font-smoothing: antialiased;
  }

  @media (max-width: 480px) {
    .button-text {
      display: none;
    }

    .logout-button {
      padding: 8px;
    }

    .logout-button .material-icons {
      margin-right: 0;
    }
  }
</style>
