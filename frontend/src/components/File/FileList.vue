<template>
  <div class="file-list-container">
    <el-dialog
      :title="$t('New Directory')"
      v-model="dirDialogVisible"
      width="30%"
    >
      <el-form>
        <el-form-item :label="$t('Enter new directory name')">
          <el-input v-model="name" :placeholder="$t('New directory name')" />
        </el-form-item>
      </el-form>
      <template v-slot:footer>
        <span class="dialog-footer">
          <el-button @click="dirDialogVisible = false">{{
            $t('Cancel')
          }}</el-button>
          <el-button type="primary" @click="onAddDir">{{
            $t('Confirm')
          }}</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog :title="$t('New File')" v-model="fileDialogVisible" width="30%">
      <el-form>
        <el-form-item :label="$t('Enter new file name')">
          <el-input v-model="name" :placeholder="$t('New file name')" />
        </el-form-item>
      </el-form>
      <template v-slot:footer>
        <span class="dialog-footer">
          <el-button size="small" @click="fileDialogVisible = false">{{
            $t('Cancel')
          }}</el-button>
          <el-button size="small" type="primary" @click="onAddFile">{{
            $t('Confirm')
          }}</el-button>
        </span>
      </template>
    </el-dialog>

    <div class="file-tree-wrapper">
      <el-tree
        ref="tree"
        class="tree"
        :data="computedFileTree"
        node-key="path"
        :highlight-current="true"
        :default-expanded-keys="expandedPaths"
        @node-contextmenu="onFileRightClick"
        @node-click="onFileClick"
        @node-expand="onDirClick"
        @node-collapse="onDirClick"
      >
        <template v-slot="{ data }">
          <span class="custom-tree-node">
            <el-popover
              v-model="isShowCreatePopoverDict[data.path]"
              trigger="manual"
              placement="right"
              popper-class="create-item-popover"
              :visible-arrow="false"
              @hide="onHideCreate(data)"
            >
              <ul class="action-item-list">
                <li class="action-item" @click="fileDialogVisible = true">
                  <font-awesome-icon icon="file-alt" color="rgba(3,47,98,.5)" />
                  <span class="action-item-text">{{ $t('Create File') }}</span>
                </li>
                <li class="action-item" @click="dirDialogVisible = true">
                  <font-awesome-icon
                    :icon="['fa', 'folder']"
                    color="rgba(3,47,98,.5)"
                  />
                  <span class="action-item-text">{{
                    $t('Create Directory')
                  }}</span>
                </li>
              </ul>
              <ul class="action-item-list">
                <li class="action-item" @click="onClickRemoveNav(data)">
                  <font-awesome-icon
                    :icon="['fa', 'trash']"
                    color="rgba(3,47,98,.5)"
                  />
                  <span class="action-item-text">{{ $t('Remove') }}</span>
                </li>
              </ul>
              <template v-slot:reference>
                <div>
                  <span class="item-icon">
                    <font-awesome-icon
                      v-if="data.is_dir"
                      :icon="['fa', 'folder']"
                      color="rgba(3,47,98,.5)"
                    />
                    <font-awesome-icon
                      v-else-if="data.path.match(/\.py$/)"
                      :icon="['fab', 'python']"
                      color="rgba(3,47,98,.5)"
                    />
                    <font-awesome-icon
                      v-else-if="data.path.match(/\.js$/)"
                      :icon="['fab', 'node-js']"
                      color="rgba(3,47,98,.5)"
                    />
                    <font-awesome-icon
                      v-else-if="data.path.match(/\.(java|jar|class)$/)"
                      :icon="['fab', 'java']"
                      color="rgba(3,47,98,.5)"
                    />
                    <font-awesome-icon
                      v-else-if="data.path.match(/\.go$/)"
                      :icon="['fab', 'go']"
                      color="rgba(3,47,98,.5)"
                    />
                    <font-awesome-icon
                      v-else-if="data.path.match(/\.zip$/)"
                      :icon="['fa', 'file-archive']"
                      color="rgba(3,47,98,.5)"
                    />
                    <font-awesome-icon
                      v-else
                      icon="file-alt"
                      color="rgba(3,47,98,.5)"
                    />
                  </span>
                  <span
                    class="item-name"
                    :class="isActiveFile(data) ? 'active' : ''"
                  >
                    {{ data.name }}
                  </span>
                </div>
              </template>
            </el-popover>
          </span>
        </template>
      </el-tree>
      <div class="add-btn-wrapper">
        <el-popover
          trigger="click"
          placement="right"
          popper-class="create-item-popover"
          :visible-arrow="false"
        >
          <ul class="action-item-list">
            <li class="action-item" @click="fileDialogVisible = true">
              <font-awesome-icon icon="file-alt" color="rgba(3,47,98,.5)" />
              <span class="action-item-text">{{ $t('Create File') }}</span>
            </li>
            <li class="action-item" @click="dirDialogVisible = true">
              <font-awesome-icon
                :icon="['fa', 'folder']"
                color="rgba(3,47,98,.5)"
              />
              <span class="action-item-text">{{ $t('Create Directory') }}</span>
            </li>
          </ul>
          <template v-slot:reference>
            <el-button
              class="add-btn"
              size="small"
              type="primary"
              :icon="Plus"
              @click="onEmptyClick"
            >
              {{ $t('Add') }}
            </el-button>
          </template>
        </el-popover>
      </div>
    </div>

    <div class="main-content">
      <div v-if="!showFile" class="file-list">
        {{ $t('Please select a file or click the add button on the left.') }}
      </div>
      <template v-else>
        <div class="top-part">
          <!--back-->
          <div class="action-container">
            <el-popover v-model="isShowDelete" trigger="click">
              <el-button
                size="small"
                type="default"
                @click="() => (this.isShowDelete = false)"
              >
                {{ $t('Cancel') }}
              </el-button>
              <el-button size="small" type="danger" @click="onFileDelete">
                {{ $t('Confirm') }}
              </el-button>
              <template v-slot:reference>
                <el-button
                  type="danger"
                  size="small"
                  style="margin-right: 10px"
                  :disabled="isDisabled"
                >
                  <font-awesome-icon :icon="['fa', 'trash']" />
                  {{ $t('Remove') }}
                </el-button>
              </template>
            </el-popover>
            <el-popover v-model="isShowRename" trigger="click">
              <el-input
                v-model="name"
                :placeholder="$t('Name')"
                style="margin-bottom: 10px"
              />
              <div style="text-align: right">
                <el-button size="small" type="warning" @click="onRenameFile">
                  {{ $t('Confirm') }}
                </el-button>
              </div>
              <template v-slot:reference>
                <div>
                  <el-button
                    type="warning"
                    size="small"
                    style="margin-right: 10px"
                    :disabled="isDisabled"
                    @click="onOpenRename"
                  >
                    <font-awesome-icon :icon="['fa', 'redo']" />
                    {{ $t('Rename') }}
                  </el-button>
                </div>
              </template>
            </el-popover>
            <el-button
              type="success"
              size="small"
              style="margin-right: 10px"
              :disabled="isDisabled"
              @click="onFileSave"
            >
              <font-awesome-icon :icon="['fa', 'save']" />
              {{ $t('Save') }}
            </el-button>
          </div>
          <!--./back-->

          <!--file path-->
          <div class="file-path-container">
            <div class="file-path">{{ currentPath }}</div>
          </div>
          <!--./file path-->
        </div>
        <file-detail />
      </template>
    </div>
  </div>
