<script>
  import { createEventDispatcher } from "svelte";

  export let setActiveView;
  export let activeView;

  let isCollapsed = false;
  const dispatch = createEventDispatcher();

  function toggleCollapse() {
    isCollapsed = !isCollapsed;
    dispatch("toggle", { collapsed: isCollapsed });
  }
</script>

<aside class:collapsed={isCollapsed}>
  <button class="collapse-btn" on:click={toggleCollapse}>
    <span class="material-icons">
      {isCollapsed ? "chevron_right" : "chevron_left"}
    </span>
  </button>
  <nav>
    <ul>
      <li>
        <button
          on:click={() => setActiveView("tasks")}
          class:active={activeView === "tasks"}
        >
          <span class="material-icons">task_alt</span>
          {#if !isCollapsed}
            <span class="button-text">Tasks</span>
          {/if}
        </button>
      </li>
      <li>
        <button
          on:click={() => setActiveView("calendar")}
          class:active={activeView === "calendar"}
        >
          <span class="material-icons">calendar_today</span>
          {#if !isCollapsed}
            <span class="button-text">Calendar</span>
          {/if}
        </button>
      </li>
    </ul>
  </nav>
</aside>

<svelte:head>
  <link
    href="https://fonts.googleapis.com/icon?family=Material+Icons"
    rel="stylesheet"
  />
</svelte:head>

<style>
  .material-icons {
    font-family: "Material Icons";
    font-weight: normal;
    font-style: normal;
    font-size: 24px;
    display: inline-block;
    line-height: 1;
    text-transform: none;
    letter-spacing: normal;
    word-wrap: normal;
    white-space: nowrap;
    direction: ltr;
    -webkit-font-smoothing: antialiased;
    text-rendering: optimizeLegibility;
    -moz-osx-font-smoothing: grayscale;
    font-feature-settings: "liga";
  }
  aside {
    background-color: #f0f0f0;
    padding: 1rem;
    transition: width 0.3s ease;
    width: 200px;
    position: relative;
  }

  aside.collapsed {
    width: 60px;
  }

  .collapse-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    background-color: #e0e0e0;
    border: none;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
  }

  nav {
    margin-top: 40px;
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
    cursor: pointer;
    display: flex;
    align-items: center;
  }

  button .button-text {
    margin-left: 10px;
  }

  button:hover,
  button.active {
    background-color: #e0e0e0;
  }

  .collapsed button {
    justify-content: center;
  }

  .collapsed button .button-text {
    display: none;
  }

  .material-icons {
    font-size: 24px;
  }
</style>
