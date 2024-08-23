<script>
  import { createEventDispatcher } from "svelte";

  export let setActiveView;
  export let activeView;
  export let isAdmin;

  let isCollapsed = false;
  const dispatch = createEventDispatcher();

  function toggleCollapse() {
    isCollapsed = !isCollapsed;
    dispatch("toggle", { collapsed: isCollapsed });
  }

  $: navItems = [
    { id: "tasks", label: "Tasks", icon: "task_alt" },
    { id: "calendar", label: "Calendar", icon: "calendar_today" },
    ...(isAdmin
      ? [{ id: "users", label: "User Management", icon: "people" }]
      : []),
  ];
</script>

<aside class:collapsed={isCollapsed}>
  <button class="collapse-btn" on:click={toggleCollapse}>
    <span class="material-icons">
      {isCollapsed ? "chevron_right" : "chevron_left"}
    </span>
  </button>
  <nav>
    <ul>
      {#each navItems as item}
        <li>
          <button
            on:click={() => setActiveView(item.id)}
            class:active={activeView === item.id}
          >
            <span class="material-icons">{item.icon}</span>
            {#if !isCollapsed}
              <span class="button-text">{item.label}</span>
            {/if}
          </button>
        </li>
      {/each}
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
