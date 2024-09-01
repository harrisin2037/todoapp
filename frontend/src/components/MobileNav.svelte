<script>
  import { onMount } from "svelte";

  export let setActiveView;
  export let activeView;
  export let isAdmin;

  let isNarrow = false;
  let isMobile = false;
  let windowWidth;

  const views = [
    { name: "Tasks", icon: "task_alt" },
    { name: "Calendar", icon: "calendar_today" },
    { name: "Templates", icon: "copy_all" },
    ...(isAdmin ? [{ name: "Users", icon: "people" }] : []),
  ];

  function toggleNavbar() {
    if (!isMobile) {
      isNarrow = !isNarrow;
    }
  }

  $: isMobile = windowWidth <= 768;
  $: if (isMobile) isNarrow = false;

  onMount(() => {
    isNarrow = window.innerWidth < 1024 && window.innerWidth > 768;
  });
</script>

<svelte:window bind:innerWidth={windowWidth} />

<nav class="nav-bar" class:narrow={isNarrow} class:mobile={isMobile}>
  {#if !isMobile}
    <div class="logo">
      <span class="material-icons">menu</span>
      {#if !isNarrow}
        <span class="app-name">Todo App</span>
      {/if}
    </div>
  {/if}
  <div class="nav-items">
    {#each views as view}
      <button
        class:active={activeView === view.name.toLowerCase()}
        on:click={() => setActiveView(view.name.toLowerCase())}
        title={view.name}
      >
        <span class="material-icons">{view.icon}</span>
        {#if !isNarrow || isMobile}
          <span class="button-text">{view.name}</span>
        {/if}
      </button>
    {/each}
  </div>
  {#if !isMobile}
    <button class="collapse-btn" on:click={toggleNavbar}>
      <span class="material-icons">
        {isNarrow ? "chevron_right" : "chevron_left"}
      </span>
    </button>
  {/if}
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
    .nav-bar.mobile {
      position: fixed;
      bottom: 0;
      left: 0;
      right: 0;
      width: 100%;
      height: auto;
      flex-direction: row;
      justify-content: space-around;
      padding: 0.5rem;
      z-index: 1000;
      background-color: #f5f5f5;
      border-top: 1px solid #e0e0e0;
    }

    .nav-bar.mobile .logo {
      display: none;
    }

    .nav-bar.mobile .nav-items {
      flex-direction: row;
      padding: 0;
      width: 100%;
    }

    .nav-bar.mobile button {
      flex-direction: column;
      align-items: center;
      padding: 0.5rem;
    }

    .nav-bar.mobile .material-icons {
      margin-right: 0;
      margin-bottom: 4px;
    }

    .nav-bar.mobile .button-text {
      font-size: 12px;
    }

    .nav-bar.mobile .collapse-btn {
      display: none;
    }
  }

  @media (max-width: 360px) {
    .nav-bar.mobile .button-text {
      display: none;
    }
  }
</style>
