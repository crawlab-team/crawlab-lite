<template>
  <section class="app-main">
    <router-view v-slot="{ Component }" :key="key">
      <transition name="fade-transform" mode="out-in">
        <component :is="Component"></component>
      </transition>
    </router-view>
  </section>
</template>

<script>
import * as Vue from 'vue'
export default {
  name: 'AppMain',
  computed: {
    key() {
      return this.$route.name !== undefined
        ? this.$route.name + +new Date()
        : this.$route + +new Date()
    },
  },
}
</script>

<style scoped>
.app-main {
  min-height: calc(100vh - 100px);
  position: relative;
  overflow: hidden;
}
</style>
