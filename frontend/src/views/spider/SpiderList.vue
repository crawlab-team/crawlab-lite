<template>
  <div class="app-container">
    <!--add dialog-->
    <el-dialog
      :title="t('Add Spider')"
      width="40%"
      v-model="addDialogVisible"
      :before-close="onAddDialogClose"
    >
      <el-form
        ref="addSpiderForm"
        :model="spiderForm"
        inline-message
        label-width="120px"
      >
        <el-form-item :label="t('Spider Name')" prop="name" required>
          <el-input
            id="spider-name"
            v-model="spiderForm.name"
            :placeholder="t('Spider Name')"
          />
        </el-form-item>
        <el-form-item :label="t('Description')" prop="description">
          <el-input
            id="spider-description"
            v-model="spiderForm.description"
            :placeholder="t('Description')"
          />
        </el-form-item>
        <!--        <el-form-item :label="t('Results')" prop="col">-->
        <!--          <el-input-->
        <!--            id="col"-->
        <!--            v-model="spiderForm.col"-->
        <!--            :placeholder="t('By default: ') + 'results_<spider_name>'"-->
        <!--          />-->
        <!--        </el-form-item>-->
        <el-form-item :label="t('Upload Zip File')" name="site">
          <el-upload
            :action="$request.baseUrl + '/spiders'"
            :before-upload="beforeUpload"
            :data="uploadForm"
            :file-list="fileList"
            :headers="{ Authorization: token }"
            :on-success="onUploadSuccess"
            :on-error="onUploadError"
          >
            <el-button id="upload" size="small" type="primary">
              <el-icon style="margin-right: 5px"
                ><font-awesome-icon icon="fa-solid fa-file-upload"
              /></el-icon>
              {{ t('Upload and Commit') }}
            </el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <el-alert :closable="false" style="margin-bottom: 10px" type="warning">
        <p>
          <i class="fa fa-exclamation-triangle" />
          {{
            t(
              'NOTE: It is best to zip your spider files from the PROJECT ROOT.',
            )
          }}
        </p>
      </el-alert>
    </el-dialog>
    <!--./add dialog-->

    <!--version list dialog-->
    <el-dialog
      v-model="versionListDialogVisible"
      :title="`${t('Spider Version List')} (${t('Spider')}: ${
        activeSpider ? activeSpider.name : ''
      })`"
      width="640px"
    >
      <el-upload
        :action="$request.baseUrl + '/spiders/' + activeSpider.id + '/versions'"
        :data="uploadForm"
        :file-list="fileList"
        :headers="{ Authorization: token }"
        :on-success="onUploadVersionSuccess"
        :on-error="onUploadError"
        align="right"
      >
        <el-button size="small" type="success">
          <el-icon style="margin-right: 5px"
            ><font-awesome-icon icon="fa-solid fa-file-upload"
          /></el-icon>
          {{ t('Upload') }}
        </el-button>
      </el-upload>
      <el-table
        :data="spiderVersionList"
        border
        class="table"
        max-height="240px"
        style="margin: 5px 10px"
      >
        <el-table-column :label="t('Upload Time')" prop="create_ts">
          <template v-slot="scope">
            {{
              `${getTime(scope.row.create_ts)}${
                scope.$index === 0 ? ' (' + t('Latest') + ')' : ''
              }`
            }}
          </template>
        </el-table-column>
        <el-table-column :label="t('MD5')" prop="md5" />
        <el-table-column :label="t('Action')" width="120px">
          <template v-slot="scope">
            <el-tooltip :content="t('Remove')" placement="top">
              <el-button
                type="danger"
                size="small"
                @click="onRemoveVersion(scope.row, $event)"
              >
                <el-icon :size="10"
                  ><font-awesome-icon icon="fa-solid fa-trash" /></el-icon
              ></el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <template v-slot:footer>
        <el-button
          size="small"
          type="primary"
          @click="versionListDialogVisible = false"
          >{{ t('Ok') }}</el-button
        >
      </template>
    </el-dialog>
    <!--./version list dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="activeSpider.id"
      @close="onCrawlConfirmDialogClose"
      @confirm="onCrawlConfirm"
    />
    <!--./crawl confirm dialog-->

    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <!--          <el-form :inline="true">-->
          <!--            <el-form-item>-->
          <!--              <el-input-->
          <!--                v-model="filter.keyword"-->
          <!--                size="small"-->
          <!--                :placeholder="t('Spider Name')"-->
          <!--                clearable-->
          <!--                @keyup.enter.native="onSearch"-->
          <!--              >-->
          <!--                <i slot="suffix" class="el-input__icon el-icon-search"></i>-->
          <!--              </el-input>-->
          <!--            </el-form-item>-->
          <!--            <el-form-item>-->
          <!--              <el-button size="small" type="success"-->
          <!--                         class="btn refresh"-->
          <!--                         @click="onRefresh">-->
          <!--                {{t('Search')}}-->
          <!--              </el-button>-->
          <!--            </el-form-item>-->
          <!--          </el-form>-->
        </div>
        <div class="right">
          <el-button
            size="small"
            type="success"
            style="font-weight: bolder"
            @click="onAdd"
          >
            <el-icon style="margin-right: 5px"
              ><font-awesome-icon icon="fa-solid fa-plus"
            /></el-icon>
            {{ t('Add Spider') }}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table
        ref="table"
        :data="spiderList"
        class="table"
        :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
        row-key="id"
        border
      >
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'type'"
            :property="col.name"
            :label="t(col.label)"
            :sortable="col.sortable"
            align="left"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ t(scope.row.type) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_5_errors'"
            :label="t(col.label)"
            :width="col.width"
            align="center"
          >
            <template v-slot="scope">
              <div :style="{ color: scope.row[col.name] > 0 ? 'red' : '' }">
                {{ scope.row[col.name] }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'cmd'"
            :label="t(col.label)"
            :width="col.width"
            align="left"
          >
            <template v-slot="scope">
              <el-input v-model="scope.row[col.name]" />
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name.match(/_ts$/)"
            :label="t(col.label)"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ getTime(scope.row[col.name]) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_status'"
            :label="t(col.label)"
            align="left"
            :width="col.width"
          >
            <template v-slot="scope">
              <template v-if="scope.row.last_status === 'ERROR'">
                <el-tooltip :content="scope.row.last_error" placement="top">
                  <status-tag :status="scope.row.last_status" />
                </el-tooltip>
              </template>
              <status-tag v-else :status="scope.row.last_status" />
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :property="col.name"
            :label="t(col.label)"
            :align="col.align || 'left'"
            :width="col.width"
          />
        </template>
        <el-table-column
          :label="t('Action')"
          align="left"
          fixed="right"
          width="130"
        >
          <template v-slot="scope">
            <!--            <el-tooltip :content="t('View')" placement="top">-->
            <!--              <el-button-->
            <!--                type="primary"-->
            <!--                icon="el-icon-search"-->
            <!--                size="mini"-->
            <!--                @click="onView(scope.row, $event)"-->
            <!--              />-->
            <!--            </el-tooltip>-->
            <el-tooltip :content="t('Run')" placement="top">
              <el-button
                type="success"
                size="small"
                @click="onCrawl(scope.row, $event)"
              >
                <el-icon :size="10"
                  ><font-awesome-icon icon="fa-solid fa-bug"
                /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip :content="t('Spider Version List')" placement="top">
              <el-button
                type="warning"
                size="small"
                @click="onViewSpiderVersions(scope.row, $event)"
              >
                <el-icon :size="10"
                  ><font-awesome-icon icon="fa-solid fa-archive"
                /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip :content="t('Remove')" placement="top">
              <el-button
                type="danger"
                size="small"
                @click="onRemove(scope.row, $event)"
              >
                <el-icon :size="10"
                  ><font-awesome-icon icon="fa-solid fa-trash" /></el-icon
              ></el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page_num"
          :page-sizes="[10, 20, 50, 100]"
          v-model:page-size="pagination.page_size"
          layout="sizes, prev, pager, next"
          :total="spiderTotal"
          @current-change="getSpiderList"
          @size-change="getSpiderList"
        />
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import { mapGetters, mapState } from 'vuex'
import dayjs from 'dayjs'
import CrawlConfirmDialog from '../../components/Dialog/CrawlConfirmDialog'
import StatusTag from '../../components/Status/StatusTag'
import { useI18n } from 'vue-i18n'

export default {
  name: 'SpiderList',
  components: {
    CrawlConfirmDialog,
    StatusTag,
  },
  setup(props) {
    const { t } = useI18n()
    return { t }
  },
  data() {
    return {
      pagination: {
        page_num: 1,
        page_size: 10,
      },
      importLoading: false,
      isEditMode: false,
      dialogVisible: false,
      addDialogVisible: false,
      crawlConfirmDialogVisible: false,
      versionListDialogVisible: false,
      activeSpider: {},
      types: [],
      spiderFormRules: {
        name: [
          { required: true, message: 'Required Field', trigger: 'change' },
        ],
      },
      fileList: [],
      refreshHandle: undefined,
      activeSpiderTaskStatus: 'RUNNING',
      isStopLoading: false,
      isRemoveLoading: false,
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm',
      'spiderList',
      'spiderTotal',
      'spiderVersionList',
      'spiderVersionTotal',
    ]),
    ...mapState('task', ['taskList']),
    ...mapGetters('user', ['token']),
    ...mapState('lang', ['lang']),
    uploadForm() {
      return {
        name: this.spiderForm.name,
        description: this.spiderForm.description || '',
      }
    },
    columns() {
      const columns = []
      columns.push({
        name: 'name',
        label: 'Name',
        width: '160',
        align: 'left',
      })
      // columns.push({ name: 'latest_tasks', label: 'Latest Tasks', width: '180' })
      columns.push({ name: 'last_status', label: 'Last Status', width: '120' })
      columns.push({ name: 'last_run_ts', label: 'Last Run', width: '140' })
      columns.push({ name: 'update_ts', label: 'Update Time', width: '140' })
      columns.push({ name: 'create_ts', label: 'Create Time', width: '140' })
      columns.push({ name: 'description', label: 'Description' })
      return columns
    },
  },
  async created() {
    // fetch spider list
    await this.getSpiderList()

    // periodically fetch spider list
    this.refreshHandle = setInterval(() => {
      this.getSpiderList()
    }, 15000)
  },
  mounted() {
    const vm = this
    this.$nextTick(() => {
      vm.$store.commit('spider/SET_SPIDER_FORM', this.spiderForm)
    })
  },
  unmounted() {
    clearInterval(this.refreshHandle)
  },
  methods: {
    onAdd() {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.addDialogVisible = true
    },
    onRefresh() {
      this.getSpiderList()
      this.$st.sendEv('爬虫列表', '刷新')
    },
    onCancel() {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.dialogVisible = false
    },
    onDialogClose() {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.dialogVisible = false
    },
    onAddDialogClose() {
      this.addDialogVisible = false
    },
    onEdit(row) {
      this.isEditMode = true
      this.$store.commit('spider/SET_SPIDER_FORM', row)
      this.dialogVisible = true
    },
    onRemove(row, ev) {
      ev.stopPropagation()
      this.$prompt(
        this.t(
          'This action cannot be undone. This will permanently delete the spider and all associated data!',
        ) +
          `<br>${
            this.lang === 'zh'
              ? '请输入 <b>' + row.name + '</b> 确认删除：'
              : 'Please type' + row.name + 'to confirm:'
          }`,
        this.t('Notification'),
        {
          dangerouslyUseHTMLString: true,
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          inputPattern: new RegExp(`^${row.name}$`),
          inputErrorMessage: this.t('Inconsistent spider name'),
        },
      )
        .then(() => {
          this.$store
            .dispatch('spider/deleteSpider', row.id)
            .then(async (response) => {
              if (!response.data || response.data.code !== 200) {
                this.$message.error(response.data.message)
                return
              }
              this.$message({
                type: 'success',
                message: this.t('Deleted successfully'),
              })
              await this.getSpiderList()
              this.$st.sendEv('爬虫列表', '删除爬虫')
            })
        })
        .catch(() => {})
    },
    onCrawl(row, ev) {
      ev.stopPropagation()
      this.crawlConfirmDialogVisible = true
      this.activeSpider = row
      this.$st.sendEv('爬虫列表', '点击运行')
    },
    onCrawlConfirm() {
      setTimeout(() => {
        this.getSpiderList()
      }, 1000)
    },
    onRemoveVersion(row, ev) {
      ev.stopPropagation()
      this.$confirm(
        this.t('Are you sure to delete this version?'),
        this.t('Notification'),
        {
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          this.$store
            .dispatch('spider/deleteSpiderVersion', {
              spider_id: this.activeSpider.id,
              version_id: row.id,
            })
            .then(async (response) => {
              if (!response.data || response.data.code !== 200) {
                this.$message.error(response.data.message)
                return
              }
              this.$message({
                type: 'success',
                message: this.t('Deleted successfully'),
              })
              await this.getSpiderVersionList(this.activeSpider.id)
              this.$st.sendEv('爬虫版本列表', '删除爬虫版本')
            })
        })
        .catch(() => {})
    },
    onView(row, ev) {
      ev.stopPropagation()
      this.$router.push('/spiders/' + row.id)
      this.$st.sendEv('爬虫列表', '查看爬虫')
    },
    onUploadSuccess(res) {
      // clear fileList
      this.fileList = []

      // fetch spider list
      setTimeout(() => {
        this.getSpiderList()
      }, 500)

      this.$message.success(this.t('Uploaded spider files successfully'))
      this.addDialogVisible = false
    },
    onUploadVersionSuccess(res) {
      // clear fileList
      this.fileList = []

      // fetch spider list
      setTimeout(() => {
        this.getSpiderVersionList(this.activeSpider.id)
      }, 500)

      this.$message.success(this.t('Uploaded spider files successfully'))
    },
    onUploadError() {
      this.$message.error(this.t('Failed to upload spider files'))
    },
    beforeUpload(file) {
      return new Promise((resolve, reject) => {
        this.$refs['addSpiderForm'].validate((res) => {
          if (res) {
            resolve()
          } else {
            reject(new Error('form validation error'))
          }
        })
      })
    },
    getTime(str) {
      if (!str || str.match('^0001')) return 'NA'
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss')
    },
    onRowClick(row, column, event) {
      this.onView(row, event)
    },
    async getSpiderList() {
      await this.$store.dispatch('spider/getSpiderList', this.pagination)
    },
    async getSpiderVersionList(spiderId) {
      await this.$store.dispatch('spider/getSpiderVersionList', {
        spider_id: spiderId,
      })
    },
    async onViewSpiderVersions(row, ev) {
      ev.stopPropagation()
      this.activeSpider = row
      this.versionListDialogVisible = true
      await this.getSpiderVersionList(row.id)
    },
    async onStop(row, ev) {
      ev.stopPropagation()
      const res = await this.$store.dispatch('task/cancelTask', row.id)
      if (res.data && res.data.code === 200) {
        this.$message.success(`Task "${row.id}" has been sent signal to stop`)
        await this.getSpiderList()
      }
    },
    onCrawlConfirmDialogClose() {
      this.crawlConfirmDialogVisible = false
    },
  },
}
</script>

<style lang="scss" scoped>
.el-dialog {
  .el-select {
    width: 100%;
  }
}
.filter {
  display: flex;
  justify-content: space-between;
  .filter-search {
    width: 240px;
  }

  .right {
    .btn {
      margin-left: 10px;
    }
  }
}
.table {
  margin-top: 8px;
  border-radius: 5px;
  .el-button {
    padding: 7px;
  }
}
.delete-confirm {
  background-color: red;
}
.add-spider-wrapper {
  display: flex;
  justify-content: center;
  .add-spider-item {
    cursor: pointer;
    width: 180px;
    font-size: 18px;
    height: 120px;
    margin: 0 20px;
    flex-basis: 40%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .add-spider-item.primary {
    color: #409eff;
    background: rgba(64, 158, 255, 0.1);
    border: 1px solid rgba(64, 158, 255, 0.1);
  }

  .add-spider-item.success {
    color: #67c23a;
    background: rgba(103, 194, 58, 0.1);
    border: 1px solid rgba(103, 194, 58, 0.1);
  }

  .add-spider-item.info {
    color: #909399;
    background: #f4f4f5;
    border: 1px solid #e9e9eb;
  }
}
.el-autocomplete {
  width: 100%;
}
</style>

<style scoped>
.el-table >>> tr {
  cursor: pointer;
}

.el-table >>> .latest-tasks .el-tag {
  margin: 3px 3px 0 0;
}
</style>
