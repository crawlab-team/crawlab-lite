<template>
  <el-tree ref="documentation-tree" :data="docData" node-key="fullUrl">
    <template v-slot="{ node, data }">
      <span
        class="custom-tree-node"
        :class="[data.active ? 'active' : '', `level-${data.level}`]"
      >
        <template
          v-if="data.level === 1 && data.children && data.children.length"
        >
          <span>{{ node.label }}</span>
        </template>
        <template v-else>
          <span>
            <a
              :href="data.fullUrl"
              target="_blank"
              style="display: block"
              @click="onClickDocumentationLink"
            >
              {{ node.label }}
            </a>
          </span>
        </template>
      </span>
    </template>
  </el-tree>
</template>

<script>
import * as Vue from 'vue'
import { mapState } from 'vuex'

export default {
  name: 'Documentation',
  data() {
    return {
      data: [],
    }
  },
  computed: {
    ...mapState('doc', ['docData']),
    pathLv1() {
      if (this.$route.path === '/') return '/'
      const m = this.$route.path.match(/(^\/\w+)/)
      return m[1]
    },
    currentDoc() {
      // find current doc
      let currentDoc
      for (let i = 0; i < this.$utils.doc.docs.length; i++) {
        const doc = this.$utils.doc.docs[i]
        if (this.pathLv1 === doc.path) {
          currentDoc = doc
          break
        }
      }
      return currentDoc
    },
  },
  watch: {
    pathLv1() {
      this.update()
    },
  },
  async created() {},
  mounted() {
    this.update()
  },
  methods: {
    isActiveNode(d) {
      // check match
      if (!this.currentDoc) return false
      return !!d.url.match(this.currentDoc.pattern)
    },
    update() {
      // expand related documentation list
      setTimeout(() => {
        this.docData.forEach((d) => {
          // parent node
          const isActive = this.isActiveNode(d)
          const node = this.$refs['documentation-tree'].getNode(d)
          node.expanded = isActive
          d['active'] = isActive

          // child nodes
          d.children.forEach((c) => {
            const node = this.$refs['documentation-tree'].getNode(c)
            const isActive = this.isActiveNode(c)
            if (!node.parent.expanded && isActive) {
              node.parent.expanded = true
            }
            c['active'] = isActive
          })
        })
      }, 100)
    },
    async getDocumentationData() {
      // fetch api data
      await this.$store.dispatch('doc/getDocData')
    },
    onClickDocumentationLink() {
      this.$st.sendEv('全局', '点击右侧文档链接')
    },
  },
}
</script>

<style scoped>
.el-tree >>> .custom-tree-node.active {
  color: #409eff;
}
.el-tree >>> .custom-tree-node.level-1 {
  font-weight: bolder;
}
</style>
