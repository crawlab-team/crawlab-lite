<script>
import { mapState } from 'vuex'
import TaskList from '@/views/task/TaskList'

export default {
  name: 'ScheduleTaskList',
  extends: TaskList,
  computed: {
    ...mapState('schedule', ['scheduleForm']),
  },
  methods: {
    async getTaskList() {
      this.isFilterSpiderDisabled = true
      this.filter = Object.assign({}, this.filter, {
        spider_id: this.scheduleForm.spider_id,
        schedule_id: this.scheduleForm.id,
      })
      const params = Object.assign({}, this.pagination, this.filter)
      await this.$store.dispatch('task/getTaskList', params)
    },
    open() {
      clearInterval(this.refreshHandle)
      this.getTaskList()
      this.refreshHandle = setInterval(() => {
        this.getTaskList()
      }, 5000)
    },
    close() {
      clearInterval(this.refreshHandle)
    },
  },
}
</script>
