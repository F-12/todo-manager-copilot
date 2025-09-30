import React, { useState, useEffect } from 'react';

export interface Todo {
  id: number;
  title: string;
  completed: boolean;
}

const API_URL = '/api/todos';

const TodoList: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTitle, setNewTitle] = useState('');

  useEffect(() => {
    fetch(API_URL)
      .then(res => res.json())
      .then(setTodos);
  }, []);

  const addTodo = async () => {
    if (!newTitle.trim()) return;
    const res = await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTitle })
    });
    const todo = await res.json();
    setTodos([...todos, todo]);
    setNewTitle('');
  };

  const toggleTodo = async (id: number) => {
    const res = await fetch(`${API_URL}/${id}`, {
      method: 'PATCH',
    });
    const updated = await res.json();
    setTodos(todos.map(t => t.id === id ? updated : t));
  };

  const deleteTodo = async (id: number) => {
    await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
    setTodos(todos.filter(t => t.id !== id));
  };

  return (
    <div>
      <h2>Todo List</h2>
      <input
        value={newTitle}
        onChange={e => setNewTitle(e.target.value)}
        placeholder="Add new todo"
      />
      <button onClick={addTodo}>Add</button>
      <ul>
        {todos.map(todo => (
          <li key={todo.id}>
            <span
              style={{ textDecoration: todo.completed ? 'line-through' : 'none', cursor: 'pointer' }}
              onClick={() => toggleTodo(todo.id)}
            >
              {todo.title}
            </span>
            <button onClick={() => deleteTodo(todo.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TodoList;
