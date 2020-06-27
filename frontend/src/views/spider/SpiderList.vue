<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="spider-list"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <v-tour
      name="spider-list-add"
      :steps="tourAddSteps"
      :callbacks="tourAddCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <!--add dialog-->
    <el-dialog
      :title="$t('Add Spider')"
      width="40%"
      :visible.sync="addDialogVisible"
      :before-close="onAddDialogClose"
    >
      <el-form ref="addSpiderForm" :model="spiderForm" inline-message label-width="120px">
        <el-form-item :label="$t('Spider Name')" prop="name" required>
          <el-input id="spider-name" v-model="spiderForm.name" :placeholder="$t('Spider Name')" />
        </el-form-item>
        <el-form-item :label="$t('Description')" prop="description">
          <el-input id="spider-description" v-model="spiderForm.description" :placeholder="$t('Description')" />
        </el-form-item>
        <!--        <el-form-item :label="$t('Results')" prop="col">-->
        <!--          <el-input-->
        <!--            id="col"-->
        <!--            v-model="spiderForm.col"-->
        <!--            :placeholder="$t('By default: ') + 'results_<spider_name>'"-->
        <!--          />-->
        <!--        </el-form-item>-->
        <el-form-item :label="$t('Upload Zip File')" label-width="120px" name="site">
          <el-upload
            :action="$request.baseUrl + '/spiders'"
            :before-upload="beforeUpload"
            :data="uploadForm"
            :file-list="fileList"
            :headers="{Authorization:token}"
            :on-success="onUploadSuccess"
          >
            <el-button id="upload" icon="el-icon-upload" size="small" type="primary">
              {{ $t('Upload') }}
            </el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <el-alert
        :closable="false"
        style="margin-bottom: 10px"
        type="warning"
      >
        <p>{{ $t('You can click "Add" to create an empty spider and upload files later.') }}</p>
        <p>{{ $t('OR, you can also click "Upload" and upload a zip file containing your spider project.') }}</p>
        <p>
          <i class="fa fa-exclamation-triangle" /> {{ $t('NOTE: When uploading a zip file, please zip your' +
            ' spider files from the ROOT DIRECTORY.') }}
        </p>
      </el-alert>
    </el-dialog>
    <!--./add dialog-->

    <!--running tasks dialog-->
    <el-dialog
      :visible.sync="isRunningTasksDialogVisible"
      :title="`${$t('Latest Tasks')} (${$t('Spider')}: ${activeSpider ? activeSpider.name : ''})`"
      width="920px"
    >
      <el-tabs v-model="activeSpiderTaskStatus">
        <el-tab-pane name="pending" :label="$t('Pending')" />
        <el-tab-pane name="running" :label="$t('Running')" />
        <el-tab-pane name="finished" :label="$t('Finished')" />
        <el-tab-pane name="error" :label="$t('Error')" />
        <el-tab-pane name="cancelled" :label="$t('Cancelled')" />
        <el-tab-pane name="abnormal" :label="$t('Abnormal')" />
      </el-tabs>
      <template slot-scope="scope">
        <el-table
          :data="getTasks(scope.row)"
          border
          class="table"
          max-height="240px"
          style="margin: 5px 10px"
        >
          <el-table-column
            :label="$t('Create Time')"
            prop="create_ts"
            width="140px"
          />
          <el-table-column
            :label="$t('Start Time')"
            prop="start_ts"
            width="140px"
          />
          <el-table-column
            :label="$t('Finish Time')"
            prop="finish_ts"
            width="140px"
          />
          <el-table-column
            :label="$t('Cmd')"
            prop="param"
            width="120px"
          />
          <el-table-column
            :label="$t('Status')"
            width="120px"
          >
            <template slot-scope="scope">
              <template
                v-if="scope.row.status === 'ERROR'"
              >
                <el-tooltip :content="scope.row.error" placement="top">
                  <status-tag :status="scope.row.status" />
                </el-tooltip>
              </template>
              <status-tag v-else :status="scope.row.status" />
            </template>
          </el-table-column>
          <el-table-column
            :label="$t('Results Count')"
            prop="result_count"
            width="80px"
          />
          <el-table-column
            :label="$t('Action')"
            width="auto"
          >
            <template slot-scope="scope">
              <el-button
                v-if="['pending', 'running'].includes(scope.row.status)"
                icon="el-icon-video-pause"
                size="mini"
                type="danger"
                @click="onStop(scope.row, $event)"
              />
            </template>
          </el-table-column>
        </el-table>
      </template>

      <template slot="footer">
        <el-button size="small" type="primary" @click="isRunningTasksDialogVisible = false">{{ $t('Ok') }}</el-button>
      </template>
    </el-dialog>
    <!--./running tasks dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="activeSpiderId"
      :spiders="selectedSpiders"
      :multiple="isMultiple"
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
          <!--                :placeholder="$t('Spider Name')"-->
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
          <!--                {{$t('Search')}}-->
          <!--              </el-button>-->
          <!--            </el-form-item>-->
          <!--          </el-form>-->
        </div>
        <div class="right">
          <el-button
            v-if="this.selectedSpiders.length"
            size="small"
            type="danger"
            icon="el-icon-video-play"
            class="btn add"
            style="font-weight: bolder"
            @click="onCrawlSelectedSpiders"
          >
            {{ $t('Run') }}
          </el-button>
          <el-button
            v-if="this.selectedSpiders.length"
            size="small"
            type="info"
            :icon="isStopLoading ? 'el-icon-loading' : 'el-icon-video-pause'"
            class="btn add"
            style="font-weight: bolder"
            @click="onStopSelectedSpiders"
          >
            {{ $t('Stop') }}
          </el-button>
          <el-button
            v-if="this.selectedSpiders.length"
            size="small"
            type="danger"
            :icon="isRemoveLoading ? 'el-icon-loading' : 'el-icon-delete'"
            class="btn add"
            style="font-weight: bolder"
            @click="onRemoveSelectedSpiders"
          >
            {{ $t('Remove') }}
          </el-button>
          <el-button
            size="small"
            type="success"
            icon="el-icon-plus"
            class="btn add"
            style="font-weight: bolder"
            @click="onAdd"
          >
            {{ $t('Add Spider') }}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--tabs-->
      <!--      <el-tabs v-model="filter.type" @tab-click="onClickTab" class="tabs">-->
      <!--        <el-tab-pane :label="$t('All')" name="all" class="all"></el-tab-pane>-->
      <!--        <el-tab-pane :label="$t('Customized')" name="customized" class="customized"></el-tab-pane>-->
      <!--        <el-tab-pane :label="$t('Configurable')" name="configurable" class="configurable"></el-tab-pane>-->
      <!--        <el-tab-pane :label="$t('Long Task')" name="long-task" class="long-task"></el-tab-pane>-->
      <!--      </el-tabs>-->
      <!--./tabs-->

      <!--legend-->
      <status-legend />
      <!--./legend-->

      <!--table list-->
      <el-table
        ref="table"
        :data="spiderList"
        class="table"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        row-key="id"
        border
        @selection-change="onSpiderSelect"
      >
        <el-table-column
          type="selection"
          width="45"
          align="center"
          reserve-selection
        />
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'type'"
            :key="col.name"
            :property="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            align="left"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ $t(scope.row.type) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_5_errors'"
            :key="col.name"
            :label="$t(col.label)"
            :width="col.width"
            align="center"
          >
            <template slot-scope="scope">
              <div :style="{color:scope.row[col.name]>0?'red':''}">
                {{ scope.row[col.name] }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'cmd'"
            :key="col.name"
            :label="$t(col.label)"
            :width="col.width"
            align="left"
          >
            <template slot-scope="scope">
              <el-input v-model="scope.row[col.name]" />
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name.match(/_ts$/)"
            :key="col.name"
            :label="$t(col.label)"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getTime(scope.row[col.name]) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_status'"
            :key="col.name"
            :label="$t(col.label)"
            align="left"
            :width="col.width"
          >
            <template slot-scope="scope">
              <template
                v-if="scope.row.last_status === 'ERROR'"
              >
                <el-tooltip :content="scope.row.last_error" placement="top">
                  <status-tag :status="scope.row.last_status" />
                </el-tooltip>
              </template>
              <status-tag v-else :status="scope.row.last_status" />
            </template>
          </el-table-column>
          <!--          <el-table-column-->
          <!--            v-else-if="['is_scrapy', 'is_long_task'].includes(col.name)"-->
          <!--            :key="col.name"-->
          <!--            :label="$t(col.label)"-->
          <!--            align="left"-->
          <!--            :width="col.width"-->
          <!--          >-->
          <!--            <template slot-scope="scope">-->
          <!--              <el-switch-->
          <!--                v-if="scope.row.type === 'customized'"-->
          <!--                v-model="scope.row[col.name]"-->
          <!--                active-color="#13ce66"-->
          <!--                disabled-->
          <!--              />-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <el-table-column
            v-else-if="col.name === 'latest_tasks'"
            :key="col.name"
            :label="$t(col.label)"
            :width="col.width"
            :align="col.align"
            class-name="latest-tasks"
          >
            <template slot-scope="scope">
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'pending') > 0"
                type="primary"
                size="small"
              >
                <i class="el-icon-loading" />
                {{ getTaskCountByStatus(scope.row, 'pending') }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'running') > 0"
                type="warning"
                size="small"
              >
                <i class="el-icon-loading" />
                {{ getTaskCountByStatus(scope.row, 'running') }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'finished') > 0"
                type="success"
                size="small"
              >
                <i class="el-icon-check" />
                {{ getTaskCountByStatus(scope.row, 'finished') }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'error') > 0"
                type="danger"
                size="small"
              >
                <i class="el-icon-error" />
                {{ getTaskCountByStatus(scope.row, 'error') }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'cancelled') > 0"
                type="info"
                size="small"
              >
                <i class="el-icon-video-pause" />
                {{ getTaskCountByStatus(scope.row, 'cancelled') }}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'abnormal') > 0"
                type="danger"
                size="small"
              >
                <i class="el-icon-warning" />
                {{ getTaskCountByStatus(scope.row, 'abnormal') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :key="col.name"
            :property="col.name"
            :label="$t(col.label)"
            :align="col.align || 'left'"
            :width="col.width"
          />
        </template>
        <el-table-column :label="$t('Action')" align="left" fixed="right" width="170">
          <template slot-scope="scope">
            <!--            <el-tooltip :content="$t('View')" placement="top">-->
            <!--              <el-button-->
            <!--                type="primary"-->
            <!--                icon="el-icon-search"-->
            <!--                size="mini"-->
            <!--                :disabled="isDisabled(scope.row)"-->
            <!--                @click="onView(scope.row, $event)"-->
            <!--              />-->
            <!--            </el-tooltip>-->
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button
                type="danger"
                icon="el-icon-delete"
                size="mini"
                :disabled="isDisabled(scope.row)"
                @click="onRemove(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip v-if="!isShowRun(scope.row)" :content="$t('No command line')" placement="top">
              <el-button
                disabled
                type="success"
                icon="fa fa-bug"
                size="mini"
                @click="onCrawl(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip v-else :content="$t('Run')" placement="top">
              <el-button
                type="success"
                icon="fa fa-bug"
                size="mini"
                :disabled="isDisabled(scope.row)"
                @click="onCrawl(scope.row, $event)"
              />
            </el-tooltip>
            <!--            <el-tooltip :content="$t('Latest Tasks')" placement="top">-->
            <!--              <el-button-->
            <!--                type="warning"-->
            <!--                icon="fa fa-tasks"-->
            <!--                size="mini"-->
            <!--                :disabled="isDisabled(scope.row)"-->
            <!--                @click="onViewRunningTasks(scope.row, $event)"-->
            <!--              />-->
            <!--            </el-tooltip>-->
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          :current-page.sync="pagination.pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pagination.pageSize"
          layout="sizes, prev, pager, next"
          :total="spiderTotal"
          @current-change="onPageNumChange"
          @size-change="onPageSizeChange"
        />
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import dayjs from 'dayjs'
  import CrawlConfirmDialog from '../../components/Common/CrawlConfirmDialog'
  import StatusTag from '../../components/Status/StatusTag'
  import StatusLegend from '../../components/Status/StatusLegend'

  export default {
    name: 'SpiderList',
    components: {
      StatusLegend,
      CrawlConfirmDialog,
      StatusTag
    },
    data() {
      return {
        pagination: {
          pageNum: 1,
          pageSize: 10
        },
        importLoading: false,
        addConfigurableLoading: false,
        isEditMode: false,
        dialogVisible: false,
        addDialogVisible: false,
        crawlConfirmDialogVisible: false,
        isRunningTasksDialogVisible: false,
        activeSpiderId: undefined,
        activeSpider: undefined,
        filter: {
          keyword: '',
          type: 'all'
        },
        types: [],
        spiderFormRules: {
          name: [{ required: true, message: 'Required Field', trigger: 'change' }]
        },
        fileList: [],
        activeTabName: 'customized',
        tourSteps: [
          // {
          //   target: '#tab-customized',
          //   content: this.$t('View a list of <strong>Customized Spiders</strong>'),
          //   params: {
          //     highlight: false
          //   }
          // },
          // {
          //   target: '#tab-configurable',
          //   content: this.$t('View a list of <strong>Configurable Spiders</strong>'),
          //   params: {
          //     highlight: false
          //   }
          // },
          // {
          //   target: '.table',
          //   content: this.$t('You can view your created spiders here.<br>Click a table row to view <strong>spider details</strong>.'),
          //   params: {
          //     placement: 'top'
          //   }
          // },
          {
            target: '.btn.add',
            content: this.$t('Click to add a new spider.')
          }
        ],
        tourCallbacks: {
          onStop: () => {
            this.$utils.tour.finishTour('spider-list')
          },
          onPreviousStep: (currentStep) => {
            this.$utils.tour.prevStep('spider-list', currentStep)
          },
          onNextStep: (currentStep) => {
            this.$utils.tour.nextStep('spider-list', currentStep)
          }
        },
        tourAddSteps: [
          // {
          //   target: '#tab-customized',
          //   content: this.$t('<strong>Customized Spider</strong> is a highly customized spider, which is able to run on any programming language and any web crawler framework.'),
          //   params: {
          //     placement: 'bottom',
          //     highlight: false
          //   }
          // },
          // {
          //   target: '#tab-configurable',
          //   content: this.$t('<strong>Configurable Spider</strong> is a spider defined by config data, aimed at streamlining spider development and improving dev efficiency.'),
          //   params: {
          //     placement: 'bottom',
          //     highlight: false
          //   }
          // },
          {
            target: '#spider-name',
            content: this.$t('Unique identifier for the spider'),
            params: {
              placement: 'right'
            }
          },
          {
            target: '#upload',
            content: this.$t('Upload a zip file containing all spider files to create the spider.'),
            params: {
              placement: 'right'
            }
          }
        ],
        tourAddCallbacks: {
          onStop: () => {
            this.$utils.tour.finishTour('spider-list-add')
          },
          onPreviousStep: (currentStep) => {
            if (currentStep === 7) {
              this.activeTabName = 'customized'
            }
            this.$utils.tour.prevStep('spider-list-add', currentStep)
          },
          onNextStep: (currentStep) => {
            if (currentStep === 6) {
              this.activeTabName = 'configurable'
            }
            this.$utils.tour.nextStep('spider-list-add', currentStep)
          }
        },
        handle: undefined,
        activeSpiderTaskStatus: 'running',
        selectedSpiders: [],
        isStopLoading: false,
        isRemoveLoading: false,
        isMultiple: false
      }
    },
    computed: {
      ...mapState('spider', [
        'importForm',
        'spiderList',
        'spiderForm',
        'spiderTotal',
        'templateList'
      ]),
      ...mapGetters('user', [
        'userInfo',
        'token'
      ]),
      ...mapState('lang', [
        'lang'
      ]),
      uploadForm() {
        return {
          name: this.spiderForm.name,
          description: this.spiderForm.description || ''
        }
      },
      columns() {
        const columns = []
        columns.push({ name: 'name', label: 'Name', width: '160', align: 'left' })
        // columns.push({ name: 'latest_tasks', label: 'Latest Tasks', width: '180' })
        columns.push({ name: 'last_status', label: 'Last Status', width: '120' })
        columns.push({ name: 'last_run_ts', label: 'Last Run', width: '140' })
        columns.push({ name: 'update_ts', label: 'Update Time', width: '140' })
        columns.push({ name: 'create_ts', label: 'Create Time', width: '140' })
        columns.push({ name: 'description', label: 'Description' })
        return columns
      }
    },
    async created() {
      // fetch spider list
      await this.getList()

      // fetch template list
      // await this.$store.dispatch('spider/getTemplateList')

      // periodically fetch spider list
      this.handle = setInterval(() => {
        this.getList()
      }, 15000)
    },
    mounted() {
      const vm = this
      this.$nextTick(() => {
        vm.$store.commit('spider/SET_SPIDER_FORM', this.spiderForm)
      })

      if (!this.$utils.tour.isFinishedTour('spider-list')) {
        this.$utils.tour.startTour(this, 'spider-list')
      }
    },
    destroyed() {
      clearInterval(this.handle)
    },
    methods: {
      onPageSizeChange(val) {
        this.pagination.pageSize = val
        this.getList()
      },
      onPageNumChange(val) {
        this.pagination.pageNum = val
        this.getList()
      },
      onSearch() {
        this.getList()
      },
      onAdd() {
        this.$store.commit('spider/SET_SPIDER_FORM', {
          template: this.templateList[0]
        })
        this.addDialogVisible = true

        setTimeout(() => {
          if (!this.$utils.tour.isFinishedTour('spider-list-add')) {
            this.$utils.tour.startTour(this, 'spider-list-add')
          }
        }, 300)
      },
      onRefresh() {
        this.getList()
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
        this.$confirm(this.$t('Are you sure to delete this spider?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(async() => {
          await this.$store.dispatch('spider/deleteSpider', row.id)
          this.$message({
            type: 'success',
            message: this.$t('Deleted successfully')
          })
          await this.getList()
          this.$st.sendEv('爬虫列表', '删除爬虫')
        })
      },
      onCrawl(row, ev) {
        ev.stopPropagation()
        this.crawlConfirmDialogVisible = true
        this.activeSpiderId = row.id
        this.$st.sendEv('爬虫列表', '点击运行')
      },
      onCrawlConfirm() {
        setTimeout(() => {
          this.getList()
        }, 1000)
      },
      onView(row, ev) {
        ev.stopPropagation()
        this.$router.push('/spiders/' + row.id)
        this.$st.sendEv('爬虫列表', '查看爬虫')
      },
      isShowRun(row) {
        if (!this.isCustomized(row)) return true
        return !!row.cmd
      },
      isCustomized(row) {
        return row.type === 'customized'
      },
      onUploadSuccess(res) {
        // clear fileList
        this.fileList = []

        // fetch spider list
        setTimeout(() => {
          this.getList()
        }, 500)

        this.$message.success(this.$t('Uploaded spider files successfully'))
        this.addDialogVisible = false
      },
      beforeUpload(file) {
        return new Promise((resolve, reject) => {
          this.$refs['addSpiderForm'].validate(res => {
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
      onClickTab(tab) {
        this.filter.type = tab.name
        this.getList()
      },
      async getList() {
        const params = {
          page_num: this.pagination.pageNum,
          page_size: this.pagination.pageSize
        // keyword: this.filter.keyword,
        }
        await this.$store.dispatch('spider/getSpiderList', params)

        // 更新当前爬虫（任务列表）
        this.updateActiveSpider()
      },
      getTasksByStatus(row, status) {
        if (!row.latest_tasks) return []
        return row.latest_tasks.filter(d => d.status === status)
      },
      getTaskCountByStatus(row, status) {
        return this.getTasksByStatus(row, status).length
      },
      updateActiveSpider() {
        if (this.activeSpider) {
          for (let i = 0; i < this.spiderList.length; i++) {
            const spider = this.spiderList[i]
            if (this.activeSpider.id === spider.id) {
              this.activeSpider = spider
            }
          }
        }
      },
      onViewRunningTasks(row, ev) {
        ev.stopPropagation()
        this.activeSpider = row
        this.isRunningTasksDialogVisible = true
      },
      getTasks(row) {
        if (!this.activeSpider.latest_tasks) {
          return []
        }
        return this.activeSpider.latest_tasks
          .filter(d => d.status === this.activeSpiderTaskStatus)
          .map(d => {
            d = JSON.parse(JSON.stringify(d))
            d.create_ts = d.create_ts.match('^0001') ? 'NA' : dayjs(d.create_ts).format('YYYY-MM-DD HH:mm:ss')
            d.start_ts = d.start_ts.match('^0001') ? 'NA' : dayjs(d.start_ts).format('YYYY-MM-DD HH:mm:ss')
            d.finish_ts = d.finish_ts.match('^0001') ? 'NA' : dayjs(d.finish_ts).format('YYYY-MM-DD HH:mm:ss')
            return d
          })
      },
      onViewTask(row) {
        this.$router.push(`/tasks/${row.id}`)
        this.$st.sendEv('爬虫列表', '任务列表', '查看任务')
      },
      async onStop(row, ev) {
        ev.stopPropagation()
        const res = await this.$store.dispatch('task/cancelTask', row.id)
        if (res.data.code === 200) {
          this.$message.success(`Task "${row.id}" has been sent signal to stop`)
          this.getList()
        }
      },
      onSpiderSelect(spiders) {
        this.selectedSpiders = spiders
      },
      async onRemoveSelectedSpiders() {
        this.$confirm(this.$t('Are you sure to delete selected items?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(async() => {
          this.isRemoveLoading = true
          try {
            const res = await this.$request.delete('/spiders', {
              spider_ids: this.selectedSpiders.map(d => d.id)
            })
            if (res.data.code === 200) {
              this.$message.success('Delete successfully')
              this.$refs['table'].clearSelection()
              await this.getList()
            }
          } finally {
            this.isRemoveLoading = false
          }
          this.$st.sendEv('爬虫列表', '批量删除爬虫')
        })
      },
      async onStopSelectedSpiders() {
        this.$confirm(this.$t('Are you sure to stop selected items?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(async() => {
          this.isStopLoading = true
          try {
            const res = await this.$request.post('/spiders-cancel', {
              spider_ids: this.selectedSpiders.map(d => d.id)
            })
            if (res.data.code === 200) {
              this.$message.success('Sent signals to cancel selected tasks')
              await this.getList()
            }
          } finally {
            this.isStopLoading = false
          }
          this.$st.sendEv('爬虫列表', '批量删除爬虫')
        })
      },
      onCrawlSelectedSpiders() {
        this.crawlConfirmDialogVisible = true
        this.isMultiple = true
      },
      onCrawlConfirmDialogClose() {
        this.crawlConfirmDialogVisible = false
        this.isMultiple = false
      },
      isDisabled(row) {
        return row.is_public && row.username !== this.userInfo.username
      }
    }
  }
</script>

<style scoped lang="scss">
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
      background: rgba(64, 158, 255, .1);
      border: 1px solid rgba(64, 158, 255, .1);
    }

    .add-spider-item.success {
      color: #67c23a;
      background: rgba(103, 194, 58, .1);
      border: 1px solid rgba(103, 194, 58, .1);
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