</template>

<script>
import * as Vue from 'vue'
import { mapState } from 'vuex'
import FileDetail from './FileDetail'

export default {
  data() {
    return {
      isEdit: false,
      showFile: false,
      name: '',
      isShowAdd: false,
      isShowDelete: false,
      isShowRename: false,
      isShowCreatePopoverDict: {},
      currentFilePath: '.',
      ignoreFileRegexList: ['__pycache__', 'md5.txt', '.pyc', '.git'],
      activeFileNode: {},
      dirDialogVisible: false,
      fileDialogVisible: false,
      nodeExpandedDict: {},
      isShowDeleteNav: false,
      ElIconPlus,
    }
  },
  name: 'FileList',
  components: { FileDetail },
  computed: {
    ...mapState('spider', ['fileTree', 'spiderForm']),
    ...mapState('file', ['fileList']),
    currentPath: {
      set(value) {
        this.$store.commit('file/SET_CURRENT_PATH', value)
      },
      get() {
        return this.$store.state.file.currentPath
      },
    },
    computedFileTree() {
      if (!this.fileTree || !this.fileTree.children) return []
      let nodes = this.sortFiles(this.fileTree.children)
      nodes = this.filterFiles(nodes)
      return nodes
    },
    expandedPaths() {
      return Object.keys(this.nodeExpandedDict)
        .map((path) => {
          return {
            path,
            expanded: this.nodeExpandedDict[path],
          }
        })
        .filter((d) => d.expanded)
        .map((d) => d.path)
    },
  },
  async created() {
    await this.getFileTree()
  },
  mounted() {
    document.querySelector('body').addEventListener('click', (ev) => {
      this.isShowCreatePopoverDict = {}
    })
  },
  unmounted() {
    document.querySelector('body').removeEventListener('click')
  },
  methods: {
    onEdit() {
      this.isEdit = true
    },
    onItemClick(item) {
      if (item.is_dir) {
        // 目录
        this.$store.dispatch('file/getFileList', { path: item.path })
      } else {
        // 文件
        this.showFile = true
        this.$store.commit('file/SET_FILE_CONTENT', '')
        this.$store.commit('file/SET_CURRENT_PATH', item.path)
        this.$store.dispatch('file/getFileContent', { path: item.path })
      }
      this.$st.sendEv('爬虫详情', '文件', '点击')
    },
    async onFileSave() {
      await this.$store.dispatch('file/saveFileContent', {
        path: this.currentPath,
      })
      this.$message.success(this.$t('Saved file successfully'))
      this.$st.sendEv('爬虫详情', '文件', '保存')
    },
    async onAddFile() {
      if (!this.name) {
        this.$message.error(this.$t('Name cannot be empty'))
        return
      }
      const arr = this.activeFileNode.path.split('/')
      if (this.activeFileNode.is_dir) {
        arr.push(this.name)
      } else {
        arr[arr.length - 1] = this.name
      }
      const path = arr.join('/')
      await this.$store.dispatch('file/addFile', { path })
      await this.$store.dispatch('spider/getFileTree')
      this.isShowAdd = false
      this.fileDialogVisible = false
      this.showFile = true
      this.$store.commit('file/SET_FILE_CONTENT', '')
      this.$store.commit('file/SET_CURRENT_PATH', path)
      await this.$store.dispatch('file/getFileContent', { path })
      this.$st.sendEv('爬虫详情', '文件', '添加')
    },
    async onAddDir() {
      if (!this.name) {
        this.$message.error(this.$t('Name cannot be empty'))
        return
      }
      const arr = this.activeFileNode.path.split('/')
      if (this.activeFileNode.is_dir) {
        arr.push(this.name)
      } else {
        arr[arr.length - 1] = this.name
      }
      const path = arr.join('/')
      await this.$store.dispatch('file/addDir', { path })
      await this.$store.dispatch('spider/getFileTree')
      this.isShowAdd = false
      this.dirDialogVisible = false
      this.$st.sendEv('爬虫详情', '文件', '添加')
    },
    async onFileDelete() {
      await this.$store.dispatch('file/deleteFile', {
        path: this.currentFilePath,
      })
      await this.$store.dispatch('spider/getFileTree')
      this.$message.success(this.$t('Deleted successfully'))
      this.isShowDelete = false
      this.showFile = false
      this.$st.sendEv('爬虫详情', '文件', '删除')
    },
    onOpenRename() {
      const arr = this.currentFilePath.split('/')
      this.name = arr[arr.length - 1]
    },
    async onRenameFile() {
      await this.$store.dispatch('file/renameFile', {
        path: this.currentFilePath,
        newPath: this.name,
      })
      await this.$store.dispatch('spider/getFileTree')
      const arr = this.currentFilePath.split('/')
      arr[arr.length - 1] = this.name
      this.currentFilePath = arr.join('/')
      this.$store.commit('file/SET_CURRENT_PATH', this.currentFilePath)
      this.$message.success(this.$t('Renamed successfully'))
      this.isShowRename = false
      this.$st.sendEv('爬虫详情', '文件', '重命名')
    },
    async getFileTree() {
      const arr = this.$route.path.split('/')
      const id = arr[arr.length - 1]
      await this.$store.dispatch('spider/getFileTree', { id })
    },
    async onFileClick(data) {
      if (data.is_dir) {
        return
      }
      this.currentFilePath = data.path
      this.onItemClick(data)
    },
    onDirClick(data, node) {
      const vm = this
      setTimeout(() => {
        vm.$set(vm.nodeExpandedDict, data.path, node.expanded)
      }, 0)
    },
    sortFiles(nodes) {
      nodes.forEach((node) => {
        if (node.is_dir) {
          if (!node.children) node.children = []
          node.children = this.sortFiles(node.children)
        }
      })
      return nodes.sort((a, b) => {
        if ((a.is_dir && b.is_dir) || (!a.is_dir && !b.is_dir)) {
          return a.name > b.name ? 1 : -1
        } else {
          return a.is_dir ? -1 : 1
        }
      })
    },
    filterFiles(nodes) {
      return nodes.filter((node) => {
        if (node.is_dir) {
          node.children = this.filterFiles(node.children)
        }
        for (let i = 0; i < this.ignoreFileRegexList.length; i++) {
          const regex = this.ignoreFileRegexList[i]
          if (node.name.match(regex)) {
            return false
          }
        }
        return true
      })
    },
    isActiveFile(node) {
      return node.path === this.currentFilePath
    },
    onFileRightClick(ev, data) {
      this.isShowCreatePopoverDict = {}
      this.isShowCreatePopoverDict[data.path] = true
      this.activeFileNode = data
      this.$st.sendEv('爬虫详情', '文件', '右键点击导航栏')
    },
    onEmptyClick() {
      const data = { path: '' }
      this.isShowCreatePopoverDict = {}
      this.isShowCreatePopoverDict[data.path] = true
      this.activeFileNode = data
      this.$st.sendEv('爬虫详情', '文件', '空白点击添加')
    },
    onHideCreate(data) {
      this.isShowCreatePopoverDict[data.path] = false
      this.name = ''
    },
    onClickRemoveNav(data) {
      this.$confirm(
        this.$t('Are you sure to delete this file/directory?'),
        this.$t('Notification'),
        {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          confirmButtonClass: 'danger',
          type: 'warning',
        },
      )
        .then(() => {
          this.onFileDeleteNav(data.path)
        })
        .catch(() => {})
    },
    async onFileDeleteNav(path) {
      await this.$store.dispatch('file/deleteFile', { path })
      await this.$store.dispatch('spider/getFileTree')
      this.$message.success(this.$t('Deleted successfully'))
      this.isShowDelete = false
      this.showFile = false
      this.$st.sendEv('爬虫详情', '文件', '删除')
    },
    clickSpider(filepath) {
      const node = this.$refs['tree'].getNode(filepath)
      const data = node.data
      this.onFileClick(data)
      node.parent.expanded = true
      this.nodeExpandedDict[node.parent.data.path] = true
      node.parent.parent.expanded = true
      this.nodeExpandedDict[node.parent.parent.data.path] = true
    },
    clickPipeline() {
      const filename = 'pipelines.py'
      for (let i = 0; i < this.computedFileTree.length; i++) {
        const dataLv1 = this.computedFileTree[i]
        const nodeLv1 = this.$refs['tree'].getNode(dataLv1.path)
        if (dataLv1.is_dir) {
          for (let j = 0; j < dataLv1.children.length; j++) {
            const dataLv2 = dataLv1.children[j]
            if (dataLv2.path.match(filename)) {
              this.onFileClick(dataLv2)
              nodeLv1.expanded = true
              this.nodeExpandedDict[dataLv1.path] = true
              return
            }
          }
        }
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.file-list-container {
  margin-left: -15px;
  height: 100%;
  min-height: 100%;
  .top-part {
    display: flex;
    height: 33px;
    margin-bottom: 10px;

    .file-path-container {
      width: 100%;
      padding: 5px;
      margin: 0 10px 0 0;
      border-radius: 5px;
      border: 1px solid #eaecef;
      display: flex;
      justify-content: space-between;
      color: rgba(3, 47, 98, 1);

      .left {
        width: 100%;
        display: flex;

        .el-icon-back {
          margin-right: 10px;
          cursor: pointer;
        }

        .el-input {
          /*height: 22px;*/
          width: 100%;
          line-height: 10px;
        }
      }

      .el-icon-edit {
        cursor: pointer;
      }
    }

    .action-container {
      text-align: right;
      display: flex;
      /*padding: 1px 5px;*/
      /*height: 24px;*/

      .el-button {
        margin: 0;
      }
    }
  }

  .file-list {
    padding: 10px;
    list-style: none;
    height: 100%;
    overflow-y: auto;
    min-height: 100%;
    /*border-radius: 5px;*/
    /*border: 1px solid #eaecef;*/
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
  }
}
</style>

<style scoped>
.file-path >>> .el-input__inner {
  font-size: 14px;
  line-height: 18px;
  height: 18px;
  border-top: none;
  border-left: none;
  border-right: none;
  border-bottom: 2px solid #409eff;
  border-radius: 0;
}

.CodeMirror-line {
  padding-right: 20px;
}

.item {
  border-bottom: 1px solid #eaecef;
}

.item-icon {
  display: inline-block;
  width: 18px;
}

.item-name {
  font-size: 14px;
  color: rgba(3, 47, 98, 1);
}

.add-type-list {
  text-align: right;
  margin-top: 10px;
}

.add-type {
  cursor: pointer;
  font-weight: bolder;
}

.file-tree-wrapper {
  float: left;
  width: 240px;
  height: calc(100vh - 200px);
  overflow: auto;
}

.file-tree-wrapper >>> .el-tree-node__content {
  height: 30px;
}

.file-tree-wrapper >>> .el-tree-node__content .item-name.active {
  font-weight: bolder;
}

.main-content {
  float: left;
  width: calc(100% - 240px);
  height: calc(100vh - 200px);
  border-left: 1px solid #eaecef;
  padding-left: 20px;
}
</style>

<style>
.create-item-popover {
  padding: 0;
  margin: 0;
}

.action-item-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.action-item-title {
  padding-top: 10px;
  padding-left: 15px;
  padding-bottom: 5px;
}

.action-item-list .action-item {
  display: flex;
  align-items: center;
  height: 35px;
  padding: 0 0 0 10px;
  margin: 0;
  cursor: pointer;
}

.action-item-list .action-item:last-child {
  border-bottom: 1px solid #eaecef;
}

.action-item-list .action-item:hover {
  background: #f5f7fa;
}

.action-item-list .action-item svg {
  width: 20px;
}

.action-item-list .action-item .action-item-text {
  margin-left: 5px;
}

.add-btn-wrapper {
  width: 220px;
  border-top: 1px solid #eaecef;
  margin: 10px 10px;
}

.add-btn-wrapper .add-btn {
  width: 80px;
  margin-left: calc(120px - 40px - 10px);
  margin-top: 20px;
}
</style>
