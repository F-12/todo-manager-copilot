<template>
  <div>
    <h2>Todo List</h2>
    <input v-model="newTitle" placeholder="Add new todo" @keyup.enter="addTodo" />
    <button @click="addTodo">Add</button>
    <ul>
      <li v-for="todo in todos" :key="todo.id">
        <span
          :style="{ textDecoration: todo.completed ? 'line-through' : 'none', cursor: 'pointer' }"
          @click="toggleTodo(todo.id)"
        >
          {{ todo.title }}
        </span>
        <button @click="deleteTodo(todo.id)">Delete</button>
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';

interface Todo {
  id: number;
  title: string;
  completed: boolean;
}

const API_URL = '/api/todos';
const todos = ref<Todo[]>([]);
const newTitle = ref('');

const fetchTodos = async () => {
  const res = await fetch(API_URL);
  todos.value = await res.json();
};

const addTodo = async () => {
  if (!newTitle.value.trim()) return;
  const res = await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title: newTitle.value })
  });
  const todo = await res.json();
  todos.value.push(todo);
  newTitle.value = '';
};

const toggleTodo = async (id: number) => {
  const res = await fetch(`${API_URL}/${id}`, {
    method: 'PATCH',
  });
  const updated = await res.json();
  todos.value = todos.value.map(t => t.id === id ? updated : t);
};

const deleteTodo = async (id: number) => {
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
  todos.value = todos.value.filter(t => t.id !== id);
};

onMounted(fetchTodos);
</script>

<style scoped>
ul {
  padding-left: 0;
}
li {
  list-style: none;
  margin-bottom: 8px;
}
button {
  margin-left: 8px;
}
</style>
