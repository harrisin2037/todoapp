<script>
  export let setActiveView;
</script>

<nav class="mobile-nav">
  <button on:click={() => setActiveView("tasks")}>Tasks</button>
  <button on:click={() => setActiveView("calendar")}>Calendar</button>
</nav>

<style>
  .mobile-nav {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: #f0f0f0;
    display: flex;
    justify-content: space-around;
    padding: 0.5rem;
  }

  button {
    padding: 0.5rem 1rem;
    background-color: transparent;
    border: none;
    cursor: pointer;
  }

  button:hover {
    background-color: #e0e0e0;
  }
</style>
