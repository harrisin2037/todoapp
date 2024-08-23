<script>
  import { onMount, onDestroy } from "svelte";
  import { fade } from "svelte/transition";
  import { API_BASE_URL } from "../config";

  let notifications = [];
  let socket;

  onMount(() => {
    connectWebSocket();
  });

  onDestroy(() => {
    if (socket) {
      socket.close();
    }
  });

  function connectWebSocket() {
    socket = new WebSocket(
      `${API_BASE_URL.replace("http", "ws")}/ws`
    );

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      addNotification(data.message);
    };

    socket.onclose = (event) => {
      console.log("WebSocket connection closed:", event);
      setTimeout(connectWebSocket, 5000);
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  function addNotification(message) {
    const id = Date.now();
    notifications = [...notifications, { id, message }];
    setTimeout(() => {
      removeNotification(id);
    }, 10000);
  }

  function removeNotification(id) {
    notifications = notifications.filter((n) => n.id !== id);
  }
</script>

<div class="notifications">
  {#each notifications as notification (notification.id)}
    <div class="notification" transition:fade={{ duration: 300 }}>
      {notification.message}
    </div>
  {/each}
</div>

<style>
  .notifications {
    position: fixed;
    top: 1rem;
    right: 1rem;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    z-index: 1000;
  }

  .notification {
    background-color: #17a2b8;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 5px;
    margin-bottom: 0.5rem;
    max-width: 300px;
    word-wrap: break-word;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
</style>
