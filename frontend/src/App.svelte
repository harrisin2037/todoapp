<script>
  import { onMount, onDestroy } from "svelte";
  import { isLoading } from "./loadingStore";
  import Sidebar from "./components/Sidebar.svelte";
  import Header from "./components/Header.svelte";
  import TaskList from "./components/TaskList.svelte";
  import Calendar from "./components/Calendar.svelte";
  import Notification from "./components/Notification.svelte";
  import MobileNav from "./components/MobileNav.svelte";
  import Auth from "./components/Auth.svelte";
  import UserManagement from "./components/UserManagement.svelte";
  import TaskTemplateManager from "./components/TaskTemplateManager.svelte";
  import { API_BASE_URL } from "./config";
  import { userStore } from "./userStore";

  let isMobile;
  let activeView = "tasks";
  let todos = [];
  let pollingInterval;
  let isAuthenticated = false;
  let isAdmin = false;
  let socket;
  let onlineUsers = [];

  onMount(async () => {
    checkMobile();
    await checkAuth();
    window.addEventListener("resize", checkMobile);

    socket = new WebSocket(`${API_BASE_URL.replace("http", "ws")}/ws`);

    socket.addEventListener("open", () => {
      console.log("WebSocket connection opened");
      updateMyOnlineStatus();
      whoAreOnline();
      setInterval(updateMyOnlineStatus, 30000);
      isLoading.set(false);
    });

    socket.addEventListener("message", (event) => {
      const data = JSON.parse(event.data);
      if (data.action === "onlineStatus") {
        onlineUsers = data.onlineUsers;
      }
    });

    socket.addEventListener("close", () => {
      console.log("WebSocket connection closed");
    });

    socket.addEventListener("error", (error) => {
      console.error("WebSocket error:", error);
    });

    await userStore.loadUsers();
  });

  onDestroy(() => {
    clearInterval(pollingInterval);
    window.removeEventListener("resize", checkMobile);
    if (socket) {
      socket.close();
    }
  });

  function updateMyOnlineStatus() {
    socket.send(
      JSON.stringify({
        token: localStorage.getItem("token"),
      })
    );
  }

  function whoAreOnline() {
    socket.send(
      JSON.stringify({
        token: localStorage.getItem("token"),
        action: "who",
      })
    );
  }

  function checkMobile() {
    isMobile = window.innerWidth < 768;
  }

  async function checkAuth() {
    const token = localStorage.getItem("token");
    isAuthenticated = !!token;
    if (isAuthenticated) {
      await fetchTodos();
      await checkAdminStatus();
      startPolling();
    }
  }

  async function checkAdminStatus() {
    const token = localStorage.getItem("token");
    if (!token) {
      isAdmin = false;
      return;
    }
    try {
      const response = await fetch(`${API_BASE_URL}/roles`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        isAdmin = false;
        isLoading.set(true);
        return;
      }

      const data = await response.json();
      data.role === "admin" ? (isAdmin = true) : (isAdmin = false);
    } catch (error) {
      console.error("Error checking admin status:", error);
      isAdmin = false;
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
    await checkAdminStatus();
    startPolling();
  }

  function handleLogout() {
    isAuthenticated = false;
    isAdmin = false;
    localStorage.removeItem("token");
    clearInterval(pollingInterval);
    todos = [];
  }

  function handleTodoUpdate() {
    fetchTodos();
  }
</script>

<main
  class={`${isMobile ? "mobile" : "desktop"} ${!isAuthenticated ? "auth" : ""}`}
>
  {#if $isLoading}
    <div class="loading">Playbook</div>
  {:else if isAuthenticated}
    {#if !isMobile}
      <Sidebar {setActiveView} {activeView} {isAdmin} />
    {/if}

    <div class="content">
      <Header on:logout={handleLogout} {onlineUsers} />

      {#if activeView === "tasks"}
        <TaskList {todos} {fetchTodos} allUsers={$userStore} />
      {:else if activeView === "calendar"}
        <Calendar {todos} />
      {:else if activeView === "users" && isAdmin}
        <UserManagement />
      {:else if activeView === "templates"}
        <TaskTemplateManager />
      {:else}
        <h1>404 - Not Found</h1>
        <!-- need a component for the 404/400/500 -->
      {/if}

      <Notification />
    </div>

    {#if isMobile}
      <MobileNav {setActiveView} {activeView} {isAdmin} />
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

  .loading {
    font-size: 38px;
    font-weight: bold;
    color: #b288c1;
    animation: loadingAnimation 1.5s infinite;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    text-align: center;
  }

  @keyframes loadingAnimation {
    0% {
      opacity: 1;
      color: #a59eab;
    }
    50% {
      opacity: 0.5;
      color: #b2aee6;
    }
    100% {
      opacity: 1;
      color: #ccb8f5;
    }
  }
</style>
