<script>
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  function handleLogout() {
    dispatch("logout");
  }
</script>

<header>
  <h1>Todo App</h1>
  <button on:click={handleLogout}>Logout</button>
</header>

<style>
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1em;
    background-color: #f0f0f0;
  }
</style>
