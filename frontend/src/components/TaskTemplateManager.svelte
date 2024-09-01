<script>
  import { onMount } from "svelte";
  import { API_BASE_URL } from "../config.js";

  let templates = [];
  let newTemplate = {
    name: "",
    description: "",
    fields: [],
  };
  let editingTemplate = null;

  onMount(fetchTemplates);

  async function fetchTemplates() {
    const token = localStorage.getItem("token");
    const response = await fetch(`${API_BASE_URL}/task-templates`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    if (response.ok) {
      templates = await response.json();
    } else {
      console.error("Failed to fetch task templates");
    }
  }

  async function createTemplate() {
    const token = localStorage.getItem("token");
    const response = await fetch(`${API_BASE_URL}/task-templates`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(newTemplate),
    });

    if (response.ok) {
      alert("Template created successfully");
      await fetchTemplates();
      newTemplate = { name: "", description: "", fields: [] };
    } else {
      alert("Failed to create template");
    }
  }

  async function updateTemplate() {
    const token = localStorage.getItem("token");
    const response = await fetch(
      `${API_BASE_URL}/task-templates/${editingTemplate.id}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(editingTemplate),
      }
    );

    if (response.ok) {
      alert("Template updated successfully");
      await fetchTemplates();
      editingTemplate = null;
    } else {
      alert("Failed to update template");
    }
  }

  async function deleteTemplate(id) {
    if (confirm("Are you sure you want to delete this template?")) {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE_URL}/task-templates/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        alert("Template deleted successfully");
        await fetchTemplates();
      } else {
        alert("Failed to delete template");
      }
    }
  }

  function addField(template) {
    template.fields = [
      ...template.fields,
      { name: "", type: "text", default: "" },
    ];
  }

  function removeField(template, index) {
    template.fields = template.fields.filter((_, i) => i !== index);
  }

  function startEditing(template) {
    editingTemplate = JSON.parse(JSON.stringify(template));
  }

  function cancelEditing() {
    editingTemplate = null;
  }
</script>

<main>
  <h2>Task Template Manager</h2>

  <section>
    <h3>Current Task Templates</h3>
    {#if templates.length === 0}
      <p>No templates available.</p>
    {:else}
      <ul>
        {#each templates as template (template.id)}
          <li class="template">
            <h4>{template.name}</h4>
            <p>{template.description}</p>
            <button on:click={() => startEditing(template)}>Edit</button>
            <button on:click={() => deleteTemplate(template.id)}>Delete</button>
          </li>
        {/each}
      </ul>
    {/if}
  </section>

  <section>
    <h3>Create New Template</h3>
    <form on:submit|preventDefault={createTemplate}>
      <div>
        <label for="templateName">Template Name:</label>
        <input
          id="templateName"
          bind:value={newTemplate.name}
          placeholder="Template Name"
          required
        />
      </div>
      <div>
        <label for="templateDescription">Description:</label>
        <textarea
          id="templateDescription"
          bind:value={newTemplate.description}
          placeholder="Description"
        ></textarea>
      </div>

      <h4>Fields</h4>
      {#each newTemplate.fields as field, index}
        <div class="field">
          <input bind:value={field.name} placeholder="Field Name" required />
          <select bind:value={field.type}>
            <option value="text">Text</option>
            <option value="number">Number</option>
            <option value="checkbox">Checkbox</option>
            <option value="list">List</option>
          </select>
          <input bind:value={field.default} placeholder="Default Value" />
          <button type="button" on:click={() => removeField(newTemplate, index)}
            >Remove</button
          >
        </div>
      {/each}
      <button type="button" on:click={() => addField(newTemplate)}
        >Add Field</button
      >

      <button type="submit">Create Template</button>
    </form>
  </section>

  {#if editingTemplate}
    <div class="modal">
      <div class="modal-content">
        <h3>Edit Template: {editingTemplate.name}</h3>
        <form on:submit|preventDefault={updateTemplate}>
          <div>
            <label for="editTemplateName">Template Name:</label>
            <input
              id="editTemplateName"
              bind:value={editingTemplate.name}
              required
            />
          </div>
          <div>
            <label for="editTemplateDescription">Description:</label>
            <textarea
              id="editTemplateDescription"
              bind:value={editingTemplate.description}
            ></textarea>
          </div>

          <h4>Fields</h4>
          {#each editingTemplate.fields as field, index}
            <div class="field">
              <input
                bind:value={field.name}
                placeholder="Field Name"
                required
              />
              <select bind:value={field.type}>
                <option value="text">Text</option>
                <option value="number">Number</option>
                <option value="checkbox">Checkbox</option>
                <option value="list">List</option>
              </select>
              <input bind:value={field.default} placeholder="Default Value" />
              <button
                type="button"
                on:click={() => removeField(editingTemplate, index)}
                >Remove</button
              >
            </div>
          {/each}
          <button type="button" on:click={() => addField(editingTemplate)}
            >Add Field</button
          >

          <div class="modal-actions">
            <button type="submit">Save Changes</button>
            <button type="button" on:click={cancelEditing}>Cancel</button>
          </div>
        </form>
      </div>
    </div>
  {/if}
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  h2,
  h3,
  h4 {
    margin-bottom: 10px;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .field {
    display: flex;
    gap: 10px;
    align-items: center;
  }

  .template {
    border: 1px solid #ccc;
    padding: 10px;
    margin-bottom: 10px;
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
    padding: 20px;
    border-radius: 5px;
    max-width: 80%;
    max-height: 80%;
    overflow-y: auto;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
  }

  button {
    cursor: pointer;
    padding: 5px 10px;
  }
</style>
