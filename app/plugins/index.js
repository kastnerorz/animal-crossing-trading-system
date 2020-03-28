export const state = () => ({
    user: {},
    token: ''
})

export const mutations = {
    saveUser(state, user) {
        state.user = user
    },

    saveToken(state, token) {
        state.token = token
    }
}