<template>
  <div :class="classObj" class="app-wrapper">
    <div
      v-if="device === 'mobile' && sidebar.opened"
      class="drawer-bg"
      @click="handleClickOutside"
    />

    <!--sidebar-->
    <sidebar class="sidebar-container" />
    <!--./sidebar-->

    <!--main container-->
    <div class="main-container">
      <navbar />
      <tags-view />
      <app-main />
    </div>
    <!--./main container-->

    <!--documentation-->
    <div class="documentation">
      <!--      <el-tooltip-->
      <!--        :content="t('Click to view related Documentation')"-->
      <!--      >-->
      <!--        <i class="el-icon-question" @click="onClickDocumentation"></i>-->
      <!--      </el-tooltip>-->
      <el-drawer
        :title="t('Related Documentation')"
        v-model="isShowDocumentation"
        :before-close="onCloseDocumentation"
        size="300px"
      >
        <documentation />
      </el-drawer>
    </div>
    <!--./documentation-->
  </div>
</template>

<script>
import * as Vue from 'vue'
import { AppMain, Navbar, Sidebar, TagsView } from './components'
import ResizeMixin from './mixin/ResizeHandler'
import Documentation from '../../components/Documentation/Documentation'
import { useI18n } from 'vue-i18n'

export default {
  name: 'Layout',
  components: {
    Documentation,
    Navbar,
    Sidebar,
    TagsView,
    AppMain,
  },
  mixins: [ResizeMixin],
  data() {
    return {
      isShowDocumentation: false,
    }
  },
  computed: {
    sidebar() {
      return this.$store.state.app.sidebar
    },
    device() {
      return this.$store.state.app.device
    },
    classObj() {
      return {
        hideSidebar: !this.sidebar.opened,
        openSidebar: this.sidebar.opened,
        withoutAnimation: this.sidebar.withoutAnimation,
        mobile: this.device === 'mobile',
      }
    },
  },
  async created() {
    // await this.$store.dispatch('doc/getDocData')
  },
  methods: {
    handleClickOutside() {
      this.$store.dispatch('CloseSideBar', { withoutAnimation: false })
    },
    onClickDocumentation() {
      this.isShowDocumentation = true
      this.$st.sendEv('全局', '打开右侧文档')
    },
    onCloseDocumentation() {
      this.isShowDocumentation = false
      this.$st.sendEv('全局', '关闭右侧文档')
    },
  },
  setup(props) {
    const { t } = useI18n()
    return { t }
  },
}
</script>

<style lang="scss" rel="stylesheet/scss" scoped>
@import '../../../src/styles/mixin.scss';
.app-wrapper {
  @include clearfix;
  position: relative;
  height: 100%;
  width: 100%;
  background: white;
  &.mobile.openSidebar {
    position: fixed;
    top: 0;
  }
}
.drawer-bg {
  background: #000;
  opacity: 0.3;
  width: 100%;
  top: 0;
  height: 100%;
  position: absolute;
  z-index: 999;
}
.documentation {
  z-index: 9999;
  position: fixed;
  right: 25px;
  bottom: 20px;
  font-size: 32px;
  cursor: pointer;
  color: #909399;
}
</style>

<style scoped>
.documentation >>> .el-drawer__body {
  overflow: auto;
}

.documentation >>> span[role='heading']:focus {
  outline: none;
}

.documentation >>> .el-tree-node__content {
  height: 40px;
  line-height: 40px;
}

.documentation >>> .custom-tree-node {
  display: block;
  width: 100%;
  height: 40px;
  line-height: 40px;
  font-size: 14px;
}

.documentation >>> .custom-tree-node a {
  display: block;
}

.documentation >>> .custom-tree-node:hover a {
  text-decoration: underline;
}
</style>
