<script>
  export let setActiveView;
</script>

<aside>
  <nav>
    <ul>
      <li><button on:click={() => setActiveView("tasks")}>Tasks</button></li>
      <li>
        <button on:click={() => setActiveView("calendar")}>Calendar</button>
      </li>
    </ul>
  </nav>
</aside>

<style>
  aside {
    background-color: #f0f0f0;
    padding: 1rem;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    margin-bottom: 1rem;
  }

  button {
    width: 100%;
    padding: 0.5rem;
    background-color: transparent;
    border: none;
    text-align: left;
    cursor: pointer;
  }

  button:hover {
    background-color: #e0e0e0;
  }
</style>
