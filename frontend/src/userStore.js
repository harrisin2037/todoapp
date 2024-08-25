import { writable } from 'svelte/store';
import { API_BASE_URL } from './config';

function createUserStore() {
    const { subscribe, set, update } = writable([]);

    return {
        subscribe,
        setUsers: (users) => set(users),
        addUser: (user) => update(users => [...users, user]),
        updateUser: (id, updatedUser) => update(users => users.map(user => user.id === id ? { ...user, ...updatedUser } : user)),
        removeUser: (id) => update(users => users.filter(user => user.id !== id)),
        loadUsers: async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/users`, {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem("token")}`,
                    },
                });
                if (response.ok) {
                    const data = await response.json();
                    console.log('Loaded users:', data.users);
                    set(data.users);
                } else {
                    console.error('Failed to load users');
                }
            } catch (error) {
                console.error('Error loading users:', error);
            }
        }
    };
}

export const userStore = createUserStore();