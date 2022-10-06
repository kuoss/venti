import { useAuthStore } from "@/stores/auth";

export default {
  install: ({ config }) => {
    config.globalProperties.$auth = useAuthStore()
    // if (useAuthStore().loggedIn) {
    //   useAuthStore().fetchUser()
    // }
  },
}