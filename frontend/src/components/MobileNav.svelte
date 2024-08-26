<script>
  import { onMount } from "svelte";

  export let setActiveView;
  export let activeView;
  export let isAdmin;

  let isNarrow = false;
  let navBar;

  const views = [
    { name: "Tasks", icon: "task_alt" },
    { name: "Calendar", icon: "calendar_today" },
    ...(isAdmin ? [{ name: "Users", icon: "people" }] : []),
  ];

  function toggleNavbar() {
    isNarrow = !isNarrow;
    navBar.style.width = isNarrow ? "60px" : "240px";
  }

  onMount(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      for (let entry of entries) {
        isNarrow = entry.contentRect.width < 200;
      }
    });

    resizeObserver.observe(navBar);

    return () => {
      resizeObserver.disconnect();
    };
  });
</script>

<nav class="nav-bar" class:narrow={isNarrow} bind:this={navBar}>
  <div class="logo">
    <span class="material-icons">menu</span>
    {#if !isNarrow}
      <span class="app-name">Todo App</span>
    {/if}
  </div>
  <div class="nav-items">
    {#each views as view}
      <button
        class:active={activeView === view.name.toLowerCase()}
        on:click={() => setActiveView(view.name.toLowerCase())}
        title={view.name}
      >
        <span class="material-icons">{view.icon}</span>
        {#if !isNarrow}
          <span class="button-text">{view.name}</span>
        {/if}
      </button>
    {/each}
  </div>
  <button class="collapse-btn" on:click={toggleNavbar}>
    <span class="material-icons"
      >{isNarrow ? "chevron_right" : "chevron_left"}</span
    >
  </button>
</nav>

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

  .nav-bar {
    display: flex;
    flex-direction: column;
    width: 240px;
    height: 100vh;
    background-color: #f5f5f5;
    transition: width 0.3s ease;
    position: relative;
    border-right: 1px solid #e0e0e0;
  }

  .nav-bar.narrow {
    width: 60px;
  }

  .logo {
    display: flex;
    align-items: center;
    padding: 1rem;
  }

  .logo .material-icons {
    font-size: 24px;
    color: #7b68ee;
    margin-right: 8px;
  }

  .app-name {
    font-size: 18px;
    font-weight: 600;
    color: #2c3e50;
  }

  .nav-items {
    display: flex;
    flex-direction: column;
    padding: 1rem 0;
  }

  button {
    display: flex;
    align-items: center;
    padding: 0.75rem 1rem;
    background-color: transparent;
    border: none;
    cursor: pointer;
    color: #6b6b6b;
    transition: all 0.3s ease;
    width: 100%;
    text-align: left;
  }

  button:hover {
    background-color: #e0e0e0;
  }

  button.active {
    color: #7b68ee;
    background-color: #e0e0e0;
  }

  .material-icons {
    font-size: 24px;
    margin-right: 16px;
  }

  .button-text {
    font-size: 14px;
  }

  .collapse-btn {
    position: absolute;
    bottom: 1rem;
    left: 0;
    right: 0;
    display: flex;
    justify-content: center;
    padding: 0.5rem;
    background-color: transparent;
    border: none;
    cursor: pointer;
    color: #7b68ee;
  }

  @media (max-width: 768px) {
    .nav-bar {
      position: fixed;
      bottom: 0;
      left: 0;
      right: 0;
      width: 100%;
      height: auto;
      flex-direction: row;
      justify-content: space-around;
      padding: 0.5rem;
    }

    .logo {
      display: none;
    }

    .nav-items {
      flex-direction: row;
      padding: 0;
      width: 100%;
    }

    button {
      flex-direction: column;
      align-items: center;
      padding: 0.5rem;
    }

    .material-icons {
      margin-right: 0;
      margin-bottom: 4px;
    }

    .button-text {
      font-size: 12px;
    }

    .collapse-btn {
      display: none;
    }
  }

  @media (max-width: 360px) {
    .button-text {
      display: none;
    }
  }
</style>
