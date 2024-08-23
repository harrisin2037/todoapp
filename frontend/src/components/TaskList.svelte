<script>
  import TodoItem from "./TodoItem.svelte";
  import { API_BASE_URL } from "../config";

  export let todos;

  let newTodo = { name: "", description: "", due_date: "" };
  let showModal = false;
  let showPending = true;
  let showCompleted = false;
  let searchQuery = "";
  let filterStatus = "all";
  let sortBy = "due_date";
  let sortOrder = "asc";

  $: filteredTodos = todos
    .filter((todo) => {
      const matchesSearch = todo.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesFilter =
        filterStatus === "all" || todo.status === filterStatus;
      return matchesSearch && matchesFilter;
    })
    .sort((a, b) => {
      const factor = sortOrder === "asc" ? 1 : -1;
      if (sortBy === "name") {
        return factor * a.name.localeCompare(b.name);
      } else if (sortBy === "due_date") {
        return factor * (new Date(a.due_date) - new Date(b.due_date));
      }
      return 0;
    });

  $: pendingTodos = filteredTodos.filter((todo) => todo.status !== "completed");
  $: completedTodos = filteredTodos.filter(
    (todo) => todo.status === "completed"
  );

  async function fetchTodos() {
    const response = await fetch(`${API_BASE_URL}/todos`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });

    if (response.ok) {
      todos = await response.json();
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  async function addTodo() {
    const response = await fetch(`${API_BASE_URL}/todos`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify(newTodo),
    });
    if (response.ok) {
      await fetchTodos();
      newTodo = { name: "", description: "", due_date: "" };
      showModal = false;
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  function toggleModal() {
    showModal = !showModal;
  }

  function togglePending() {
    showPending = !showPending;
  }

  function toggleCompleted() {
    showCompleted = !showCompleted;
  }
</script>

<div class="task-list">
  <div class="task-list-header">
    <h2>Tasks</h2>
    <button class="add-btn" on:click={toggleModal}>
      <span class="plus-icon">+</span> Add Task
    </button>
  </div>

  <div class="search-filter-sort">
    <input
      type="text"
      bind:value={searchQuery}
      placeholder="Search tasks..."
      class="search-input"
    />
    <select bind:value={filterStatus} class="filter-select">
      <option value="all">All</option>
      <option value="pending">Pending</option>
      <option value="completed">Completed</option>
    </select>
    <select bind:value={sortBy} class="sort-select">
      <option value="due_date">Due Date</option>
      <option value="name">Name</option>
    </select>
    <button
      class="sort-order-btn"
      on:click={() => (sortOrder = sortOrder === "asc" ? "desc" : "asc")}
    >
      {sortOrder === "asc" ? "▲" : "▼"}
    </button>
  </div>

  <div class="task-list-content">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="tasks-header" on:click={togglePending}>
      <h3>Pending Tasks ({pendingTodos.length})</h3>
      <span class="expand-icon">{showPending ? "▼" : "▶"}</span>
    </div>

    {#if showPending}
      <div class="tasks-list">
        {#each pendingTodos as todo (todo.ID)}
          <TodoItem {todo} on:update={fetchTodos} on:delete={fetchTodos} />
        {/each}
      </div>
    {/if}

    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="tasks-header" on:click={toggleCompleted}>
      <h3>Completed Tasks ({completedTodos.length})</h3>
      <span class="expand-icon">{showCompleted ? "▼" : "▶"}</span>
    </div>

    {#if showCompleted}
      <div class="tasks-list">
        {#each completedTodos as todo (todo.ID)}
          <TodoItem {todo} on:update={fetchTodos} on:delete={fetchTodos} />
        {/each}
      </div>
    {/if}
  </div>

  {#if showModal}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal" on:click={toggleModal}>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="modal-content" on:click|stopPropagation>
        <h3>Add New Task</h3>
        <form on:submit|preventDefault={addTodo}>
          <input
            type="text"
            bind:value={newTodo.name}
            placeholder="Task name"
            required
          />
          <textarea
            bind:value={newTodo.description}
            placeholder="Description"
            rows="3"
            class="limited-size"
          ></textarea>
          <input type="date" bind:value={newTodo.due_date} />
          <div class="form-actions">
            <button type="button" class="cancel-btn" on:click={toggleModal}
              >Cancel</button
            >
            <button type="submit" class="submit-btn">Add Task</button>
          </div>
        </form>
      </div>
    </div>
  {/if}
</div>

<style>
  .task-list {
    background-color: #f7f9fb;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .task-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background-color: #ffffff;
    border-bottom: 1px solid #e8ecee;
  }

  .task-list-header h2 {
    font-size: 1.2rem;
    font-weight: 600;
    color: #2c3e50;
  }

  .add-btn {
    background-color: #7b68ee;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    font-weight: 500;
    transition: background-color 0.3s ease;
  }

  .add-btn:hover {
    background-color: #6c5ce7;
  }

  .plus-icon {
    margin-right: 0.5rem;
    font-size: 1.2rem;
  }

  .search-filter-sort {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    padding: 1rem;
    background-color: #ffffff;
    border-bottom: 1px solid #e8ecee;
  }

  .search-input,
  .filter-select,
  .sort-select {
    flex: 1;
    min-width: 120px;
    padding: 0.5rem;
    border: 1px solid #e8ecee;
    border-radius: 4px;
    font-size: 1rem;
  }

  .sort-order-btn {
    padding: 0.5rem 1rem;
    background-color: #e8ecee;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
  }

  .task-list-content {
    flex-grow: 1;
    overflow-y: auto;
    padding: 1rem;
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

  .modal-content h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    color: #2c3e50;
  }

  form {
    display: flex;
    flex-direction: column;
  }

  input,
  textarea {
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
    transition: background-color 0.3s ease;
  }

  .cancel-btn {
    background-color: #e8ecee;
    color: #2c3e50;
  }

  .cancel-btn:hover {
    background-color: #d1d8e0;
  }

  .submit-btn {
    background-color: #7b68ee;
    color: white;
  }

  .submit-btn:hover {
    background-color: #6c5ce7;
  }

  .expand-icon {
    font-size: 0.8rem;
    color: #7b68ee;
  }

  .tasks-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: #e8ecee;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    user-select: none;
    margin-top: 1rem;
  }

  .tasks-header h3 {
    margin: 0;
    font-size: 1rem;
    color: #2c3e50;
  }

  .tasks-list {
    margin-top: 0.5rem;
  }

  .limited-size {
    max-width: 100%;
    max-height: 200px;
    width: 100%;
    box-sizing: border-box;
  }

  @media (max-width: 600px) {
    .search-filter-sort {
      flex-direction: column;
    }

    .search-input,
    .filter-select,
    .sort-select,
    .sort-order-btn {
      width: 100%;
    }
  }
</style>
