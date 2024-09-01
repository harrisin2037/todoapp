import { writable } from 'svelte/store';

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
                console.log('Fetching users from:', `${import.meta.env.VITE_API_BASE_URL}/users`);
                const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/users`, {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem("token")}`,
                    },
                });
                console.log('Response status:', response.status);
                console.log('Response headers:', Object.fromEntries(response.headers));

                const text = await response.text();
                console.log('Raw response:', text);

                if (response.ok) {
                    try {
                        const data = JSON.parse(text);
                        console.log('Parsed data:', data);
                        if (Array.isArray(data.users)) {
                            set(data.users);
                        } else {
                            console.error('Unexpected data structure:', data);
                        }
                    } catch (parseError) {
                        console.error('Error parsing JSON:', parseError);
                    }
                } else {
                    console.error('Failed to load users. Status:', response.status);
                }
            } catch (error) {
                console.error('Error loading users:', error);
            }
        }
    };
}

export const userStore = createUserStore();