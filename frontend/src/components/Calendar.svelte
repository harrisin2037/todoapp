<script>
  import { onMount } from "svelte";
  import TaskPopup from "./TaskPopUp.svelte";
  import { API_BASE_URL } from "../config";

  let todos = [];
  let currentDate = new Date();
  let currentMonth = currentDate.getMonth();
  let currentYear = currentDate.getFullYear();
  let selectedDate = null;

  onMount(async () => {
    await fetchTodos();
  });

  async function fetchTodos() {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE_URL}/todos`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (!response.ok) throw new Error("Failed to fetch todos");
      todos = await response.json();
      console.log("Fetched todos:", todos);
    } catch (error) {
      console.error("Error fetching todos:", error);
    }
  }

  function previousMonth() {
    if (currentMonth === 0) {
      currentMonth = 11;
      currentYear--;
    } else {
      currentMonth--;
    }
  }

  function nextMonth() {
    if (currentMonth === 11) {
      currentMonth = 0;
      currentYear++;
    } else {
      currentMonth++;
    }
  }

  function getDaysInMonth(month, year) {
    return new Date(year, month + 1, 0).getDate();
  }

  function getFirstDayOfMonth(month, year) {
    return new Date(year, month, 1).getDay();
  }

  function getFilteredTodos(date) {
    const filteredTodos = todos.filter((todo) => {
      const todoDueDate = new Date(todo.due_date);
      return (
        todoDueDate.getFullYear() === currentYear &&
        todoDueDate.getMonth() === currentMonth &&
        todoDueDate.getDate() === date
      );
    });
    console.log(
      `Filtered todos for ${currentYear}-${currentMonth + 1}-${date}:`,
      filteredTodos
    );
    return filteredTodos;
  }

  function openTaskPopup(date) {
    selectedDate = date;
  }

  function closeTaskPopup() {
    selectedDate = null;
  }

  async function toggleTodoCompletion(todoId) {
    try {
      const todo = todos.find((t) => t.id === todoId);
      const updatedStatus =
        todo.status === "completed" ? "pending" : "completed";

      const response = await fetch(`${API_BASE_URL}/todos/${todoId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ status: updatedStatus }),
      });

      if (!response.ok) throw new Error("Failed to update todo");

      todos = todos.map((t) =>
        t.id === todoId ? { ...t, status: updatedStatus } : t
      );
    } catch (error) {
      console.error("Error updating todo:", error);
    }
  }

  $: monthName = new Date(currentYear, currentMonth).toLocaleString("default", {
    month: "long",
  });

  $: firstDayOfMonth = getFirstDayOfMonth(currentMonth, currentYear);
  $: daysInMonth = getDaysInMonth(currentMonth, currentYear);
  $: calendarDays = [
    ...Array(firstDayOfMonth).fill(null),
    ...Array(daysInMonth).keys(),
  ].map((n) => (n === null ? null : n + 1));
  $: {
    if (todos) {
      calendarDays = [
        ...Array(firstDayOfMonth).fill(null),
        ...Array(daysInMonth).keys(),
      ].map((n) => (n === null ? null : n + 1));
    }
  }
</script>

<div class="calendar">
  <div class="calendar-header">
    <button on:click={previousMonth} class="nav-button">
      <svg viewBox="0 0 24 24" width="24" height="24">
        <path d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z" />
      </svg>
    </button>
    <h2>{monthName} {currentYear}</h2>
    <button on:click={nextMonth} class="nav-button">
      <svg viewBox="0 0 24 24" width="24" height="24">
        <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z" />
      </svg>
    </button>
  </div>

  <div class="calendar-grid">
    {#each ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"] as day}
      <div class="day-header">{day}</div>
    {/each}

    {#each calendarDays as date}
      {#if date !== null}
        {@const dayTodos = getFilteredTodos(date)}
        {@const isToday =
          date === currentDate.getDate() &&
          currentMonth === currentDate.getMonth() &&
          currentYear === currentDate.getFullYear()}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
          class="day {isToday ? 'today' : ''}"
          on:click={() => openTaskPopup(date)}
        >
          <span class="date">{date}</span>
          {#if dayTodos.length > 0}
            <div class="task-indicators">
              {#each dayTodos.slice(0, 3) as todo}
                <span
                  class="task-indicator"
                  style="background-color: {todo.status === 'completed'
                    ? '#4caf50'
                    : '#ff9800'};"
                ></span>
              {/each}
              {#if dayTodos.length > 3}
                <span class="task-indicator more">+{dayTodos.length - 3}</span>
              {/if}
            </div>
          {/if}
        </div>
      {:else}
        <div class="day empty"></div>
      {/if}
    {/each}
  </div>

  {#if selectedDate !== null}
    <TaskPopup
      date={new Date(currentYear, currentMonth, selectedDate)}
      tasks={getFilteredTodos(selectedDate)}
      onClose={closeTaskPopup}
      onToggleComplete={toggleTodoCompletion}
    />
  {/if}
</div>

<style>
  .calendar {
    max-width: 1000px;
    margin: 0 auto;
    font-family: "Roboto", sans-serif;
    background-color: #ffffff;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    padding: 20px;
    margin-top: 20px;
  }

  .calendar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .calendar-header h2 {
    font-size: 24px;
    font-weight: 500;
    color: #2c3e50;
  }

  .nav-button {
    background: none;
    border: none;
    cursor: pointer;
    color: #7b68ee;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    transition: background-color 0.3s;
  }

  .nav-button:hover {
    background-color: #f0f0f0;
  }

  .nav-button svg {
    fill: currentColor;
  }

  .calendar-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 10px;
  }

  .day-header {
    text-align: center;
    font-weight: 500;
    padding: 10px;
    color: #7f8c8d;
    font-size: 14px;
  }

  .day {
    border: 1px solid #e8ecee;
    border-radius: 8px;
    padding: 10px;
    min-height: 100px;
    position: relative;
    cursor: pointer;
    transition:
      background-color 0.3s,
      box-shadow 0.3s;
  }

  .day:hover {
    background-color: #f9f9f9;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  }

  .day.empty {
    background-color: #f9f9f9;
    border: none;
  }

  .day.today {
    background-color: #e8f5e9;
    border-color: #4caf50;
  }

  .date {
    position: absolute;
    top: 5px;
    left: 5px;
    font-weight: 500;
    color: #2c3e50;
  }

  .task-indicators {
    position: absolute;
    bottom: 5px;
    left: 5px;
    right: 5px;
    display: flex;
    flex-wrap: wrap;
    gap: 3px;
  }

  .task-indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .task-indicator.more {
    background-color: #7f8c8d;
    color: white;
    font-size: 10px;
    width: auto;
    height: auto;
    padding: 2px 4px;
    border-radius: 10px;
  }
</style>
