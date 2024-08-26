<script>
  import { onMount } from "svelte";
  import { API_BASE_URL } from "../config";

  let users = [];
  let newUser = { username: "", email: "", password: "", role: "user" };
  let editingUser = null;
  let showModal = false;
  let ws;

  $: currentUser = editingUser || newUser;

  onMount(() => {
    fetchUsers();
    // setupWebSocket();
  });

  async function fetchUsers() {
    const response = await fetch(`${API_BASE_URL}/admin/users`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });

    if (response.ok) {
      users = (await response.json()).users;
      console.log("users: ", users);
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  function setupWebSocket() {
    const wsUrl = `${API_BASE_URL.replace(/^http/, "ws")}/ws`;
    console.log("Connecting to WebSocket:", wsUrl);

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log("WebSocket connection opened");
    };

    ws.onmessage = (event) => {
      console.log("WebSocket message received:", event.data);
      try {
        const data = JSON.parse(event.data);
        if (data.type === "online_users") {
          console.log("Updating online users:", data.users);
          updateOnlineStatus(data.users);
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.onclose = (event) => {
      console.log("WebSocket connection closed:", event);
    };
  }

  function updateOnlineStatus(onlineUsers) {
    users = users.map((user) => ({
      ...user,
      online: onlineUsers.includes(user.id.toString()),
    }));
  }

  async function addUser() {
    const response = await fetch(`${API_BASE_URL}/admin/users`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify(newUser),
    });

    if (response.ok) {
      await fetchUsers();
      newUser = { username: "", email: "", password: "", role: "user" };
      showModal = false;
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  async function updateUser() {
    const response = await fetch(
      `${API_BASE_URL}/admin/users/${editingUser.id}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify(editingUser),
      }
    );

    if (response.ok) {
      await fetchUsers();
      editingUser = null;
      showModal = false;
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  async function deleteUser(userId) {
    if (confirm("Are you sure you want to delete this user?")) {
      const response = await fetch(`${API_BASE_URL}/admin/users/${userId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });

      if (response.ok) {
        await fetchUsers();
      } else {
        const errorData = await response.json();
        alert(`Error: ${errorData.error}`);
      }
    }
  }

  function openModal(user = null) {
    if (user) {
      editingUser = { ...user };
    } else {
      editingUser = null;
    }
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    editingUser = null;
  }
</script>

<div class="user-management">
  <h2>User Management</h2>
  <button class="add-btn" on:click={() => openModal()}>Add User</button>

  <table>
    <thead>
      <tr>
        <th>Username</th>
        <th>Email</th>
        <th>Online</th>
        <th>Role</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each users as user (user.id)}
        <tr>
          <td data-label="Username">{user.username}</td>
          <td data-label="Email">{user.email}</td>
          <td data-label="Online">
            <span class="status-indicator {user.online ? 'online' : 'offline'}"
            ></span>
            {user.online ? "Online" : "Offline"}
          </td>
          <td data-label="Role">{user.role}</td>
          <td data-label="Actions">
            <button class="edit-btn" on:click={() => openModal(user)}>
              <span class="material-icons">edit</span>
            </button>
            <button class="delete-btn" on:click={() => deleteUser(user.id)}>
              <span class="material-icons">delete</span>
            </button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>

  {#if showModal}
    <div class="modal">
      <div class="modal-content">
        <h3>{editingUser ? "Edit User" : "Add User"}</h3>
        <form on:submit|preventDefault={editingUser ? updateUser : addUser}>
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            placeholder="Username"
            bind:value={currentUser.username}
            required
          />
          {#if !editingUser}
            <label for="password">Password</label>
            <input
              type="password"
              id="password"
              placeholder="Password"
              bind:value={currentUser.password}
              required
            />
          {/if}
          <label for="email">Email</label>
          <input
            type="email"
            id="email"
            placeholder="Email"
            bind:value={currentUser.email}
            required
          />
          <label for="role">Role</label>
          <select id="role" bind:value={currentUser.role}>
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </select>
          <div class="form-actions">
            <button type="button" class="cancel-btn" on:click={closeModal}
              >Cancel</button
            >
            <button type="submit" class="submit-btn"
              >{editingUser ? "Update" : "Add"}</button
            >
          </div>
        </form>
      </div>
    </div>
  {/if}
</div>

<style>
  .user-management {
    padding: 2rem;
  }

  h2 {
    margin-bottom: 1rem;
  }

  .add-btn {
    background-color: #7b68ee;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    margin-bottom: 1rem;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    text-align: left;
    padding: 0.5rem;
    border-bottom: 1px solid #e8ecee;
  }

  th {
    background-color: #f7f9fb;
    font-weight: 600;
  }

  .edit-btn,
  .delete-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    margin: 0 4px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: #555;
  }

  .edit-btn:hover,
  .delete-btn:hover {
    color: #000;
  }

  .material-icons {
    font-size: 24px;
  }

  .modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .modal-content {
    background-color: white;
    padding: 2rem;
    border-radius: 8px;
    width: 90%;
    max-width: 500px;
  }

  form {
    display: flex;
    flex-direction: column;
  }

  label {
    margin-bottom: 0.5rem;
  }

  input,
  select {
    margin-bottom: 1rem;
    padding: 0.7rem;
    border: 1px solid #e8ecee;
    border-radius: 4px;
    font-size: 1rem;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
  }

  .cancel-btn,
  .submit-btn {
    padding: 0.7rem 1.5rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
  }

  .cancel-btn {
    background-color: #e8ecee;
    color: #2c3e50;
  }

  .submit-btn {
    background-color: #7b68ee;
    color: white;
  }

  .status-indicator {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    margin-right: 5px;
  }

  .status-indicator.online {
    background-color: #4caf50;
  }

  .status-indicator.offline {
    background-color: #f44336;
  }

  @media (max-width: 600px) {
    .user-management {
      padding: 1rem;
    }

    table,
    thead,
    tbody,
    th,
    td,
    tr {
      display: block;
    }

    thead tr {
      position: absolute;
      top: -9999px;
      left: -9999px;
    }

    tr {
      margin-bottom: 1rem;
      border: 1px solid #e8ecee;
    }

    td {
      border: none;
      position: relative;
      padding-left: 50%;
    }

    td:before {
      position: absolute;
      top: 6px;
      left: 6px;
      width: 45%;
      padding-right: 10px;
      white-space: nowrap;
      content: attr(data-label);
      font-weight: bold;
    }
  }
</style>
