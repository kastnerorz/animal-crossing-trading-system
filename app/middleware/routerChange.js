export default function({ store, route, redirect }) {
  if (typeof store.state.loadingComponent === 'object') {
    store.commit('closeLoading')
  }
}
