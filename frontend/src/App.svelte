<script>
  import { onMount, onDestroy } from "svelte";
  import Sidebar from "./components/Sidebar.svelte";
  import Header from "./components/Header.svelte";
  import TaskList from "./components/TaskList.svelte";
  import Calendar from "./components/Calendar.svelte";
  import Notification from "./components/Notification.svelte";
  import MobileNav from "./components/MobileNav.svelte";
  import Auth from "./components/Auth.svelte";
  import { API_BASE_URL } from "./config";

  let isMobile;
  let activeView = "tasks";
  let todos = [];
  let pollingInterval;
  let isAuthenticated = false;
  let isLoading = true;

  onMount(async () => {
    checkMobile();
    await checkAuth();
    window.addEventListener("resize", checkMobile);
    isLoading = false;
  });

  onDestroy(() => {
    clearInterval(pollingInterval);
    window.removeEventListener("resize", checkMobile);
  });

  function checkMobile() {
    isMobile = window.innerWidth < 768;
  }

  async function checkAuth() {
    const token = localStorage.getItem("token");
    isAuthenticated = !!token;
    if (isAuthenticated) {
      await fetchTodos();
      startPolling();
    }
  }

  function startPolling() {
    pollingInterval = setInterval(async () => {
      await fetchTodos();
    }, 15000);
  }

  async function fetchTodos() {
    const token = localStorage.getItem("token");
    const response = await fetch(`${API_BASE_URL}/todos`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    if (response.ok) {
      todos = await response.json();
    } else {
      isAuthenticated = false;
      localStorage.removeItem("token");
    }
  }

  function setActiveView(view) {
    activeView = view;
  }

  async function handleLogin() {
    isAuthenticated = true;
    await fetchTodos();
    startPolling();
  }

  function handleLogout() {
    isAuthenticated = false;
    localStorage.removeItem("token");
    clearInterval(pollingInterval);
    todos = [];
  }
</script>

<main
  class={`${isMobile ? "mobile" : "desktop"} ${!isAuthenticated ? "auth" : ""}`}
>
  {#if isLoading}
    <div class="loading">Loading...</div>
  {:else if isAuthenticated}
    {#if !isMobile}
      <Sidebar {setActiveView} {activeView} />
    {/if}

    <div class="content">
      <Header on:logout={handleLogout} />

      {#if activeView === "tasks"}
        <TaskList {todos} {fetchTodos} />
      {:else if activeView === "calendar"}
        <Calendar {todos} />
      {/if}

      <Notification />
    </div>

    {#if isMobile}
      <MobileNav {setActiveView} {activeView} />
    {/if}
  {:else}
    <div class="auth-wrapper">
      <Auth on:login={handleLogin} />
    </div>
  {/if}
</main>

<style>
  :global(body, html) {
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
  }

  main {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    width: 100%;
  }

  .content {
    flex-grow: 1;
    overflow-y: auto;
  }

  .mobile .content {
    padding-bottom: 60px;
  }

  .desktop {
    display: grid;
    grid-template-columns: auto 1fr;
  }

  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    font-size: 1.2em;
    color: #333;
  }

  .auth-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    width: 100%;
    background-color: #fafbfc;
  }

  .desktop.auth {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .desktop .auth-wrapper {
    grid-column: 1 / -1;
  }
</style>
