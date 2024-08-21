<script>
  import { onMount } from "svelte";
  import TodoItem from "./components/TodoItem.svelte";

  import { API_BASE_URL } from "./config";

  let todos = [];

  let filterStatus = "";
  let sortBy = "";
  let sortOrder = "asc";

  let newTodo = {
    name: "",
    description: "",
    due_date: "",
  };

  onMount(async () => {
    await fetchTodos();
  });

  $: filteredTodos = todos.filter((todo) => {
    return filterStatus ? todo.status === filterStatus : true;
  });

  $: sortedTodos = filteredTodos.sort((a, b) => {
    if (sortBy) {
      if (sortOrder === "asc") {
        return a[sortBy] > b[sortBy] ? 1 : -1;
      } else {
        return a[sortBy] < b[sortBy] ? 1 : -1;
      }
    }
    return 0;
  });

  async function fetchTodos(statuses = [], sortBy = "", order = "") {
    let queryParams = [];

    if (statuses.length) {
      queryParams.push(`status=${statuses.join(",")}`);
    }
    if (sortBy) {
      queryParams.push(`sort_by=${sortBy}`);
    }
    if (order) {
      queryParams.push(`order=${order}`);
    }

    const queryString = queryParams.length ? `?${queryParams.join("&")}` : "";
    const response = await fetch(`${API_BASE_URL}/todos${queryString}`);
    todos = await response.json();
  }

  async function addTodo() {
    const response = await fetch(`${API_BASE_URL}/todos`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(newTodo),
    });
    if (response.ok) {
      await fetchTodos();
      newTodo = { name: "", description: "", due_date: "" };
    } else if (response.status === 400) {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  async function updateTodo(todo) {
    const response = await fetch(`${API_BASE_URL}/todos/${todo.ID}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(todo),
    });
    if (response.ok) {
      await fetchTodos();
    }
  }

  async function deleteTodo(id) {
    const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
      method: "DELETE",
    });
    if (response.ok) {
      await fetchTodos();
    }
  }

  function handleUpdate(event) {
    updateTodo(event.detail);
  }

  function handleDelete(id) {
    deleteTodo(id);
  }

  function applyFilters(status, sort, order) {
    fetchTodos(status ? [status] : [], sort, order);
  }
</script>

<main>
  <h1>Todo App</h1>
  <form on:submit|preventDefault={addTodo}>
    <input
      type="text"
      name="name"
      id="name"
      class="todo-input"
      placeholder="Enter todo name"
      bind:value={newTodo.name}
    />
    <input
      type="text"
      name="description"
      id="description"
      class="todo-input"
      placeholder="Enter todo description"
      bind:value={newTodo.description}
    />
    <input
      type="date"
      name="due_date"
      id="due_date"
      class="todo-input"
      placeholder="Enter due date"
      bind:value={newTodo.due_date}
    />
    <button type="submit">Add Todo</button>
  </form>
  <div>
    <label for="status">Filter by status:</label>
    <select id="status" bind:value={filterStatus}>
      <option value="">All</option>
      <option value="completed">Completed</option>
      <option value="pending">Pending</option>
    </select>

    <label for="sort">Sort by:</label>
    <select id="sort" bind:value={sortBy}>
      <option value="">None</option>
      <option value="name">name</option>
      <option value="due_date">Due Date</option>
    </select>

    <label for="order">Order:</label>
    <select id="order" bind:value={sortOrder}>
      <option value="asc">Ascending</option>
      <option value="desc">Descending</option>
    </select>

    <button on:click={() => applyFilters(filterStatus, sortBy, sortOrder)}
      >Apply</button
    >
  </div>

  {#each sortedTodos as todo (todo.ID)}
    <TodoItem
      {todo}
      on:update={handleUpdate}
      on:delete={() => handleDelete(todo.ID)}
    />
  {/each}
</main>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .todo-input {
    padding: 0.5rem;
    font-size: 1rem;
  }

  button {
    padding: 0.5rem;
    font-size: 1rem;
    cursor: pointer;
  }

  @media (min-width: 600px) {
    form {
      flex-direction: row;
      flex-wrap: wrap;
      gap: 0.5rem;
    }

    .todo-input {
      flex: 1 1 calc(33.333% - 1rem);
    }

    button {
      flex: 1 1 100%;
    }
  }

  @media (max-width: 599px) {
    form {
      flex-direction: column;
    }

    .todo-input {
      flex: 1 1 100%;
    }

    button {
      flex: 1 1 100%;
    }
  }
</style>
