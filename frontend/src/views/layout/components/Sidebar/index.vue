<template>
  <el-scrollbar wrap-class="scrollbar-wrapper">
    <div class="sidebar-logo" :class="isCollapse ? 'collapsed' : ''">
      <span>C</span
      ><span v-show="!isCollapse"
        >rawlab<span class="version">v{{ version }}</span>
        <span class="text-lite">Lite</span>
      </span>
    </div>
    <el-menu
      :show-timeout="200"
      :default-active="routeLevel1"
      :collapse="isCollapse"
      :background-color="variables.menuBg"
      :text-color="variables.menuText"
      :active-text-color="variables.menuActiveText"
      mode="vertical"
    >
      <sidebar-item
        v-for="route in routes"
        :key="route.path"
        :class="route.path.replace('/', '')"
        :item="route"
        :base-path="route.path"
      />
    </el-menu>
  </el-scrollbar>
</template>

<script>
import * as Vue from 'vue'
import { mapGetters, mapState } from 'vuex'
import variables from '@/styles/variables.scss'
import SidebarItem from './SidebarItem'

export default {
  components: { SidebarItem },
  computed: {
    ...mapState('user', ['adminPaths']),
    ...mapGetters(['sidebar']),
    routeLevel1() {
      const pathArray = this.$route.path.split('/')
      return `/${pathArray[1]}`
    },
    routes() {
      return this.$router.options.routes.filter((d) => {
        return !this.adminPaths.includes(d.path)
      })
    },
    variables() {
      return variables
    },
    isCollapse() {
      return !this.sidebar.opened
    },
    version() {
      return (
        this.$store.state.version.version || window.sessionStorage.getItem('v')
      )
    },
  },
  async created() {},
  mounted() {},
}
</script>

<style>
#app .sidebar-container .el-menu {
  height: calc(100% - 50px);
}
.sidebar-container .sidebar-logo {
  height: 52px;
  display: flex;
  align-items: center;
  padding-left: 20px;
  color: #fff;
  background: rgb(48, 65, 86);
  font-size: 24px;
  font-weight: 600;
  font-family: 'Verdana', serif;
}
.sidebar-container .sidebar-logo.collapsed {
  padding-left: 8px;
}
.sidebar-container .sidebar-logo .version {
  margin-left: 5px;
  font-weight: normal;
  font-size: 10px;
}
.text-lite {
  font-size: 11px;
  position: absolute;
  top: 15px;
  right: 25px;
}
</style>
