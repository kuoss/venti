import { defineStore } from 'pinia'

export const useUIDStore = defineStore('uid', {
    state: () => ({
        lastUID: 1,
        currentUID: 0,
    })
})

